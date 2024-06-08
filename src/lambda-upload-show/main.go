package main

import (
	"common"
	"context"
	"errors"
	"fmt"
	"log"
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

type SingleVideo struct {
    Url string               `json:"url"`
    Method string            `json:"method"`
    EpisodeNumber uint64     `json:"episodeNumber"`
    SeasonNumber uint64     `json:"seasonNumber"`
}

type UploadShowResponse struct {
    Videos []SingleVideo   `json:"urls"`
}

var s3PresignClient *s3.PresignClient;
var dynamodbClient *dynamodb.Client
const expiration = 3600 // 60m

func uploadShow(ctx context.Context, event common.Show) (UploadShowResponse, error) {
    showUUID := uuid.New().String()
    event.UUID = showUUID

    res := UploadShowResponse{
        Videos: []SingleVideo{},
    }

    // create pre signed urls for each episode in each season
    for seasonIndex := 0; seasonIndex < len(event.Seasons); seasonIndex++ {
        for episodeIndex := 0; episodeIndex < len(event.Seasons[seasonIndex].Episodes); episodeIndex++ {
            timestamp := time.Now().Unix()
            fileName := fmt.Sprintf("%s-%d-%d-%d.%s",
                showUUID, event.Seasons[seasonIndex].SeasonNumber,
                event.Seasons[seasonIndex].Episodes[episodeIndex].EpisodeNumber,
                timestamp, event.Seasons[seasonIndex].Episodes[episodeIndex].Video.FileType)
            // having '/' in the name causes s3 to treat it as a folder
            fileName = strings.ReplaceAll(fileName, "/", "-")

            event.Seasons[seasonIndex].Episodes[episodeIndex].Video.FileName = fileName

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
                log.Printf("Error getting presigned url for uploading show for \"%s\": %v", fileName, err)
                return UploadShowResponse{}, errors.New("Error creating presign url")
            }

            url := SingleVideo {
                Url: request.URL,
                Method: request.Method,
                SeasonNumber: event.Seasons[seasonIndex].SeasonNumber,
                EpisodeNumber: event.Seasons[seasonIndex].Episodes[episodeIndex].EpisodeNumber,
            }

            res.Videos = append(res.Videos, url)
        }
    }

    // add the show to the db
    marshaledShow, err := attributevalue.MarshalMap(event)
    if err != nil {
        log.Printf("Error marshaling show: %v", err)
        return UploadShowResponse{}, errors.New("Error marshaling show")
    }

    tableName := common.ShowTableName
    input := &dynamodb.PutItemInput{
        TableName: &tableName,
        Item: marshaledShow,
    }

     _, err = dynamodbClient.PutItem(context.TODO(), input)
    if err != nil {
        log.Printf("Error putting show: %v", err)
        return UploadShowResponse{}, errors.New("Error putting show")
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

    lambda.Start(uploadShow)
}
