package main

import (
	"common"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"log"
)

var dynamoDbClient *dynamodb.Client

func createReview(
	ctx context.Context,
	sqsEvent events.SQSEvent,
) error {
	var review common.Review
	var err error
	for _, message := range sqsEvent.Records {
		err = json.Unmarshal([]byte(message.Body), &review)
		if err != nil {
			log.Printf("Error unmarshalling sqs message: %v\n", err)
			return err
		}

		movieTableName := common.MovieTableName
		getMovieResult, err := dynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
			TableName: &movieTableName,
			Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: review.TargetUUID}},
		})
		if err != nil {
			log.Printf("Error getting review target: %v", err)
			return err
		}
		if getMovieResult.Item == nil {
			showTableName := common.ShowTableName
			getShowResult, err := dynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
				TableName: &showTableName,
				Key: map[string]types.AttributeValue{
					"id": &types.AttributeValueMemberS{Value: review.TargetUUID},
				},
			})
			if err != nil {
				log.Printf("Error getting review target: %v", err)
				return err
			}
			if getShowResult.Item == nil {
				return err
			}
		}

		review.UUID = uuid.New().String()

		marshaledReview, err := attributevalue.MarshalMap(review)
		if err != nil {
			log.Printf("Error marshalling review: %v", err)
			return err
		}

		reviewTableName := common.ReviewTableName
		_, err = dynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: &reviewTableName,
			Item:      marshaledReview,
		})
		if err != nil {
			log.Printf("Error putting review: %v", err)
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

	lambda.Start(createReview)
}
