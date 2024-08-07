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
	"strconv"
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

type GetShowEvent struct {
    UUID string    `json:"uuid"`
    Season uint64  `json:"season"`
    Episode uint64 `json:"episode"`
}

type GetShowResponse struct {
    Url string               `json:"url"`
    Method string            `json:"method"`
    Data common.Show         `json:"show"`
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

func getShow(ctx context.Context, incomingRequest events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    uuid, ok := incomingRequest.QueryStringParameters["uuid"]
    if (!ok) { return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil }
    _, ok = incomingRequest.QueryStringParameters["season"]
    if (!ok) { return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil }
    _, ok = incomingRequest.QueryStringParameters["episode"]
    if (!ok) { return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil }
    resolution, ok := incomingRequest.QueryStringParameters["resolution"]
    if (!ok) {
        return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil
    }
    if (resolution != common.Resolution1 && resolution != common.Resolution2 && resolution != common.Resolution3)  {
        return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil
    }

    season, err := strconv.Atoi(incomingRequest.QueryStringParameters["season"])
    if (err != nil) { return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil }

    episode, err := strconv.Atoi(incomingRequest.QueryStringParameters["episode"])
    if (err != nil) { return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil }

    tableName := common.ShowTableName
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
        log.Printf("Error getting show: %v", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error getting show")
    }

    if result.Item == nil {
        return common.ErrorResponse(http.StatusBadRequest, "Show does not exist"), nil
    }

    var show common.Show
    err = attributevalue.UnmarshalMap(result.Item, &show)
    if err != nil {
        log.Printf("Error unmarshaling show :%v", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error umarshaling show")
    }

    var filename string = ""
    seasonActualIndex := -1
    episodeActualIndex := -1
    outer: for seasonIndex := 0; seasonIndex < len(show.Seasons); seasonIndex++ {
        for episodeIndex := 0; episodeIndex < len(show.Seasons[seasonIndex].Episodes); episodeIndex++ {
            if (
            uint64(season) == show.Seasons[seasonIndex].SeasonNumber &&
            uint64(episode) == show.Seasons[seasonIndex].Episodes[episodeIndex].EpisodeNumber) {
                filename = show.Seasons[seasonIndex].Episodes[episodeIndex].Video.FileName
                seasonActualIndex = seasonIndex
                episodeActualIndex = episodeIndex

                // check if the video is processed
                if (!show.Seasons[seasonIndex].Episodes[episodeIndex].Video.Ready) {
                    return common.EmptyErrorResponse(http.StatusBadRequest), nil
                }
                break outer;
            }
        }
    }

    if filename == "" {
        return common.ErrorResponse(http.StatusBadRequest,"Episode does not exist"), nil
    }

    bucketName := common.VideoBucketName
    filename = fmt.Sprintf("%s/%s.mp4", filename, resolution)
    request, err := s3PresignClient.PresignGetObject(context.TODO(),
        &s3.GetObjectInput{
            Bucket: &bucketName,
            Key: &filename,
        },
        func (opts *s3.PresignOptions) {
            opts.Expires = time.Duration(expiration * int64(time.Second))
    })

    if err != nil {
        log.Printf("Error creating presigned url for getting episode: %v", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error creating presigned url for getting episode")
    }

    res := GetShowResponse{
        Url: request.URL,
        Method: request.Method,
        Data: show,
    }
    resString, err := json.Marshal(res)
    if (err != nil) {
        log.Printf("Error marshaling res: %v", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error marshaling res")
    }

    prefChangeItem := common.PreferenceChange{
        UpdateWeight: common.GetUpdateWeight,
        ChangeWeight: common.GetChangeWeight,
        Actors: show.Seasons[seasonActualIndex].Episodes[episodeActualIndex].Actors,
        Directors: show.Seasons[seasonActualIndex].Episodes[episodeActualIndex].Directors,
        Genres: show.Genres,
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

    lambda.Start(getShow)
}
