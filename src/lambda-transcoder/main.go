package main

import (
	"common"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)


var s3Client *s3.Client
var uploader *manager.Uploader
var downloader *manager.Downloader
var dynamodbClient *dynamodb.Client
var snsClient *sns.Client
var publishingTopicARN string

func handler(ctx context.Context, s3Event events.S3Event) error {
    for _, record := range s3Event.Records {
        bucket := record.S3.Bucket.Name
        key := record.S3.Object.Key
        log.Printf("Processing %s - %s", bucket, key)

        videoName, _ := strings.CutSuffix(key, common.OriginalSuffix)
        update := false
        if (strings.HasSuffix(videoName, common.UpdateSuffix)) {
            update = true
            videoName, _ = strings.CutSuffix(videoName, common.UpdateSuffix)
        }

        inputFile := "/tmp/input_video"
        outputFile := "/tmp/output.mp4"

        err := downloadFile(ctx, bucket, key, inputFile)
        if (err != nil) {
            log.Printf("failed to download file: %v\n", err)
            return fmt.Errorf("failed to download file %s: %v", key, err)
        }

        resolutions:= []string{common.Resolution1, common.Resolution2, common.Resolution3}
        for _, res := range resolutions {
            os.Remove(outputFile)
            err = transcodeVideo(inputFile, outputFile, fmt.Sprintf("scale=%s", res))
            if (err != nil) {
                log.Printf("failed to transcode video: %v", err)
                return fmt.Errorf("failed to transcode video: %v", err)
            }

            // base_name_of_the_file/resolution.mp4
            outputKey := fmt.Sprintf("%s/%s.mp4", videoName, res)
            err = uploadFile(ctx, bucket, outputKey, outputFile)
            if (err != nil) {
                log.Printf("failed to upload file: %v", err)
                return fmt.Errorf("failed to upload file: %v", err)
            }

        }
        // clean up temporary files
        os.Remove(inputFile)
        os.Remove(outputFile)

        // delete the original
        input := &s3.DeleteObjectInput{
            Bucket: aws.String(bucket),
            Key:    aws.String(key),
        }
        _, err = s3Client.DeleteObject(context.TODO(), input)
        if (err != nil)  {
            log.Printf("Error deleting the original: %v\n", err)
            return errors.New(fmt.Sprintf("Error deleting the original: %v\n", err))
        }

        // set video to ready
        nameParts := strings.Split(videoName, "-")
        uuid := strings.Join(nameParts[:5], "-")
        var show *common.Show
        if (len(nameParts)==6) {
            err = updateMovie(uuid)
            if (err != nil) {
                log.Printf("Error marking movie video ready: %v\n", err)
                return errors.New(fmt.Sprintf("Error marking movie video ready: %v\n", err))
            }
        } else {
            show, err = updateTvShow(uuid, nameParts[5], nameParts[6])
            if (err != nil) {
                log.Printf("Error marking show video ready: %v\n", err)
                return errors.New(fmt.Sprintf("Error marking show video ready: %v\n", err))
            }
        }

        // only emit publish notification if this is is not an update
        if (!update) {
            // emit publish notification
            var notification string
            if (len(nameParts) == 6) {
                notification = "movie~" + uuid
            } else {
                // tv shows can be uploaded in bulk therefore we need to wait for all episodes to be uploaded
                // before emitting a notification
                notification = "show~" + uuid

                showReadyForPublish := true
                outer: for _, season := range show.Seasons {
                    for _, episode := range season.Episodes {
                        if (!episode.Video.Ready) {
                            showReadyForPublish = false
                            break outer
                        }
                    }
                }
                if (!showReadyForPublish) { continue }
            }


            log.Printf("Emitting to topic: %s\n", publishingTopicARN)
            _, err = snsClient.Publish(context.TODO(), &sns.PublishInput{
                Message: aws.String(string(notification)),
                TopicArn: aws.String(publishingTopicARN),
            })

            if (err != nil) {
                log.Printf("Error emitting notification: %v\n", err)
                return errors.New(fmt.Sprintf("Error emitting notification: %v\n", err))
            }
        }
    }

    return nil
}

func updateMovie(movieUUID string) error {
    tableName := common.MovieTableName
    key := map[string]types.AttributeValue{
        "id": &types.AttributeValueMemberS{
            Value: movieUUID,
        },
    }

    updateExpression := "SET video.ready = :newValue"
    expressionAttributeValues := map[string]types.AttributeValue{
        ":newValue": &types.AttributeValueMemberBOOL{
            Value: true,
        },
    }

    input := &dynamodb.UpdateItemInput{
        TableName:                 aws.String(tableName),
        Key:                       key,
        UpdateExpression:          aws.String(updateExpression),
        ExpressionAttributeValues: expressionAttributeValues,
        ReturnValues:              types.ReturnValueNone,
    }

    _, err := dynamodbClient.UpdateItem(context.TODO(), input)
    if err != nil { return errors.New(fmt.Sprintf("Error updating movie %s: %v", movieUUID, err)) }

    return nil
}

func updateTvShow(showUUID string, seasonStr string, episodeStr string) (*common.Show, error) {
    season, err := strconv.Atoi(seasonStr)
    if (err != nil) { return nil, errors.New(fmt.Sprintf("cant parse season when updating video: %v", err)) }
    episode, err := strconv.Atoi(episodeStr)
    if (err != nil) { return nil, errors.New(fmt.Sprintf("cant parse episode when updating video: %v", err)) }

    // get show
    tableName := common.ShowTableName
    input := &dynamodb.GetItemInput{
        TableName: &tableName,
        Key: map[string]types.AttributeValue {
            "id": &types.AttributeValueMemberS {
                Value: showUUID,
            },
        },
    }

    result, err := dynamodbClient.GetItem(context.TODO(), input)
    if err != nil { return nil, errors.New(fmt.Sprintf("Error getting show: %v", err)) }
    if result.Item == nil { return nil, errors.New("Show does not exist") }

    var show common.Show
    err = attributevalue.UnmarshalMap(result.Item, &show)
    if err != nil { return nil, errors.New(fmt.Sprintf("Error unmarshaling show: %v", err)) }

    // find actual season and episode index in the item
    seasonActualIndex := -1
    episodeActualIndex := -1
    outer: for seasonIndex, seasonStruct := range show.Seasons {
        if (seasonStruct.SeasonNumber != uint64(season)) { continue }
        for episodeIndex, episodeStruct := range seasonStruct.Episodes {
            if (episodeStruct.EpisodeNumber == uint64(episode)) {
                seasonActualIndex = seasonIndex
                episodeActualIndex = episodeIndex
                break outer
            }
        }
    }


    // update the show
    show.Seasons[seasonActualIndex].Episodes[episodeActualIndex].Video.Ready = true
    key := map[string]types.AttributeValue{
        "id": &types.AttributeValueMemberS{
            Value: showUUID,
        },
    }

    updateExpression := fmt.Sprintf("SET seasons[%d].episodes[%d].video.ready = :newValue", seasonActualIndex, episodeActualIndex)
    expressionAttributeValues := map[string]types.AttributeValue{
        ":newValue": &types.AttributeValueMemberBOOL{
            Value: true,
        },
    }

    inputUpdate := &dynamodb.UpdateItemInput{
        TableName:                 aws.String(tableName),
        Key:                       key,
        UpdateExpression:          aws.String(updateExpression),
        ExpressionAttributeValues: expressionAttributeValues,
        ReturnValues:              types.ReturnValueNone,
    }

    _, err = dynamodbClient.UpdateItem(context.TODO(), inputUpdate)
    if err != nil { return nil, errors.New(fmt.Sprintf("Error updating show %s: %v", showUUID, err)) }

    return &show, nil
}

func downloadFile(ctx context.Context, bucket, key, filepath string) error {
    file, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = downloader.Download(ctx, file, &s3.GetObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
    })
    return err
}

func uploadFile(ctx context.Context, bucket, key, filepath string) error {
    file, err := os.Open(filepath)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = uploader.Upload(ctx, &s3.PutObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
        Body:   file,
    })
    return err
}

func transcodeVideo(inputFile string, outputFile string, scale string) error {
    cmd := exec.Command("/opt/bin/ffmpeg",
    "-i", inputFile,
    "-vf", scale,
    "-c:v", "libx264",
    "-preset", "medium",
    "-crf", "23",
    "-c:a", "aac",
    "-b:a", "128k",
    "-movflags", "+faststart",
    outputFile)

    stderr, err := cmd.StderrPipe()
    if err != nil {
        return err
    }

    err = cmd.Start();
    if (err != nil) {
        return err
    }

    stderrMsg, _ := io.ReadAll(stderr)
    log.Printf("FFmpeg stderr: %s\n", string(stderrMsg))

    return cmd.Wait()
}

func main() {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }

    s3Client = s3.NewFromConfig(cfg)
    uploader = manager.NewUploader(s3Client)
    downloader = manager.NewDownloader(s3Client)
    dynamodbClient = dynamodb.NewFromConfig(cfg)
    snsClient = sns.NewFromConfig(cfg)

    publishingTopicARN = os.Getenv("PUBLISHING_TOPIC_ARN")

    lambda.Start(handler)
}
