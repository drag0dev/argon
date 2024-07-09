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
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type UploadMovieResponse struct {
    Url string               `json:"url"`
    Method string            `json:"method"`
}

var s3PresignClient *s3.PresignClient;
var dynamodbClient *dynamodb.Client
const expiration = 3600 // 60m

func uploadMovie(ctx context.Context, incomingRequest events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    var event common.Movie
    err := json.Unmarshal([]byte(incomingRequest.Body), &event)
    if (err != nil) { return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil }
    if (!common.IsMovieValid(&event)) { return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil }

    movieUUID := uuid.New().String()
    event.UUID = movieUUID

    event.Video.Ready = false

    timestamp := time.Now().Unix()
    fileName := fmt.Sprintf("%s-%d", movieUUID, timestamp)

    event.Video.FileName = fileName
    fileName = fmt.Sprintf("%s%s", fileName, common.OriginalSuffix)

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
        log.Printf("Error getting presigned url for uploading movie for \"%s\": %v", fileName, err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error creating presign url")
    }

    // add the movie to the db
    marshaledMovie, err := attributevalue.MarshalMap(event)
    if err != nil {
        log.Printf("Error marshaling movie: %v", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error marshaling movie")
    }

    tableName := common.MovieTableName
    input := &dynamodb.PutItemInput{
        TableName: &tableName,
        Item: marshaledMovie,
    }

     _, err = dynamodbClient.PutItem(context.TODO(), input)
    if err != nil {
        log.Printf("Error putting movie: %v", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error putting movie")
    }

    res := UploadMovieResponse {
        Url: request.URL,
        Method: request.Method,
    }
    resString, err := json.Marshal(res)
    if (err != nil) {
        log.Printf("Error marshaling res: %v", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error marshaling movie")
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

    s3Client := s3.NewFromConfig(sdkConfig)
    s3PresignClient = s3.NewPresignClient(s3Client)

    dynamodbClient = dynamodb.NewFromConfig(sdkConfig)

    lambda.Start(uploadMovie)
}
