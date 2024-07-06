package main

import (
	"common"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type UpdateMovieRequest struct {
    UUID       string   `json:"uuid"`
    FileType   string   `json:"fileType"`
    FileSize   uint64   `json:"fileSize"`
}

type UpdateMovieResponse struct {
    Url    string    `json:"url"`
    Method string    `json:"method"`
}

var s3Client *s3.Client
var s3PresignClient *s3.PresignClient;
var dynamodbClient *dynamodb.Client
const expiration = 3600 // 60m

func updateMovie(ctx context.Context, incomingRequest events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    var event UpdateMovieRequest
    err := json.Unmarshal([]byte(incomingRequest.Body), &event)
    if (err != nil) { return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil }
    if (len(event.UUID) == 0 || len(event.FileType) == 0 || event.FileSize == 0) { return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil }

    // get the movie
    tableName := common.MovieTableName
    getInput := &dynamodb.GetItemInput{
        TableName: &tableName,
        Key: map[string]types.AttributeValue {
            "id": &types.AttributeValueMemberS {
                Value: event.UUID,
            },
        },
    }

    result, err := dynamodbClient.GetItem(context.TODO(), getInput)
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

    // cant update a movie that is not ready
    if (!movie.Video.Ready) { return common.EmptyErrorResponse(http.StatusBadRequest), nil }

    // delete the all three resoltions of the old movie
    for _, res := range []string{common.Resolution1, common.Resolution2, common.Resolution3} {
        filename := fmt.Sprintf("%s/%s.mp4", movie.Video.FileName, res)
        _, err = s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
            Bucket: aws.String(common.VideoBucketName),
            Key: &filename,
        })

        if (err != nil) {
            log.Printf("Error deleting movie %s from s3: %v\n", filename, err)
            return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New(fmt.Sprintf("Error deleting movie from s3: %v\n", err))
        }
    }
    // delete the old folder
    _, err = s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
        Bucket: aws.String(common.VideoBucketName),
        Key: &movie.Video.FileName,
    })

    if (err != nil) {
        log.Printf("Error deleting movie folder %s from s3: %v\n", movie.Video.FileName, err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New(fmt.Sprintf("Error deleting movie from s3: %v\n", err))
    }

    // update the movie item in the table
    timestamp := time.Now().Unix()
    fileName := fmt.Sprintf("%s-%d", movie.UUID, timestamp)

    movie.Video.Ready = false
    movie.Video.FileSize = event.FileSize
    movie.Video.FileType = event.FileType
    movie.Video.LastChangeTimestamp = timestamp
    movie.Video.FileName = fileName

    fileName = fmt.Sprintf("%s%s", fileName, common.OriginalSuffix)

    marshaledVideo, err := attributevalue.MarshalMap(movie.Video)
    if err != nil {
        log.Printf("Error marshaling updated video: %v", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error marshaling video")
    }
    updateInput := &dynamodb.UpdateItemInput{
        TableName: aws.String(common.MovieTableName),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: movie.UUID},
        },
        UpdateExpression: aws.String("SET video = :val"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":val": &types.AttributeValueMemberM{Value: marshaledVideo},
        },
    }
    _, err = dynamodbClient.UpdateItem(context.TODO(), updateInput)
    if err != nil {
        log.Printf("Error putting video: %v", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error putting video")
    }

    // create pre signed url
    request, err := s3PresignClient.PresignPutObject(context.TODO(),
    &s3.PutObjectInput{
        Bucket: aws.String(common.VideoBucketName),
        Key: &fileName,
    },
    func(opts *s3.PresignOptions) {
        opts.Expires = time.Duration(expiration * int64(time.Second))
    })

    if err != nil {
        log.Printf("Error getting presigned url for updating movie for \"%s\": %v", fileName, err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error creating presign url")
    }

    res := UpdateMovieResponse {
        Url: request.URL,
        Method: request.Method,
    }
    resString, err := json.Marshal(res)
    if (err != nil) {
        log.Printf("Error marshaling res: %v", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error marshaling res")
    }

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

    s3Client = s3.NewFromConfig(sdkConfig)
    s3PresignClient = s3.NewPresignClient(s3Client)

    dynamodbClient = dynamodb.NewFromConfig(sdkConfig)

    lambda.Start(updateMovie)
}
