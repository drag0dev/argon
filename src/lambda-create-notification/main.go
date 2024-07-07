package main

import (
	"common"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	sesTypes "github.com/aws/aws-sdk-go-v2/service/sesv2/types"

"github.com/bits-and-blooms/bloom/v3"
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
        FromEmailAddress: aws.String(common.SenderEmail),
    }
}

type Wrapper struct {
    Message string `json:"message"`
}

func handlePublish(ctx context.Context, sqsEvent events.SQSEvent) error {
    for _, message := range sqsEvent.Records {
        bloomFilter := bloom.NewWithEstimates(1000000, 0.01)
        var wrapper Wrapper
        err := json.Unmarshal([]byte(message.Body), &wrapper)
        if (err != nil)  {
            log.Printf("Error unmarshaling event: %v\n", err)
            return err
        }
        var notification string = wrapper.Message

        movie := false
        show := false
        if (strings.HasPrefix(notification, "movie~")) {
            movie = true
        } else if (strings.HasPrefix(notification, "show~")) {
            show = true
        }

        parts := strings.Split(notification, "~")
        if (len(parts) != 2) {
            log.Printf("Malformed notification: %s\n", notification)
            return errors.New("Malformed notification")
        }
        uuid := parts[1]

        log.Printf("Raw notification: %s\n", notification)

        var title string
        var description string
        var subscriptionStrings []string
        if (movie) {
            movie, err := getMovie(uuid)
            if (err != nil) { return err }
            subscriptionStrings = append(subscriptionStrings, movie.Genres...)
            subscriptionStrings = append(subscriptionStrings, movie.Directors...)
            subscriptionStrings = append(subscriptionStrings, movie.Actors...)
            title = movie.Title
            description = movie.Description
        } else if (show) {
            show, err := getShow(uuid)
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
            conditions = append(conditions, expression.Name("target").Equal(expression.Value(value)))
        }

        filterExpression := conditions[0]
        for _, condition := range conditions[1:] {
            filterExpression = filterExpression.Or(condition)
        }

        expr, err := expression.NewBuilder().WithFilter(filterExpression).Build()

        if err != nil {
            fmt.Printf("Error building expression: %v\n", err)
            return err
        }


        var lastEvaluatedKey map[string]types.AttributeValue
        for {
            input := &dynamodb.ScanInput{
                TableName:                 aws.String(common.SubscriptionTableName),
                FilterExpression:          expr.Filter(),
                ExpressionAttributeNames:  expr.Names(),
                ExpressionAttributeValues: expr.Values(),
                Select:                    types.SelectSpecificAttributes,
                ProjectionExpression:      aws.String("userId,target"),
                Limit:                     aws.Int32(10),
                ExclusiveStartKey:         lastEvaluatedKey,
            }

            result, err := dynamodbClient.Scan(ctx, input)
            var throughputErr *types.ProvisionedThroughputExceededException
            if errors.As(err, &throughputErr) {
                time.Sleep(time.Duration(100000000))
                continue
            }
            if (err != nil) {
                log.Printf("Error running query: %v\n", err)
                return err
            }

            for _, item := range result.Items {
                userID, ok := item["userId"]
                if (!ok) {
                    log.Println("userId missing from the scan result")
                    return errors.New("userId missing from the scan result")
                }

                if (bloomFilter.TestAndAddString(userID.(*types.AttributeValueMemberS).Value)) {
                    continue
                }

                userInfo, err := cognitoClient.AdminGetUser(ctx, &cognitoidentityprovider.AdminGetUserInput{
                    UserPoolId: aws.String(cognitoUserPoolID),
                    Username:   aws.String(userID.(*types.AttributeValueMemberS).Value),
                })
                if err != nil {
                    log.Printf("Error getting user info: %v\n", err)
                    return err
                }

                var userEmail string = ""
                for _, attr := range userInfo.UserAttributes {
                    if *attr.Name == "email" {
                        userEmail = *attr.Value
                        break
                    }
                }

                if (len(userEmail) == 0) {
                    log.Printf("User email is missing!")
                    return errors.New("User email is missing!")
                }

                email := craftEmail(title, description, userEmail)
                log.Printf("User email: '%s'\n", userEmail)
                log.Println(email)
                _, err = sesClient.SendEmail(ctx, &email)

                if err != nil {
                    log.Printf("Error sending email: %v\n", err)
                    return err
                }
                log.Printf("Email sent successfully!")
            }

            lastEvaluatedKey = result.LastEvaluatedKey
            if (lastEvaluatedKey == nil) { break }
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
