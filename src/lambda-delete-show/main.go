package main

import (
	"common"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

func deleteSeason(season *common.Season) error {
    // delete each episode
    for _, episode := range season.Episodes {
        err := deleteVideo(&episode.Video)
        if (err != nil) {
            log.Printf("%v", err)
            return err
        }
    }
    return nil
}

// NOTE:
// if only uuid is provided - delete the whole show
// if season is provided and no episode - delete season
// if season and episode are provided - delete episode
func deleteShow(ctx context.Context, incomingRequest events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    uuid, ok := incomingRequest.QueryStringParameters["uuid"]
    if (!ok) { return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil }

    var err error

    season, seasonOk := incomingRequest.QueryStringParameters["season"]
    seasonIdx := -1
    if (seasonOk) {
        seasonIdx, err = strconv.Atoi(season)
        if (err != nil) { return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil }
    }

    episode, episodeOk := incomingRequest.QueryStringParameters["episode"]
    episodeIdx := -1
    if (episodeOk) {
        episodeIdx, err = strconv.Atoi(episode)
        if (err != nil) { return common.ErrorResponse(http.StatusBadRequest, "Malformed input"), nil }
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


    // delete the show
    if (!seasonOk && !episodeOk) {
        // check if all episodes are ready
        for _, season := range show.Seasons {
            for _, episode := range season.Episodes {
                if (!episode.Video.Ready) { return common.EmptyErrorResponse(http.StatusBadRequest), nil }
            }
        }

        // delete each season
        for _, season := range show.Seasons {
                err := deleteSeason(&season)
                if (err != nil) {
                    log.Printf("%v", err)
                    return common.EmptyErrorResponse(http.StatusInternalServerError), err
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
    }


    // delete a season
    if (seasonOk && !episodeOk) {
        seasonActualIdx := -1
        for idx, season := range show.Seasons {
            if (season.SeasonNumber == uint64(seasonIdx)) {
                seasonActualIdx = idx

                // check if all videos are ready
                for _, episode := range season.Episodes {
                    if (!episode.Video.Ready) { return common.EmptyErrorResponse(http.StatusBadRequest), nil }
                }

                err := deleteSeason(&season)
                if (err != nil) {
                    log.Printf("%v", err)
                    return common.EmptyErrorResponse(http.StatusInternalServerError), err
                }
            }
        }

        // if the season doesnt exist
        if (seasonActualIdx == -1) { return common.EmptyErrorResponse(http.StatusBadRequest), nil }

        // remove the season from the show in show table
        show.Seasons = append(show.Seasons[:seasonActualIdx], show.Seasons[seasonActualIdx+1:]...)

        seasonAttributeArray, err := attributevalue.MarshalList(show.Seasons)
        if (err != nil) {
            log.Printf("Error marshaling modified seasons : %v", err)
            return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error marshaling modified seasons")
        }

        input := &dynamodb.UpdateItemInput{
            TableName: aws.String(common.ShowTableName),
            Key: map[string]types.AttributeValue{
                "id": &types.AttributeValueMemberS{Value: show.UUID},
            },
            UpdateExpression: aws.String("SET seasons = :val"),
            ExpressionAttributeValues: map[string]types.AttributeValue{
                ":val": &types.AttributeValueMemberL{Value: seasonAttributeArray},
            },
        }

        _, err = dynamodbClient.UpdateItem(context.TODO(), input)
        if err != nil {
            log.Printf("Error updating modified seasons : %v", err)
            return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error updating modified seasons")
        }
    }

    // delete an episode
    if (seasonOk && episodeOk) {
        seasonActualIdx := -1
        episodeActualIdx := -1
        for sIdx, season := range show.Seasons {
            if (season.SeasonNumber == uint64(seasonIdx)) {
                seasonActualIdx = sIdx
                for eIdx, episode := range season.Episodes {
                    if (episode.EpisodeNumber == uint64(episodeIdx)) {
                        episodeActualIdx = eIdx

                        // is episode ready
                        if (!episode.Video.Ready) { return common.EmptyErrorResponse(http.StatusBadRequest), nil }

                        err := deleteVideo(&episode.Video)
                        if (err != nil) {
                            log.Printf("%v", err)
                            return common.EmptyErrorResponse(http.StatusInternalServerError), err
                        }
                    }
                }
            }
        }
        // if the season and episode dont exist
        if (seasonActualIdx == -1 || episodeActualIdx == -1) { return common.EmptyErrorResponse(http.StatusBadRequest), nil }

        // remove the episode from the season
        show.Seasons[seasonActualIdx].Episodes = append(show.Seasons[seasonActualIdx].Episodes[:episodeActualIdx], show.Seasons[seasonActualIdx].Episodes[episodeActualIdx+1:]...)

        seasonAttributeArray, err := attributevalue.MarshalList(show.Seasons)
        if (err != nil) {
            log.Printf("Error marshaling modified seasons : %v", err)
            return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error marshaling modified seasons")
        }

        input := &dynamodb.UpdateItemInput{
            TableName: aws.String(common.ShowTableName),
            Key: map[string]types.AttributeValue{
                "id": &types.AttributeValueMemberS{Value: show.UUID},
            },
            UpdateExpression: aws.String("SET seasons = :val"),
            ExpressionAttributeValues: map[string]types.AttributeValue{
                ":val": &types.AttributeValueMemberL{Value: seasonAttributeArray},
            },
        }

        _, err = dynamodbClient.UpdateItem(context.TODO(), input)
        if err != nil {
            log.Printf("Error updating modified seasons : %v", err)
            return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error updating modified seasons")
        }
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
