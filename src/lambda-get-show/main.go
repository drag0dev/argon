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

type GetShowEvent struct {
    UUID string    `json:"uuid"`
    Season uint64  `json:"season"`
    Episode uint64 `json:"episode"`
}

type GetShowResponse struct {
    Url string               `json:"url"`
    Method string            `json:"method"`
}

var s3PresignClient *s3.PresignClient;
var dynamodbClient *dynamodb.Client
const expiration = 300 // 5m

func getShow(ctx context.Context, event GetShowEvent) (GetShowResponse, error) {
    tableName := common.ShowTableName
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
        log.Printf("Error getting show: %v", err)
        return GetShowResponse{}, errors.New("Error getting show")
    }

    if result.Item == nil {
        return GetShowResponse{}, errors.New("Show does not exist")
    }

    var show common.Show
    err = attributevalue.UnmarshalMap(result.Item, &show)
    if err != nil {
        log.Printf("Error unmarshaling show :%v", err)
        return GetShowResponse{}, errors.New("Error umarshaling show")
    }

    var filename string = ""
    outer: for seasonIndex := 0; seasonIndex < len(show.Seasons); seasonIndex++ {
        for episodeIndex := 0; episodeIndex < len(show.Seasons[seasonIndex].Episodes); episodeIndex++ {
            if (
            event.Season == show.Seasons[seasonIndex].SeasonNumber &&
            event.Episode == show.Seasons[seasonIndex].Episodes[episodeIndex].EpisodeNumber) {
                filename = show.Seasons[seasonIndex].Episodes[episodeIndex].Video.FileName
                break outer;
            }
        }
    }

    if filename == "" {
        return GetShowResponse{}, errors.New("Episode does not exist")
    }

    bucketName := common.VideoBucketName
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
        return GetShowResponse{}, errors.New("Error creating presigned url for getting episode")
    }

    res := GetShowResponse{
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

    lambda.Start(getShow)
}
