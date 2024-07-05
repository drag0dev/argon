package main

import (
	"common"
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"log"
)

var dynamoDbClient *dynamodb.Client

func uploadSubscription(
	ctx context.Context,
	sqsEvent events.SQSEvent,
) error {
	var subscription common.Subscription
	var err error
	for _, message := range sqsEvent.Records {
		err = json.Unmarshal([]byte(message.Body), &subscription)
		if err != nil {
			log.Printf("Error unmarshalling sqs message: %v\n", err)
			return err
		}

		subscription.UUID = uuid.New().String()

		marshaledSubscription, err := attributevalue.MarshalMap(subscription)
		if err != nil {
			log.Printf("Error marshalling subscription: %v\n", err)
			return err
		}

		subscriptionTableName := common.SubscriptionTableName
		input := &dynamodb.PutItemInput{
			TableName: &subscriptionTableName,
			Item:      marshaledSubscription,
		}

		_, err = dynamoDbClient.PutItem(context.TODO(), input)
		if err != nil {
			log.Printf("Error putting subscription: %v", err)
			return errors.New("error putting subscription")
		}
	}

	return nil
}

func main() {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
	if err != nil {
		log.Fatal("Cannot load in default config.")
	}

	dynamoDbClient = dynamodb.NewFromConfig(sdkConfig)

	lambda.Start(uploadSubscription)
}
