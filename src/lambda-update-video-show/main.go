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

type UpdateShowRequest struct {
    UUID        string    `json:"uuid"`
    FileType    string    `json:"fileType"`
    FileSize    uint64    `json:"fileSize"`
    Season      uint64    `json:"season"`
    Episode     uint64    `json:"episode"`
}

type UpdateShowResponse struct {
    Url    string            `json:"url"`
    Method string            `json:"method"`
}

var s3Client *s3.Client
var s3PresignClient *s3.PresignClient;
var dynamodbClient *dynamodb.Client
const expiration = 3600 // 60m

func updateShow(ctx context.Context, incomingRequest events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    var event UpdateShowRequest
    err := json.Unmarshal([]byte(incomingRequest.Body), &event)
    if (err != nil) { return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil }
    if (len(event.UUID) == 0 || len(event.FileType) == 0 || event.FileSize == 0 ||
            event.Season == 0 || event.Episode == 0) { return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil }

    // get the show
    tableName := common.ShowTableName
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

    seasonActualIdx := -1
    episodeActualIdx := -1
    for seasonIdx, season := range show.Seasons {
        if (season.SeasonNumber == event.Season) {
            seasonActualIdx = seasonIdx
            for episodeIdx, episode := range season.Episodes {
                if (episode.EpisodeNumber == event.Episode) {
                    episodeActualIdx = episodeIdx
                }
            }
        }
    }
    // non existant season and episode
    if (seasonActualIdx == -1 || episodeActualIdx == -1) { return common.EmptyErrorResponse(http.StatusBadRequest), nil }

    // is episode ready
    if (!show.Seasons[seasonActualIdx].Episodes[episodeActualIdx].Video.Ready) { return common.EmptyErrorResponse(http.StatusBadRequest), nil }


    episodeFileName := show.Seasons[seasonActualIdx].Episodes[episodeActualIdx].Video.FileName
    // delete all three resolutions of the old episode
    for _, res := range []string{common.Resolution1, common.Resolution2, common.Resolution3} {
        filename := fmt.Sprintf("%s/%s.mp4", episodeFileName, res)
        _, err = s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
            Bucket: aws.String(common.VideoBucketName),
            Key: &filename,
        })

        if (err != nil) {
            log.Printf("Error deleting episode %s from s3: %v\n", filename, err)
            return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New(fmt.Sprintf("Error deleting episode from s3: %v\n", err))
        }
    }
    // delete the old folder
    _, err = s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
        Bucket: aws.String(common.VideoBucketName),
        Key: &episodeFileName,
    })

    if (err != nil) {
        log.Printf("Error deleting episode folder %s from s3: %v\n", episodeFileName, err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New(fmt.Sprintf("Error deleting episode from s3: %v\n", err))
    }

    // update the show item in the table
    video := &show.Seasons[seasonActualIdx].Episodes[episodeActualIdx].Video

    parts := strings.Split(video.FileName, "-")
    // all episodes in a tv show carry timestamp in the name that represents the show creation timestamp
    oldTimestamp := parts[len(parts)-1]
    fileName := fmt.Sprintf("%s-%d-%d-%s", show.UUID, event.Season, event.Episode, oldTimestamp)

    video.Ready = false
    video.FileSize = event.FileSize
    video.FileType = event.FileType
    video.LastChangeTimestamp = time.Now().Unix()
    video.FileName = fileName

    fileName = fmt.Sprintf("%s%s", fileName, common.OriginalSuffix)

    marshaledSeasons, err := attributevalue.MarshalList(show.Seasons)
    if err != nil {
        log.Printf("Error marshaling seasons show: %v", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error marshaling seasons")
    }
    updateInput := &dynamodb.UpdateItemInput{
        TableName: aws.String(common.ShowTableName),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: show.UUID},
        },
        UpdateExpression: aws.String("SET seasons = :val"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":val": &types.AttributeValueMemberL{Value: marshaledSeasons},
        },
    }
    _, err = dynamodbClient.UpdateItem(context.TODO(), updateInput)
    if err != nil {
        log.Printf("Error putting seasons: %v", err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error putting seasons")
    }

    // create presigned url
    request, err := s3PresignClient.PresignPutObject(context.TODO(),
    &s3.PutObjectInput{
        Bucket: aws.String(common.VideoBucketName),
        Key: &fileName,
    },
    func(opts *s3.PresignOptions) {
        opts.Expires = time.Duration(expiration * int64(time.Second))
    })

    if err != nil {
        log.Printf("Error getting presigned url for updating episode for \"%s\": %v", fileName, err)
        return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error creating presign url")
    }

    res := UpdateShowResponse {
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

    lambda.Start(updateShow)
}
