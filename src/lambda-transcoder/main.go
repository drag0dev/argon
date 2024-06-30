package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
    "common"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)


var s3Client *s3.Client
var uploader *manager.Uploader
var downloader *manager.Downloader

func handler(ctx context.Context, s3Event events.S3Event) error {
    for _, record := range s3Event.Records {
        bucket := record.S3.Bucket.Name
        key := record.S3.Object.Key
        log.Printf("Processing %s - %s", bucket, key)

        // name of the file without the suffix and the format of the file
        name, _ := strings.CutSuffix(key, common.OriginalSuffix)

        inputFile := "/tmp/input_video"
        outputFile := "/tmp/output.mp4"

        err := downloadFile(ctx, bucket, key, inputFile)
        if (err != nil) {
            return fmt.Errorf("failed to download file: %v", err)
        }

        resolutions:= []string{common.Resolution1, common.Resolution2, common.Resolution3}
        for _, res := range resolutions {
            // Transcode the video
            err = transcodeVideo(inputFile, outputFile, fmt.Sprintf("scale=%s", res))
            if (err != nil) {
                return fmt.Errorf("failed to transcode video: %v", err)
            }

            // base_name_of_the_file/resolution.mp4
            outputKey := fmt.Sprintf("%s/%s.mp4", name, res)
            err = uploadFile(ctx, bucket, outputKey, outputFile)
            if (err != nil) {
                return fmt.Errorf("failed to upload file: %v", err)
            }
        }
        // Clean up temporary files
        os.Remove(inputFile)
        os.Remove(outputFile)
    }

    return nil
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

    lambda.Start(handler)
}
