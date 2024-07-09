package services

import (
	"common"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	common2 "github.com/drag0dev/argon/src/ecs/pkg/common"
	"log"
)

func GetSubscriptions(userId string) ([]common.Subscription, error) {
	dynamoDbClient, err := common2.GetDynamoDbClient()
	if err != nil {
		log.Fatal("Cannot load in default config.")
		return nil, err
	}

	queryOutput, err := dynamoDbClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName: aws.String(common.SubscriptionTableName),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberS{Value: userId},
		},
		IndexName:              aws.String(common.SubscriptionTableUserIdSecondaryIndex),
		KeyConditionExpression: aws.String("userId = :userId"),
	})
	if err != nil {
		log.Printf("Error querying subscriptions: %v", err)
		return nil, err
	}
	if queryOutput.Count == 0 {
		return make([]common.Subscription, 0), nil
	}

	var subscriptions []common.Subscription
	err = attributevalue.UnmarshalListOfMaps(queryOutput.Items, &subscriptions)
	if err != nil {
		log.Printf("Error unmarshaling subscriptions: %v", err)
		return nil, err
	}

	return subscriptions, nil
}
