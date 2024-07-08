package main

import (
	"common"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
)

var dynamoDbClient *dynamodb.Client

func editMovieOrShow(editMetadataRequest *common.EditMetadataRequest) error {
	marshaledGenres, err := attributevalue.MarshalList(editMetadataRequest.Genres)
	if err != nil {
		log.Printf("Error marshaling target genres: %v", err)
	}
	marshaledActors, err := attributevalue.MarshalList(editMetadataRequest.Actors)
	if err != nil {
		log.Printf("Error marshaling target actors: %v", err)
	}
	marshaledDirectors, err := attributevalue.MarshalList(editMetadataRequest.Directors)
	if err != nil {
		log.Printf("Error marshaling target directors: %v", err)
	}

	_, err = dynamoDbClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(common.ShowTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: editMetadataRequest.TargetUUID},
		},
		UpdateExpression: aws.String("SET title = :title, description = :description, genres = :genres, " +
			"actors = :actors, directors = :directors"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":title":       &types.AttributeValueMemberS{Value: editMetadataRequest.Title},
			":description": &types.AttributeValueMemberS{Value: editMetadataRequest.Description},
			":genres":      &types.AttributeValueMemberL{Value: marshaledGenres},
			":actors":      &types.AttributeValueMemberL{Value: marshaledActors},
			":directors":   &types.AttributeValueMemberL{Value: marshaledDirectors},
		},
	})
	if err != nil {
		log.Printf("Error updating target show seasons: %v", err)
		return err
	}

	return nil
}

func handleMovie(editMetadataRequest *common.EditMetadataRequest) error {
	movieTableName := common.MovieTableName
	getMovieResult, err := dynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &movieTableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: editMetadataRequest.TargetUUID},
		},
	})
	if err != nil {
		log.Printf("Error getting edit target: %v", err)
		return err
	}
	if getMovieResult.Item == nil {
		log.Println("Edit target not found.")
		return nil
	}

	err = editMovieOrShow(editMetadataRequest)
	if err != nil {
		return err
	}

	return nil
}

func getShow(uuid string) (*dynamodb.GetItemOutput, error) {
	showTableName := common.ShowTableName
	return dynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &showTableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: uuid},
		},
	})
}

func handleShow(editMetadataRequest *common.EditMetadataRequest) error {
	getShowResult, err := getShow(editMetadataRequest.TargetUUID)
	if err != nil {
		log.Printf("Error getting edit target: %v", err)
		return err
	}
	if getShowResult.Item == nil {
		return handleMovie(editMetadataRequest)
	}

	err = editMovieOrShow(editMetadataRequest)
	if err != nil {
		return err
	}

	return nil
}

func handleEpisode(editMetadataRequest *common.EditMetadataRequest) error {
	getShowResult, err := getShow(editMetadataRequest.TargetUUID)
	if err != nil {
		log.Printf("Error getting edit target: %v", err)
		return err
	}
	if getShowResult.Item == nil {
		return nil
	}

	var show common.Show
	err = attributevalue.UnmarshalMap(getShowResult.Item, &show)
	if err != nil {
		log.Printf("Error unmarshalling edit target: %v", err)
		return err
	}

	if uint64(len(show.Seasons)) <= *editMetadataRequest.SeasonNumber {
		return nil
	}
	targetSeason := &show.Seasons[*editMetadataRequest.SeasonNumber]
	if uint64(len(targetSeason.Episodes)) <= *editMetadataRequest.EpisodeNumber {
		return nil
	}

	targetEpisode := &targetSeason.Episodes[*editMetadataRequest.EpisodeNumber]
	targetEpisode.Title = editMetadataRequest.Title
	targetEpisode.Description = editMetadataRequest.Description
	targetEpisode.Actors = editMetadataRequest.Actors
	targetEpisode.Directors = editMetadataRequest.Directors

	marshaledSeasons, err := attributevalue.MarshalList(show.Seasons)
	if err != nil {
		log.Printf("Error marshalling target show seasons: %v", err)
		return err
	}

	_, err = dynamoDbClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(common.ShowTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: editMetadataRequest.TargetUUID},
		},
		UpdateExpression: aws.String("SET seasons = :seasons"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":seasons": &types.AttributeValueMemberL{Value: marshaledSeasons},
		},
	})
	if err != nil {
		log.Printf("Error updating target show seasons: %v", err)
		return err
	}

	return nil
}

func editMetadata(ctx context.Context, sqsEvent events.SQSEvent) error {
	var editMetadataRequest common.EditMetadataRequest
	var err error
	for _, message := range sqsEvent.Records {
		err = json.Unmarshal([]byte(message.Body), &editMetadataRequest)
		if err != nil {
			log.Printf("Error unmarshalling sqs message: %v", err)
			return err
		}

		if editMetadataRequest.SeasonNumber == nil {
			err = handleShow(&editMetadataRequest)
		} else {
			err = handleEpisode(&editMetadataRequest)
		}
		if err != nil {
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

	lambda.Start(editMetadata)
}
