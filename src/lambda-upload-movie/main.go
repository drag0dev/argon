package main

import (
	"common"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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
    SignedHeader http.Header `json:"signedHeader"`
}

var s3PresignClient *s3.PresignClient;
var dynamodbClient *dynamodb.Client
const expiration = 3600 // 60m

func uploadMovie(ctx context.Context, event common.Movie) (UploadMovieResponse, error) {
    movieUUID := uuid.New().String()
    event.UUID = movieUUID

    timestamp := time.Now().Unix()
    fileName := fmt.Sprintf("%s-%d.%s", movieUUID, timestamp, event.Video.FileType)
    // having '/' in the name causes s3 to treat it as a folder
    fileName = strings.ReplaceAll(fileName, "/", "-")

    event.Video.FileName = fileName

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
        return UploadMovieResponse{}, errors.New("Error creating presign url")
    }

    // add the movie to the db
    marshaledMovie, err := attributevalue.MarshalMap(event)
    if err != nil {
        log.Printf("Error marshaling movie: %v", err)
        return UploadMovieResponse{}, errors.New("Error marshaling movie")
    }

    tableName := common.MovieTableName
    input := &dynamodb.PutItemInput{
        TableName: &tableName,
        Item: marshaledMovie,
    }

     _, err = dynamodbClient.PutItem(context.TODO(), input)
    if err != nil {
        log.Printf("Error putting movie: %v", err)
        return UploadMovieResponse{}, errors.New("Error putting movie")
    }

    res := UploadMovieResponse {
        Url: request.URL,
        Method: request.Method,
        SignedHeader: request.SignedHeader,
    }

    return res, nil
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
