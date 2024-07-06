package main

import (
	"common"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3Client *s3.Client
var dynamodbClient *dynamodb.Client

func deleteMovie(ctx context.Context, incomingRequest events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    uuid, ok := incomingRequest.QueryStringParameters["uuid"]
    if (!ok) {
        return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil
    }

    tableName := common.MovieTableName
    getMovieInput := &dynamodb.GetItemInput{
        TableName: &tableName,
        Key: map[string]types.AttributeValue {
            "id": &types.AttributeValueMemberS {
                Value: uuid,
            },
        },
    }

    result, err := dynamodbClient.GetItem(context.TODO(), getMovieInput)
    if err != nil {
        log.Printf("Error getting movie: %v", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error getting movie")
    }
    if result.Item == nil { return common.ErrorResponse(http.StatusBadRequest, "Movie does not exist"), nil }

    var movie common.Movie
    err = attributevalue.UnmarshalMap(result.Item, &movie)
    if err != nil {
        log.Printf("Error unmarshaling movie :%v", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error umarshaling movie")
    }

    if (!movie.Video.Ready) { return common.EmptyErrorResponse(http.StatusBadRequest), nil }

    // delete each resolution
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

    // delete the folder
    _, err = s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
        Bucket: aws.String(common.VideoBucketName),
        Key: &movie.Video.FileName,
    })

    if (err != nil) {
        log.Printf("Error deleting movie folder %s from s3: %v\n", movie.Video.FileName, err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New(fmt.Sprintf("Error deleting movie from s3: %v\n", err))
    }

    // delete the item from movie table
    deleteInput := &dynamodb.DeleteItemInput{
        TableName: aws.String(tableName),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS {
                Value: uuid,
            },
        },
    }

    _, err = dynamodbClient.DeleteItem(context.TODO(), deleteInput)
    if (err != nil) {
        log.Printf("Error deleting movie from dynamodb: %v\n", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New(fmt.Sprintf("Error deleting movie from dynamodb: %v\n", err))
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
    }, nil
}

func main() {
    sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
    if err != nil {
        log.Fatal("Cannot load in default config")
    }
    s3Client = s3.NewFromConfig(sdkConfig)
    dynamodbClient = dynamodb.NewFromConfig(sdkConfig)
    lambda.Start(deleteMovie)
}
