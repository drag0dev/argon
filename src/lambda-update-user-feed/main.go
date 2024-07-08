package main

import (
	"common"
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
    "math/rand/v2"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var dynamodbClient *dynamodb.Client

type kv struct {
    Key   string
    Value float64
}

func giveTopThree(input map[string]float64) []string {
    var kvs []kv
    for key, value := range input {
        kvs = append(kvs, kv{Key: key, Value: value})
    }
    sort.Slice(kvs, func(i, j int) bool {
        return kvs[i].Value > kvs[j].Value
    })

    var res []string
    numberOfItems := 3
    if (len(kvs) < 3) {
        res = make([]string, len(kvs))
        numberOfItems = len(kvs)
    } else {
        res = make([]string, 3)
    }

    for i := 0; i < numberOfItems; i++ {
        res[i] = kvs[i].Key
    }

    return res
}

func queryMovie(topThree []string, field string) ([]string, error) {
    var conditions []expression.ConditionBuilder
    for _, value := range topThree {
        conditions = append(conditions, expression.Contains(expression.Name(field), value))
    }

    filterExpression := conditions[0]
    for _, condition := range conditions[1:] {
        filterExpression = filterExpression.Or(condition)
    }

    expr, err := expression.NewBuilder().WithFilter(filterExpression).Build()

    if err != nil {
        fmt.Printf("Error building expression: %v\n", err)
        return nil, err
    }

    input := &dynamodb.ScanInput{
        TableName:                 aws.String(common.MovieTableName),
        FilterExpression:          expr.Filter(),
        ExpressionAttributeNames:  expr.Names(),
        ExpressionAttributeValues: expr.Values(),
        Limit:                     aws.Int32(10),
    }

    result, err := dynamodbClient.Scan(context.TODO(), input)

    if (err != nil) {
        log.Printf("Error getting movies: %v", err)
        return nil, err
    }

    movies := []string{}
    for _, item := range result.Items {
        id, ok := item["id"]
        if (!ok) {
            log.Println("movie id missing from the scan result")
            return nil, errors.New("movie id missing from the scan result")
        }

        idCast, ok := id.(*types.AttributeValueMemberS)
        if (!ok) {
            log.Println("movie id is not a string")
            return nil, errors.New("movie id is not a string")
        }

        movies = append(movies, idCast.Value)
    }

    res := []string{}
    if (len(movies) == 0) { return []string{}, nil }

    uniqueInts := make(map[int]struct{})
    var randomInts []int
    for len(randomInts) < min(3, len(movies)) {
        num := rand.IntN(len(movies))
        if _, exists := uniqueInts[num]; !exists {
            uniqueInts[num] = struct{}{}
            randomInts = append(randomInts, num)
        }
    }

    for _, index := range randomInts {
        res = append(res, movies[index])
    }

    return res, nil
}

func queryShow(topThree []string, field string) ([]string, error) {
    var conditions []expression.ConditionBuilder
    for _, value := range topThree {
        conditions = append(conditions, expression.Contains(expression.Name(field), value))
    }

    filterExpression := conditions[0]
    for _, condition := range conditions[1:] {
        filterExpression = filterExpression.Or(condition)
    }

    expr, err := expression.NewBuilder().WithFilter(filterExpression).Build()

    if err != nil {
        fmt.Printf("Error building expression: %v\n", err)
        return nil, err
    }

    input := &dynamodb.ScanInput{
        TableName:                 aws.String(common.ShowTableName),
        FilterExpression:          expr.Filter(),
        ExpressionAttributeNames:  expr.Names(),
        ExpressionAttributeValues: expr.Values(),
        Limit:                     aws.Int32(10),
    }

    result, err := dynamodbClient.Scan(context.TODO(), input)

    if (err != nil) {
        log.Printf("Error getting show: %v", err)
        return nil, err
    }


    shows := []string{}
    for _, item := range result.Items {
        id, ok := item["id"]
        if (!ok) {
            log.Println("show id missing from the scan result")
            return nil, errors.New("show id missing from the scan result")
        }

        idCast, ok := id.(*types.AttributeValueMemberS)
        if (!ok) {
            log.Println("show id is not a string")
            return nil, errors.New("show id is not a string")
        }

        shows = append(shows, idCast.Value)
    }

    res := []string{}
    if (len(shows) == 0) { return []string{}, nil }

    uniqueInts := make(map[int]struct{})
    var randomInts []int
    for len(randomInts) < min(3, len(shows)) {
        num := rand.IntN(len(shows))
        if _, exists := uniqueInts[num]; !exists {
            uniqueInts[num] = struct{}{}
            randomInts = append(randomInts, num)
        }
    }

    for _, index := range randomInts {
        res = append(res, shows[index])
    }

    return res, nil
}

func uniqueSlice(slice []string) []string {
    keys := make(map[string]bool)
    list := []string{}
    for _, entry := range slice {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }
    return list
}

func handleFeedUpdate(ctx context.Context, sqsEvent events.SQSEvent) error {
    for _, message := range sqsEvent.Records {
        userId := message.Body
        input := &dynamodb.GetItemInput{
            TableName: aws.String(common.UserPreferenceTableName),
            Key: map[string]types.AttributeValue {
                "userId": &types.AttributeValueMemberS {
                    Value: userId,
                },
            },
        }

        result, err := dynamodbClient.GetItem(context.TODO(), input)
        if err != nil {
            log.Printf("Error getting preferences: %v", err)
            return err
        }

        if result.Item == nil {
            log.Printf("User does not exist?")
            return errors.New("USer does not exist?")
        }

        var preference common.UserPreference
        err = attributevalue.UnmarshalMap(result.Item, &preference)
        if err != nil {
            log.Printf("Error unmarshaling preference :%v", err)
            return err
        }

        topThreeActors := giveTopThree(preference.Actors)
        topThreeGenres := giveTopThree(preference.Genres)
        topThreeDirectors := giveTopThree(preference.Directors)

        moviesByActor, err := queryMovie(topThreeActors, "actors")
        if (err != nil) {return err}
        log.Println(moviesByActor)
        moviesByGenre, err := queryMovie(topThreeGenres, "genres")
        if (err != nil) {return err}
        log.Println(moviesByGenre)
        moviesByDirector, err := queryMovie(topThreeDirectors, "directors")
        if (err != nil) {return err}
        log.Println(moviesByDirector)

        showsByActors, err := queryShow(topThreeActors, "actors")
        if (err != nil) {return err}
        log.Println(showsByActors)
        showsByGenres, err := queryShow(topThreeGenres, "genres")
        if (err != nil) {return err}
        log.Println(showsByGenres)
        showsByDirectors, err := queryShow(topThreeDirectors, "directors")
        if (err != nil) {return err}
        log.Println(showsByDirectors)

        feedMovies := []string{}
        feedMovies = append(feedMovies, moviesByActor...)
        feedMovies = append(feedMovies, moviesByGenre...)
        feedMovies = append(feedMovies, moviesByDirector...)
        feedMovies = uniqueSlice(feedMovies)

        feedShows := []string{}
        feedShows = append(feedShows, showsByActors...)
        feedShows = append(feedShows, showsByGenres...)
        feedShows = append(feedShows, showsByDirectors...)
        feedShows = uniqueSlice(feedShows)

        // shuffled the feed
        rand.Shuffle(len(feedMovies), func(i, j int) {
            feedMovies[i], feedMovies[j] = feedMovies[j], feedMovies[i]
        })
        rand.Shuffle(len(feedShows), func(i, j int) {
            feedShows[i], feedShows[j] = feedShows[j], feedShows[i]
        })


        feedStruct := common.Feed{
            UserId: userId,
            FeedShows: feedShows,
            FeedMovies: feedMovies,
        }
        feedStructMarshaled, err := attributevalue.MarshalMap(feedStruct)

        putItemInput := &dynamodb.PutItemInput{
            TableName: aws.String(common.FeedTableName),
            Item: feedStructMarshaled,
        }

        _, err = dynamodbClient.PutItem(ctx, putItemInput)
        if err != nil {
            log.Printf("Error putting feed: %v", err)
            return err
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
    lambda.Start(handleFeedUpdate)
}
