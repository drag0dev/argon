package main

import (
	"common"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/lestrrat-go/jwx/jwt"
)

var dynamoDbClient *dynamodb.Client
var sqsClient *sqs.Client
var preferenceChangeQueueClient *sqs.Client

func enqueChangePreferenceItem(prefChangeItem common.PreferenceChange, headerVal string) {
    token := strings.TrimPrefix(headerVal, "Bearer ")
    if (token == "") {
        log.Printf("Missing token")
        return
    }

    parsedToken, err := jwt.Parse([]byte(token))
    if (err != nil) {
        log.Printf("Error parsing token: %v", err)
        return
    }

    sub, ok := parsedToken.Get("sub")
    if !ok {
        log.Println("sub claim not found in token")
        return
    }

    userId, ok := sub.(string)
    if !ok {
        log.Println("userid is not string")
        return
    }
    prefChangeItem.UserId = userId

    prefChangeMarshaled, err := json.Marshal(prefChangeItem)
    if (err != nil ) {
        log.Printf("Error marshaling preference change item: %v", err)
        return
    }

    queueUrl, err := preferenceChangeQueueClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
        QueueName: aws.String(common.PreferenceUpdateQueue),
    })
    if err != nil {
        log.Printf("Error getting queue url: %v", err)
        return
    }

    sendInput := &sqs.SendMessageInput{
        MessageBody: aws.String(string(prefChangeMarshaled)),
        QueueUrl:    queueUrl.QueueUrl,
    }
    _, err = preferenceChangeQueueClient.SendMessage(context.TODO(), sendInput)
    if err != nil { log.Printf("Error enquing preference change item: %v", err) }
}

func queueReview(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var review common.Review
	err := json.Unmarshal([]byte(request.Body), &review)
	if err != nil {
		return common.ErrorResponse(http.StatusBadRequest, "Malformed input."), err
	}
	if !review.IsValid() {
		return common.ErrorResponse(http.StatusBadRequest, "Malformed input."), nil
	}

    actors := []string{}
    directors := []string{}
    genres := []string{}

	movieTableName := common.MovieTableName
	getMovieResult, err := dynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &movieTableName,
		Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: review.TargetUUID}},
	})
	if err != nil {
		log.Printf("Error getting review target: %v", err)
		return common.ErrorResponse(http.StatusInternalServerError, "Error getting review target."), err
	}
	if getMovieResult.Item == nil {
		showTableName := common.ShowTableName
		getShowResult, err := dynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
			TableName: &showTableName,
			Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: review.TargetUUID}},
		})
		if err != nil {
			log.Printf("Error getting review target: %v", err)
			return common.ErrorResponse(http.StatusInternalServerError, "Error getting review target."), err
		}
		if getShowResult.Item == nil {
			return common.ErrorResponse(http.StatusNotFound, "Review target not found."), nil
		}
        var show common.Show
        err = attributevalue.UnmarshalMap(getShowResult.Item, &show)
        if err != nil {
            log.Printf("Error unmarshaling movie :%v", err)
            return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error umarshaling movie")
        }
        actors = show.Actors
        directors = show.Directors
        genres = show.Genres
    } else {
        // we have a movie
        var movie common.Movie
        err = attributevalue.UnmarshalMap(getMovieResult.Item, &movie)
        if err != nil {
            log.Printf("Error unmarshaling movie :%v", err)
            return common.EmptyErrorResponse(http.StatusInternalServerError), errors.New("Error umarshaling movie")
        }
        actors = movie.Actors
        directors = movie.Directors
        genres = movie.Genres
    }

    changeWeight := 0
    if (review.Grade == 1) {
        changeWeight = common.ReviewChangeWeight1
    } else if (review.Grade == 2) {
        changeWeight = common.ReviewChangeWeight2
    } else if (review.Grade == 3) {
        changeWeight = common.ReviewChangeWeight3
    } else if (review.Grade == 4) {
        changeWeight = common.ReviewChangeWeight4
    } else if (review.Grade == 5) {
        changeWeight = common.ReviewChangeWeight5
    }
    prefChangeItem := common.PreferenceChange{
        UpdateWeight: common.ReviewUpdateWeight,
        ChangeWeight: float64(changeWeight),
        Actors: actors,
        Directors: directors,
        Genres: genres,
    }
    enqueChangePreferenceItem(prefChangeItem, request.Headers["Authorization"])

	message, err := json.Marshal(review)
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, "Error marshalling review."), err
	}

	queueUrl, err := sqsClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(common.ReviewQueueName),
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

func main() {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}

	dynamoDbClient = dynamodb.NewFromConfig(sdkConfig)
	sqsClient = sqs.NewFromConfig(sdkConfig)
    preferenceChangeQueueClient = sqs.NewFromConfig(sdkConfig)
	lambda.Start(queueReview)
}
