package main

import (
	"common"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/jsii-runtime-go"
	"log"
	"net/http"
)

var dynamoDbClient *dynamodb.Client
var sqsClient *sqs.Client

func queueUnsubscription(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	uuid, ok := request.QueryStringParameters["uuid"]
	if !ok {
		return common.ErrorResponse(http.StatusBadRequest, "Malformed input."), nil
	}
	if len(uuid) == 0 {
		return common.ErrorResponse(http.StatusBadRequest, "Malformed input."), nil
	}

	tableName := common.SubscriptionTableName
	getInput := &dynamodb.GetItemInput{
		Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: uuid}},
		TableName: &tableName,
	}

	getOutput, err := dynamoDbClient.GetItem(context.TODO(), getInput)
	if err != nil {
		log.Printf("Error getting subscription: %v", err)
		return common.ErrorResponse(http.StatusInternalServerError, "Error getting subscription."), err
	}
	if getOutput.Item == nil {
		return common.ErrorResponse(http.StatusNotFound, "Subscription not found."), nil
	}

	var subscription common.Subscription
	err = attributevalue.UnmarshalMap(getOutput.Item, &subscription)
	if err != nil {
		log.Printf("Error unmarshalling subscription: %v", err)
		return common.ErrorResponse(http.StatusInternalServerError, "Error unmarshalling subscription."), err
	}

	queueUrl, err := sqsClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: jsii.String(common.UnsubscriptionQueueName),
	})
	if err != nil {
		log.Printf("Error getting queue url: %v", err)
		return common.ErrorResponse(http.StatusInternalServerError, "Error getting queue url."), err
	}

	sendInput := &sqs.SendMessageInput{
		MessageBody: jsii.String(uuid),
		QueueUrl:    queueUrl.QueueUrl,
	}
	_, err = sqsClient.SendMessage(context.TODO(), sendInput)
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, "Error sending message to queue."), err
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
}

func main() {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
	if err != nil {
		log.Fatal("Cannot load in default config.")
	}

	dynamoDbClient = dynamodb.NewFromConfig(sdkConfig)
	sqsClient = sqs.NewFromConfig(sdkConfig)

	lambda.Start(queueUnsubscription)
}
