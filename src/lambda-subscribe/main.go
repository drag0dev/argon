package main

import (
	"common"
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/jsii-runtime-go"
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

		// In case a queued subscription that would be the result of this query got written to the db while this
		// subscription was in the queue
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
			log.Printf("Error querying subscriptions: %v\n", err)
			return err
		}
		if queryOutput.Count != 0 {
			log.Printf("Subscription already exists: %v\n", subscription)
			return errors.New("subscription already exists")
		}

		subscription.UUID = uuid.New().String()

		marshaledSubscription, err := attributevalue.MarshalMap(subscription)
		if err != nil {
			log.Printf("Error marshalling subscription: %v\n", err)
			return err
		}

		putInput := &dynamodb.PutItemInput{
			TableName: &subscriptionTableName,
			Item:      marshaledSubscription,
		}

		_, err = dynamoDbClient.PutItem(context.TODO(), putInput)
		if err != nil {
			log.Printf("Error putting subscription: %v", err)
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

	lambda.Start(uploadSubscription)
}
