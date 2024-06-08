package main

import (
	"common"
	"context"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type GetMovieEvent struct {
    UUID string `json:"uuid"`
}

type GetMovieResponse struct {
    Url string               `json:"url"`
    Method string            `json:"method"`
}

var s3PresignClient *s3.PresignClient;
var dynamodbClient *dynamodb.Client
const expiration = 300 // 5m

func getMovie(ctx context.Context, event GetMovieEvent) (GetMovieResponse, error) {
    tableName := common.MovieTableName
    input := &dynamodb.GetItemInput{
        TableName: &tableName,
        Key: map[string]types.AttributeValue {
            "id": &types.AttributeValueMemberS {
                Value: event.UUID,
            },
        },
    }

    result, err := dynamodbClient.GetItem(context.TODO(), input)
    if err != nil {
        log.Printf("Error getting movie: %v", err)
        return GetMovieResponse{}, errors.New("Error getting movie")
    }

    if result.Item == nil {
        return GetMovieResponse{}, errors.New("Movie does not exist")
    }

    var movie common.Movie
    err = attributevalue.UnmarshalMap(result.Item, &movie)
    if err != nil {
        log.Printf("Error unmarshaling movie :%v", err)
        return GetMovieResponse{}, errors.New("Error umarshaling movie")
    }

    bucketName := common.VideoBucketName
    request, err := s3PresignClient.PresignGetObject(context.TODO(),
        &s3.GetObjectInput{
            Bucket: &bucketName,
            Key: &movie.Video.FileName,
        },
        func (opts *s3.PresignOptions) {
            opts.Expires = time.Duration(expiration * int64(time.Second))
    })

    if err != nil {
        log.Printf("Error creating presigned url for getting movie: %v", err)
        return GetMovieResponse{}, errors.New("Error creating presigned url for getting movie")
    }

    res := GetMovieResponse{
        Url: request.URL,
        Method: request.Method,
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

    lambda.Start(getMovie)
}
