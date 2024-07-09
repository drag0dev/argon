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
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"log"
	"net/http"
)

var dynamoDbClient *dynamodb.Client
var sqsClient *sqs.Client

func queueMessage(editMetadataRequest *common.EditMetadataRequest) (events.APIGatewayProxyResponse, error) {
	message, err := json.Marshal(editMetadataRequest)
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, "Error marshalling edit request."), err
	}

	queueUrl, err := sqsClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(common.EditMetadataRequestQueueName),
	})
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, "Error getting queue URL."), err
	}

	_, err = sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: aws.String(string(message)),
		QueueUrl:    queueUrl.QueueUrl,
	})
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, "Error sending message to queue."), err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS",
			"Access-Control-Allow-Headers":     "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		},
		Body: string(message),
	}, nil
}

func handleMovie(editMetadataRequest *common.EditMetadataRequest) (events.APIGatewayProxyResponse, error) {
	movieTableName := common.MovieTableName
	getMovieResult, err := dynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &movieTableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: editMetadataRequest.TargetUUID},
		},
	})
	if err != nil {
		log.Printf("Error getting edit target: %v", err)
		return common.ErrorResponse(http.StatusInternalServerError, "Error getting edit target"), err
	}
	if getMovieResult.Item == nil {
		return common.ErrorResponse(http.StatusNotFound, "Edit target not found"), nil
	}

	return queueMessage(editMetadataRequest)
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

func handleShow(editMetadataRequest *common.EditMetadataRequest) (events.APIGatewayProxyResponse, error) {
	getShowResult, err := getShow(editMetadataRequest.TargetUUID)
	if err != nil {
		log.Printf("Error getting edit target: %v", err)
		return common.ErrorResponse(http.StatusInternalServerError, "Error getting edit target"), err
	}
	if getShowResult.Item == nil {
		return handleMovie(editMetadataRequest)
	}

	return queueMessage(editMetadataRequest)
}

func handleEpisode(editMetadataRequest *common.EditMetadataRequest) (events.APIGatewayProxyResponse, error) {
	getShowResult, err := getShow(editMetadataRequest.TargetUUID)
	if err != nil {
		log.Printf("Error getting edit target: %v", err)
		return common.ErrorResponse(http.StatusInternalServerError, "Error getting edit target"), err
	}
	if getShowResult.Item == nil {
		return common.ErrorResponse(http.StatusNotFound, "Edit target not found"), nil
	}

	var show common.Show
	err = attributevalue.UnmarshalMap(getShowResult.Item, &show)
	if err != nil {
		log.Printf("Error unmarshalling edit target: %v", err)
		return common.ErrorResponse(http.StatusInternalServerError, "Error unmarshalling edit target"), err
	}

	seasonActualIdx, episodeActualIdx := -1, -1
outer:
	for seasonIdx, season := range show.Seasons {
		if season.SeasonNumber != *editMetadataRequest.SeasonNumber {
			continue
		}

		seasonActualIdx = seasonIdx

		for episodeIdx, episode := range season.Episodes {
			if episode.EpisodeNumber != *editMetadataRequest.EpisodeNumber {
				continue
			}

			episodeActualIdx = episodeIdx
			break outer
		}
	}
	if seasonActualIdx == -1 || episodeActualIdx == -1 {
		return common.ErrorResponse(http.StatusNotFound, "Edit target not found"), nil
	}

	return queueMessage(editMetadataRequest)
}

func queueEditMetadata(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	var editMetadataRequest common.EditMetadataRequest
	err := json.Unmarshal([]byte(request.Body), &editMetadataRequest)
	if err != nil {
		return common.ErrorResponse(http.StatusBadRequest, "Malformed input."), err
	}
	if !editMetadataRequest.IsValid() {
		return common.ErrorResponse(http.StatusBadRequest, "Malformed input."), nil
	}

	var response events.APIGatewayProxyResponse
	if editMetadataRequest.SeasonNumber == nil {
		response, err = handleShow(&editMetadataRequest)
	} else {
		response, err = handleEpisode(&editMetadataRequest)
	}

	return response, err
}

func main() {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
	if err != nil {
		log.Fatal("Cannot load in default config.")
	}

	dynamoDbClient = dynamodb.NewFromConfig(sdkConfig)
	sqsClient = sqs.NewFromConfig(sdkConfig)

	lambda.Start(queueEditMetadata)
}
