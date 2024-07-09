package services

import (
	"common"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	common2 "github.com/drag0dev/argon/src/ecs/pkg/common"
	"log"
)

func GetMovies() ([]common.Movie, error) {
	dynamoDbClient, err := common2.GetDynamoDbClient()
	if err != nil {
		log.Fatal("Cannot load in default config.")
		return nil, err
	}

	scanOutput, err := dynamoDbClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(common.MovieTableName),
	})
	if err != nil {
		log.Printf("Error querying movies: %v", err)
		return nil, err
	}
	if scanOutput.Count == 0 {
		return make([]common.Movie, 0), nil
	}

	var movies []common.Movie
	err = attributevalue.UnmarshalListOfMaps(scanOutput.Items, &movies)
	if err != nil {
		log.Printf("Error unmarshaling movies: %v", err)
		return nil, err
	}

	return movies, nil
}

func GetShows() ([]common.Show, error) {
	dynamoDbClient, err := common2.GetDynamoDbClient()
	if err != nil {
		log.Fatal("Cannot load in default config.")
		return nil, err
	}

	scanOutput, err := dynamoDbClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(common.ShowTableName),
	})
	if err != nil {
		log.Printf("Error querying shows: %v", err)
		return nil, err
	}
	if scanOutput.Count == 0 {
		return make([]common.Show, 0), nil
	}

	var shows []common.Show
	err = attributevalue.UnmarshalListOfMaps(scanOutput.Items, &shows)
	if err != nil {
		log.Printf("Error unmarshaling shows: %v", err)
		return nil, err
	}

	return shows, nil
}
