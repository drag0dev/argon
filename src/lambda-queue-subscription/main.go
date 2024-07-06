package main

import (
	"common"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/jsii-runtime-go"
	"log"
	"net/http"
)

var dynamoDbClient *dynamodb.Client
var sqsClient *sqs.Client

func queueSubscription(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	var subscription common.Subscription
	err := json.Unmarshal([]byte(request.Body), &subscription)
	if err != nil {
		return common.ErrorResponse(http.StatusBadRequest, "Malformed input."), nil
	}
	if !subscription.IsValid() {
		return common.ErrorResponse(http.StatusBadRequest, "Malformed input."), nil
	}

	message, err := json.Marshal(subscription)
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, "Error marshalling subscription."), nil
	}

	queueUrl, err := sqsClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: jsii.String(common.SubscriptionQueueName),
	})
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, "Error getting queue URL."), nil
	}

	input := &sqs.SendMessageInput{
		MessageBody: jsii.String(string(message)),
		QueueUrl:    queueUrl.QueueUrl,
	}
	_, err = sqsClient.SendMessage(context.TODO(), input)
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, "Error sending message to queue."), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(message),
	}, nil
}

func main() {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
	if err != nil {
		log.Fatal("Cannot load in default config.")
	}

	dynamoDbClient = dynamodb.NewFromConfig(sdkConfig)
	sqsClient = sqs.NewFromConfig(sdkConfig)

	lambda.Start(queueSubscription)
}
