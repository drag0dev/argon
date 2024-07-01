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

func deleteVideo(video *common.Video) error {
    // delete each resolution
    for _, res := range []string{common.Resolution1, common.Resolution2, common.Resolution3} {
        filename := fmt.Sprintf("%s/%s.mp4", video.FileName, res)
        _, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
            Bucket: aws.String(common.VideoBucketName),
            Key: &filename,
        })

        if (err != nil) {
            log.Printf("Error deleting episode %s from s3: %v\n", filename, err)
            return errors.New(fmt.Sprintf("Error deleting episode from s3: %v\n", err))
        }
    }

    // delete the folder
    _, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
        Bucket: aws.String(common.VideoBucketName),
        Key: &video.FileName,
    })

    if (err != nil) {
        log.Printf("Error deleting episode folder %s from s3: %v\n", video.FileName, err)
        return errors.New(fmt.Sprintf("Error deleting episode from s3: %v\n", err))
    }

    return nil
}

func deleteShow(ctx context.Context, incomingRequest events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    uuid, ok := incomingRequest.QueryStringParameters["uuid"]
    if (!ok) {
        return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil
    }

    tableName := common.ShowTableName
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
        log.Printf("Error getting show: %v", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error getting show")
    }
    if result.Item == nil { return common.ErrorResponse(http.StatusBadRequest, "Show does not exist"), nil }

    var show common.Show
    err = attributevalue.UnmarshalMap(result.Item, &show)
    if err != nil {
        log.Printf("Error unmarshaling show :%v", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error umarshaling show")
    }

    // check if all episodes are ready
    for _, season := range show.Seasons {
        for _, episode := range season.Episodes {
            if (!episode.Video.Ready) { return common.EmptyErrorResponse(http.StatusBadRequest), nil }
        }
    }

    // delete each episode
    for _, season := range show.Seasons {
        for _, episode := range season.Episodes {
            err := deleteVideo(&episode.Video)
            if (err != nil) {
                log.Printf("%v", err)
                return common.EmptyErrorResponse(http.StatusInternalServerError), err
            }
        }
    }

    // delete the item from show table
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
        log.Printf("Error deleting show from dynamodb: %v\n", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New(fmt.Sprintf("Error deleting show from dynamodb: %v\n", err))
    }

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
    }, nil
}

func main() {
    sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
    if err != nil {
        log.Fatal("Cannot load in default config")
    }
    s3Client = s3.NewFromConfig(sdkConfig)
    dynamodbClient = dynamodb.NewFromConfig(sdkConfig)
    lambda.Start(deleteShow)
}
