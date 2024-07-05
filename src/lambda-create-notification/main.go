package main

import (
	"common"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
        sesTypes "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

var dynamodbClient *dynamodb.Client
var sesClient *sesv2.Client
var cognitoClient *cognitoidentityprovider.Client
var cognitoUserPoolID string
var sesEmailIdentity string

func getMovie(uuid string) (*common.Movie, error) {
    tableName := common.MovieTableName
    input := &dynamodb.GetItemInput{
        TableName: &tableName,
        Key: map[string]types.AttributeValue {
            "id": &types.AttributeValueMemberS {
                Value: uuid,
            },
        },
    }

    result, err := dynamodbClient.GetItem(context.TODO(), input)
    if err != nil {
        log.Printf("Error getting movie: %v\n", err)
        return nil, errors.New(fmt.Sprintf("Error getting movie: %v", err))
    }

    if result.Item == nil {
        return nil, nil
    }

    var movie common.Movie
    err = attributevalue.UnmarshalMap(result.Item, &movie)
    if err != nil {
        log.Printf("Error unmarshaling movie: %v\n", err)
        return nil, errors.New(fmt.Sprintf("Error unmarshaling movie: %v", err))
    }

    return &movie, nil
}

func getShow(uuid string) (*common.Show, error) {
    tableName := common.ShowTableName
    input := &dynamodb.GetItemInput{
        TableName: &tableName,
        Key: map[string]types.AttributeValue {
            "id": &types.AttributeValueMemberS {
                Value: uuid,
            },
        },
    }

    result, err := dynamodbClient.GetItem(context.TODO(), input)
    if err != nil {
        log.Printf("Error getting show: %v\n", err)
        return nil, errors.New(fmt.Sprintf("Error getting show: %v", err))
    }

    if result.Item == nil {
        return nil, nil
    }

    var show common.Show
    err = attributevalue.UnmarshalMap(result.Item, &show)
    if err != nil {
        log.Printf("Error unmarshaling show: %v\n", err)
        return nil, errors.New(fmt.Sprintf("Error unmarshaling show: %v", err))
    }

    return &show, nil
}

func craftEmail(title string, description string, destinationEmail string) sesv2.SendEmailInput {
    body := fmt.Sprintf("%s just released!\n\n%s\n\nWatch now!", title, description)
    return sesv2.SendEmailInput{
        Destination: &sesTypes.Destination{
            ToAddresses: []string{destinationEmail},
        },
        Content: &sesTypes.EmailContent{
            Simple: &sesTypes.Message{
                Body: &sesTypes.Body{
                    Text: &sesTypes.Content{
                        Data: aws.String(body),
                    },
                },
                Subject: &sesTypes.Content{
                    Data: aws.String(fmt.Sprintf("%s just released!", title)),
                },
            },
        },
        FromEmailAddress: aws.String("your-verified-email@example.com"),
    }
}

func handlePublish(ctx context.Context, sqsEvent events.SQSEvent) error {
    for _, message := range sqsEvent.Records {
        var notification common.PublishNotification
        err := json.Unmarshal([]byte(message.Body), &notification)
        if err != nil {
            log.Printf("Error unmarshaling notification: %v\n", err)
            return err
        }

        var title string
        var description string
        var subscriptionStrings []string
        if (len(notification.MovieUUID) > 0) {
            movie, err := getMovie(notification.MovieUUID)
            if (err != nil) { return err }
            subscriptionStrings = append(subscriptionStrings, movie.Genres...)
            subscriptionStrings = append(subscriptionStrings, movie.Directors...)
            subscriptionStrings = append(subscriptionStrings, movie.Actors...)
            title = movie.Title
            description = movie.Description
        } else if (len(notification.ShowUUID) > 0) {
            show, err := getShow(notification.ShowUUID)
            if (err != nil) { return err }
            subscriptionStrings = append(subscriptionStrings, show.Genres...)
            subscriptionStrings = append(subscriptionStrings, show.Directors...)
            subscriptionStrings = append(subscriptionStrings, show.Actors...)
            title = show.Title
            description = show.Description
        } else {
            log.Printf("Nothing published, this is not supposed to happen!\n")
        }

        if (len(subscriptionStrings) == 0) { continue }

        // dynamodb query
        var conditions []expression.ConditionBuilder
        for _, value := range subscriptionStrings {
            conditions = append(conditions, expression.Name("value").Equal(expression.Value(value)))
        }
        var filterExpression expression.ConditionBuilder
        filterExpression = conditions[0]
        for _, condition := range conditions[1:] {
            filterExpression = filterExpression.Or(condition)
        }

        expr, err := expression.NewBuilder().WithFilter(filterExpression).Build()
        if err != nil {
            fmt.Printf("Error building query: %v\n", err)
            return err
        }


        var lastEvaluatedKey map[string]types.AttributeValue
        for {
            input := &dynamodb.QueryInput{
                TableName:                 aws.String("TODO"),
                IndexName:                 aws.String("TODO"),
                KeyConditionExpression:    expr.KeyCondition(),
                ExpressionAttributeNames:  expr.Names(),
                ExpressionAttributeValues: expr.Values(),
                Select:                    types.SelectSpecificAttributes,
                ProjectionExpression:      aws.String("TODO userid"),
                Limit:                     aws.Int32(1000),
                ExclusiveStartKey:         lastEvaluatedKey,
            }

            result, err := dynamodbClient.Query(ctx, input)
            if (err != nil) {
                log.Printf("Error running query: %v\n", err)
                return err
            }

            for _, item := range result.Items {
                if userID, ok := item["useid TODO"]; ok {
                    userInfo, err := cognitoClient.AdminGetUser(ctx, &cognitoidentityprovider.AdminGetUserInput{
                        UserPoolId: aws.String(cognitoUserPoolID),
                        Username:   aws.String(userID.(*types.AttributeValueMemberS).Value),
                    })
                    if err != nil {
                        log.Printf("Error getting user info: %v\n", err)
                        continue
                    }

                    var userEmail string
                    for _, attr := range userInfo.UserAttributes {
                        if *attr.Name == "email" {
                            userEmail = *attr.Value
                            break
                        }
                    }

                    email := craftEmail(title, description, userEmail)
                    _, err = sesClient.SendEmail(ctx, &email)

                    if err != nil {
                        log.Printf("Error sending email: %v\n", err)
                        return err
                    }
                }
                lastEvaluatedKey = result.LastEvaluatedKey
                if (lastEvaluatedKey == nil) { break }
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
    cognitoClient = cognitoidentityprovider.NewFromConfig(sdkConfig)
    sesClient = sesv2.NewFromConfig(sdkConfig)
    cognitoUserPoolID = os.Getenv("COGNITO_USER_POOL_ID")

    lambda.Start(handlePublish)
}
