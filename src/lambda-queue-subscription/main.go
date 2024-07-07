package main

import (
	"common"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
		return common.ErrorResponse(http.StatusBadRequest, "Malformed input."), err
	}
	if !subscription.IsValid() {
		return common.ErrorResponse(http.StatusBadRequest, "Malformed input."), nil
	}
	subscription.UserUUIDType = fmt.Sprintf("%s#%d", subscription.UserUUID, subscription.Type)

	subscriptionTableName := common.SubscriptionTableName
	queryInput := &dynamodb.QueryInput{
		TableName: &subscriptionTableName,
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userIdType": &types.AttributeValueMemberS{Value: subscription.UserUUIDType},
			":target":     &types.AttributeValueMemberS{Value: subscription.Target},
		},
		IndexName:              jsii.String(common.SubscriptionTableSecondaryIndex),
		KeyConditionExpression: aws.String("userIdType = :userIdType and target = :target"),
	}
	queryOutput, err := dynamoDbClient.Query(context.TODO(), queryInput)
	if err != nil {
		log.Printf("Error querying subscriptions: %v", err)
		return common.ErrorResponse(http.StatusInternalServerError, "Error querying subscriptions."), err
	}
	if queryOutput.Count != 0 {
		return common.ErrorResponse(http.StatusBadRequest, "Subscription already exists."), nil
	}

	message, err := json.Marshal(subscription)
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, "Error marshalling subscription."), err
	}

	queueUrl, err := sqsClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: jsii.String(common.SubscriptionQueueName),
	})
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, "Error getting queue URL."), err
	}

	input := &sqs.SendMessageInput{
		MessageBody: jsii.String(string(message)),
		QueueUrl:    queueUrl.QueueUrl,
	}
	_, err = sqsClient.SendMessage(context.TODO(), input)
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, "Error sending message to queue."), err
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