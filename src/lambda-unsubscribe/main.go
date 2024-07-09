package main

import (
	"common"
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
)

var dynamoDbClient *dynamodb.Client

func deleteSubscription(
	ctx context.Context,
	sqsEvent events.SQSEvent,
) error {
	var err error
	for _, message := range sqsEvent.Records {
		if len(message.Body) == 0 {
			return errors.New("received empty id from sqs")
		}

		deleteInput := &dynamodb.DeleteItemInput{
			TableName: aws.String(common.SubscriptionTableName),
			Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: message.Body}},
		}

		_, err = dynamoDbClient.DeleteItem(context.TODO(), deleteInput)
		if err != nil {
			log.Printf("Error deleting subscription: %v", err)
			return err
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

	lambda.Start(deleteSubscription)
}
