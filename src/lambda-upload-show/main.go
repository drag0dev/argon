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

func uploadShow(ctx context.Context, incomingRequest events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    var event common.Show
    err := json.Unmarshal([]byte(incomingRequest.Body), &event)
    if (err != nil) { return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil }
    if (!common.IsShowValid(&event)) { return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil }

    showUUID := uuid.New().String()
    event.UUID = showUUID

    res := UploadShowResponse{
        Videos: []SingleVideo{},
    }

    // create pre signed urls for each episode in each season
    for seasonIndex := 0; seasonIndex < len(event.Seasons); seasonIndex++ {
        for episodeIndex := 0; episodeIndex < len(event.Seasons[seasonIndex].Episodes); episodeIndex++ {
            timestamp := time.Now().Unix()
            fileName := fmt.Sprintf("%s-%d-%d-%d",
                showUUID, event.Seasons[seasonIndex].SeasonNumber,
                event.Seasons[seasonIndex].Episodes[episodeIndex].EpisodeNumber,
                timestamp)

            event.Seasons[seasonIndex].Episodes[episodeIndex].Video.FileName = fileName
            event.Seasons[seasonIndex].Episodes[episodeIndex].Video.Ready = false

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
                log.Printf("Error getting presigned url for uploading show for \"%s\": %v", fileName, err)
                return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error creating presign url")
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
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error marshaling show")
    }

    tableName := common.ShowTableName
    input := &dynamodb.PutItemInput{
        TableName: &tableName,
        Item: marshaledShow,
    }

     _, err = dynamodbClient.PutItem(context.TODO(), input)
    if err != nil {
        log.Printf("Error putting show: %v", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error putting show")
    }

    resString, err := json.Marshal(res)
    if (err != nil) {
        log.Printf("Error marshaling res: %v", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error marshaling res")
    }

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Headers:    map[string]string{"Content-Type": "application/json"},
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

    lambda.Start(uploadShow)
}
