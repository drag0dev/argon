package main

import (
	"common"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/lestrrat-go/jwx/jwt"
)

type GetMovieResponse struct {
    Url string               `json:"url"`
    Method string            `json:"method"`
    Data common.Movie        `json:"movie"`
}

var s3PresignClient *s3.PresignClient;
var dynamodbClient *dynamodb.Client
var preferenceChangeQueueClient *sqs.Client
const expiration = 300 // 5m

func enqueChangePreferenceItem(prefChangeItem common.PreferenceChange, headerVal string) {
    token := strings.TrimPrefix(headerVal, "Bearer ")
    if (token == "") {
        log.Printf("Missing token")
        return
    }

    parsedToken, err := jwt.Parse([]byte(token))
    if (err != nil) {
        log.Printf("Error parsing token: %v", err)
        return
    }

    sub, ok := parsedToken.Get("sub")
    if !ok {
        log.Println("sub claim not found in token")
        return
    }

    userId, ok := sub.(string)
    if !ok {
        log.Println("userid is not string")
        return
    }
    prefChangeItem.UserId = userId

    prefChangeMarshaled, err := json.Marshal(prefChangeItem)
    if (err != nil ) {
        log.Printf("Error marshaling preference change item: %v", err)
        return
    }

    queueUrl, err := preferenceChangeQueueClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
        QueueName: aws.String(common.PreferenceUpdateQueue),
    })
    if err != nil {
        log.Printf("Error getting queue url: %v", err)
        return
    }

    sendInput := &sqs.SendMessageInput{
        MessageBody: aws.String(string(prefChangeMarshaled)),
        QueueUrl:    queueUrl.QueueUrl,
    }
    _, err = preferenceChangeQueueClient.SendMessage(context.TODO(), sendInput)
    if err != nil { log.Printf("Error enquing preference change item: %v", err) }
}

func getMovie(ctx context.Context, incomingRequest events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    uuid, ok := incomingRequest.QueryStringParameters["uuid"]
    if (!ok) {
        return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil
    }
    resolution, ok := incomingRequest.QueryStringParameters["resolution"]
    if (!ok) {
        return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil
    }

    if (resolution != common.Resolution1 && resolution != common.Resolution2 && resolution != common.Resolution3)  {
        return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil
    }

    tableName := common.MovieTableName
    input := &dynamodb.GetItemInput{
        TableName: &tableName,
        Key: map[string]types.AttributeValue {
            "id": &types.AttributeValueMemberS {
                Value: uuid,
            },
        },
    }

    result, err := dynamodbClient.GetItem(context.TODO(), input)
    if err != nil {
        log.Printf("Error getting movie: %v", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error getting movie")
    }

    if result.Item == nil {
        return common.ErrorResponse(http.StatusBadRequest, "Movie does not exist"), nil
    }

    var movie common.Movie
    err = attributevalue.UnmarshalMap(result.Item, &movie)
    if err != nil {
        log.Printf("Error unmarshaling movie :%v", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error umarshaling movie")
    }

    if (!movie.Video.Ready) {
        return common.EmptyErrorResponse(http.StatusBadRequest), nil
    }

    bucketName := common.VideoBucketName
    fileName := fmt.Sprintf("%s/%s.mp4", movie.Video.FileName, resolution)
    request, err := s3PresignClient.PresignGetObject(context.TODO(),
        &s3.GetObjectInput{
            Bucket: &bucketName,
            Key: &fileName,
        },
        func (opts *s3.PresignOptions) {
            opts.Expires = time.Duration(expiration * int64(time.Second))
    })

    if err != nil {
        log.Printf("Error creating presigned url for getting movie: %v", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error creating presigned url for getting movie")
    }

    res := GetMovieResponse{
        Url: request.URL,
        Method: request.Method,
        Data: movie,
    }
    resString, err := json.Marshal(res)
    if (err != nil) {
        log.Printf("Error marshaling res response: %v", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error marshaling res response for getting movie")
    }

    prefChangeItem := common.PreferenceChange{
        UpdateWeight: common.GetUpdateWeight,
        ChangeWeight: common.GetChangeWeight,
        Actors: movie.Actors,
        Directors: movie.Directors,
        Genres: movie.Genres,
    }
    enqueChangePreferenceItem(prefChangeItem, incomingRequest.Headers["Authorization"])

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Headers: map[string]string{
            "Content-Type": "application/json",
            "Access-Control-Allow-Origin": "*",
            "Access-Control-Allow-Credentials": "true",
            "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
            "Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
        },
        Body: string(resString),
    }, nil
}

func main() {
    sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
    if err != nil {
        log.Fatal("Cannot load in default config")
    }

    s3Client := s3.NewFromConfig(sdkConfig)
    s3PresignClient = s3.NewPresignClient(s3Client)
    preferenceChangeQueueClient = sqs.NewFromConfig(sdkConfig)
    dynamodbClient = dynamodb.NewFromConfig(sdkConfig)

    lambda.Start(getMovie)
}
