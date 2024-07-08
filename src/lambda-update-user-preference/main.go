package main

import (
	"common"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var dynamodbClient *dynamodb.Client

type Wrapper struct {
    Message string `json:"message"`
}

type UserPreferences struct {
    Actors map[string]float64
    Directors map[string]float64
    Genres map[string]float64
    UpdateWeight float64
}

func handlePublish(ctx context.Context, sqsEvent events.SQSEvent) error {
    // parse all items
    var parsedItems []common.PreferenceChange
    for _, message := range sqsEvent.Records {

        log.Println(message.Body)

        var tempItem common.PreferenceChange
        err := json.Unmarshal([]byte(message.Body), &tempItem)
        if (err != nil)  { log.Printf("Error unmarshaling into item: %v\n", err) }

        parsedItems = append(parsedItems, tempItem)
    }

    var users map[string]UserPreferences = make(map[string]UserPreferences)
    for _, item := range parsedItems {
        value, ok := users[item.UserId]
        if (!ok) {
            users[item.UserId] = UserPreferences{
                Actors: make(map[string]float64),
                Directors: make(map[string]float64),
                Genres: make(map[string]float64),
            }
        }
        value, _ = users[item.UserId]

        // apply genres
        for _, genre := range item.Genres {
            value.Genres[genre] += item.ChangeWeight
        }
        // apply directors
        for _, director := range item.Directors {
            value.Directors[director] += item.ChangeWeight
        }
        // apply actors
        for _, actor := range item.Actors {
            value.Actors[actor] += item.ChangeWeight
        }

        // apply update weight
        value.UpdateWeight += item.UpdateWeight
        users[item.UserId] = value
    }

    for userId, pref := range users {
        var attributeNames map[string]string = make (map[string]string)
        var attributeValues map[string]types.AttributeValue = make(map[string]types.AttributeValue)

        // init fields if they dont exist
        updateExpression := `
        SET #counterName = if_not_exists(#counterName, :zero) + :counterValue,
             #actorsName = if_not_exists(#actorsName, :empty_map),
             #directorsName = if_not_exists(#directorsName, :empty_map),
             #genresName = if_not_exists(#genresName, :empty_map)
            `
        attributeNames["#counterName"] = "updateCounter"
        attributeNames["#actorsName"] = "actors"
        attributeNames["#directorsName"] = "directors"
        attributeNames["#genresName"] = "genres"
        attributeValues[":zero"] = &types.AttributeValueMemberN{Value: "0"}
        attributeValues[":empty_map"] = &types.AttributeValueMemberM{Value: make(map[string]types.AttributeValue)}
        attributeValues[":counterValue"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", pref.UpdateWeight)}
        input := &dynamodb.UpdateItemInput{
            TableName: aws.String(common.UserPreferenceTableName),
            Key: map[string]types.AttributeValue{"userId": &types.AttributeValueMemberS{Value: userId}},
            UpdateExpression: aws.String(updateExpression),
            ExpressionAttributeNames: attributeNames,
            ExpressionAttributeValues: attributeValues,
            ReturnValues: types.ReturnValueNone,
        }
        _, errr := dynamodbClient.UpdateItem(ctx, input)
        if (errr != nil) {
            log.Printf("Error init fields : %v\n", errr)
            return errr
        }



        attributeNames = make (map[string]string)
        attributeValues = make(map[string]types.AttributeValue)
        attributeValues[":zero"] = &types.AttributeValueMemberN{Value: "0"}
        updateExpression = "SET "

        idx := 0
        for actor, val := range pref.Actors {
            updateExpression += fmt.Sprintf(" actors.#key%d = if_not_exists(actors.#key%d, :zero) + :value%d,", idx, idx, idx)
            attributeNames[fmt.Sprintf("#key%d", idx)] = actor
            attributeValues[fmt.Sprintf(":value%d", idx)] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", val)}
            idx += 1
        }
        for director, val := range pref.Directors {
            updateExpression += fmt.Sprintf(" directors.#key%d = if_not_exists(directors.#key%d, :zero) + :value%d,", idx, idx, idx)
            attributeNames[fmt.Sprintf("#key%d", idx)] = director
            attributeValues[fmt.Sprintf(":value%d", idx)] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", val)}
            idx += 1
        }
        for genre, val := range pref.Genres {
            updateExpression += fmt.Sprintf(" genres.#key%d = if_not_exists(genres.#key%d, :zero) + :value%d,", idx, idx, idx)
            attributeNames[fmt.Sprintf("#key%d", idx)] = genre
            attributeValues[fmt.Sprintf(":value%d", idx)] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", val)}
            idx += 1
        }

        updateExpression, _ = strings.CutSuffix(updateExpression, ",")

        log.Println(updateExpression)

        input = &dynamodb.UpdateItemInput{
            TableName: aws.String(common.UserPreferenceTableName),
            Key: map[string]types.AttributeValue{"userId": &types.AttributeValueMemberS{Value: userId}},
            UpdateExpression: aws.String(updateExpression),
            ExpressionAttributeNames: attributeNames,
            ExpressionAttributeValues: attributeValues,
            ReturnValues: types.ReturnValueAllNew,
        }

        var updatedPreference *dynamodb.UpdateItemOutput = nil
        var err error
        for i := 0; i < 5; i++ {
            updatedPreference, err = dynamodbClient.UpdateItem(ctx, input)
            var throughputErr *types.ProvisionedThroughputExceededException
            if errors.As(err, &throughputErr) {
                time.Sleep(time.Duration(100000000))
                continue
            } else if err != nil {
                log.Printf("Error updating user preference: %v\n", err)
            }
        }

        if (updatedPreference == nil) {
            return errors.New("Could not update user preference")
        }

        updatedCounter, ok := updatedPreference.Attributes["updateCounter"]
        if (!ok) { log.Printf("Missing updated update counter in the updated preference!") }
        updatedCounterCast, ok := updatedCounter.(*types.AttributeValueMemberN)
        if (!ok) { log.Printf("Updated update counter not a number type") }
        updatedCounterVal, err := strconv.Atoi(updatedCounterCast.Value)
        if (err != nil) { log.Printf("Updated update counter not a number") }

        if (updatedCounterVal > common.UpdateCounterThreshhold) {
            log.Printf("New feed")
            updateExpression = "SET #counterName = :counterValue"
            input := &dynamodb.UpdateItemInput{
                TableName: aws.String(common.UserPreferenceTableName),
                Key: map[string]types.AttributeValue{"userId": &types.AttributeValueMemberS{Value: userId}},
                UpdateExpression: aws.String(updateExpression),
                ExpressionAttributeNames: map[string]string{ "#counterName": "updateCounter"},
                ExpressionAttributeValues: map[string]types.AttributeValue{
                    ":counterValue": &types.AttributeValueMemberN{Value: "0"},
                },
                ReturnValues: types.ReturnValueNone,
            }
            _, err := dynamodbClient.UpdateItem(ctx, input)
            if err != nil {
                log.Printf("Error updating update counter: %v", err)
            }
        }
    }

    return nil
}

func main() {
    sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
    if err != nil {
        log.Fatal("Cannot load in default config")
    }
    dynamodbClient = dynamodb.NewFromConfig(sdkConfig)
    lambda.Start(handlePublish)
}
