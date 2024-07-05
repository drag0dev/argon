package main

import (
	"common"
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
)

var dynamoDbClient *dynamodb.Client

func uploadSubscription(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	var subscription common.Subscription
	err := json.Unmarshal([]byte(request.Body), &subscription)
	if err != nil {
		return common.ErrorResponse(http.StatusBadRequest, "Malformed input."), nil
	}

	// TODO: Validate subscription

	subscription.UUID = uuid.New().String()

	marshaledSubscription, err := attributevalue.MarshalMap(subscription)
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, "Error marshaling subscription."), nil
	}

	subscriptionTableName := common.SubscriptionTableName
	input := &dynamodb.PutItemInput{
		TableName: &subscriptionTableName,
		Item:      marshaledSubscription,
	}

	_, err = dynamoDbClient.PutItem(context.TODO(), input)
	if err != nil {
		log.Printf("Error putting subscription: %v", err)
		return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("error putting subscription")
	}

	result, err := json.Marshal(subscription)
	if err != nil {
		log.Printf("Error marshaling subscription: %v", err)
		return common.EmptyErrorResponse(http.StatusInternalServerError),
			errors.New("error marshaling subscription")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"content-type": "application/json"},
		Body:       string(result),
	}, nil
}

func main() {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
	if err != nil {
		log.Fatal("Cannot load in default config.")
	}

	dynamoDbClient = dynamodb.NewFromConfig(sdkConfig)

	lambda.Start(uploadSubscription)
}
