package main

import (
	"common"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"log"
	"net/http"
)

var dynamoDbClient *dynamodb.Client
var sqsClient *sqs.Client

func queueReview(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var review common.Review
	err := json.Unmarshal([]byte(request.Body), &review)
	if err != nil {
		return common.ErrorResponse(http.StatusBadRequest, "Malformed input."), err
	}
	if !review.IsValid() {
		return common.ErrorResponse(http.StatusBadRequest, "Malformed input."), nil
	}

	movieTableName := common.MovieTableName
	getMovieResult, err := dynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &movieTableName,
		Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: review.TargetUUID}},
	})
	if err != nil {
		log.Printf("Error getting review target: %v", err)
		return common.ErrorResponse(http.StatusInternalServerError, "Error getting review target."), err
	}
	if getMovieResult.Item == nil {
		showTableName := common.ShowTableName
		getShowResult, err := dynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
			TableName: &showTableName,
			Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: review.TargetUUID}},
		})
		if err != nil {
			log.Printf("Error getting review target: %v", err)
			return common.ErrorResponse(http.StatusInternalServerError, "Error getting review target."), err
		}
		if getShowResult.Item == nil {
			return common.ErrorResponse(http.StatusNotFound, "Review target not found."), nil
		}
	}

	message, err := json.Marshal(review)
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, "Error marshalling review."), err
	}

	queueUrl, err := sqsClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(common.ReviewQueueName),
	})
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, "Error getting queue URL."), err
	}

	_, err = sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: aws.String(string(message)),
		QueueUrl:    queueUrl.QueueUrl,
	})
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, "Error sending message to queue."), err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS",
			"Access-Control-Allow-Headers":     "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		},
		Body: string(message),
	}, nil
}

func main() {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}

	dynamoDbClient = dynamodb.NewFromConfig(sdkConfig)
	sqsClient = sqs.NewFromConfig(sdkConfig)

	lambda.Start(queueReview)
}
