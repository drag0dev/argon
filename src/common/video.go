package common

import (
    _ "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
    _ "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Video struct {
    // path on S3
    FileName string             `dynamodbav:"fileName"`
    FileType string             `dynamodbav:"fileType"`
    // size in bytes
    FileSize uint64             `dynamodbav:"fileSize"`
    CreationTimestamp int64     `dynamodbav:"creationTimestamp"`
    LastChangeTimestamp int64   `dynamodbav:"lastChangeTimestamp"`
}

type Episode struct {
    EpisodeNumber uint64        `dynamodbav:"episodeNumber"`
    Title string                `dynamodbav:"title"`
    Description string          `dynamodbav:"description"`
    // String Set type
    Actors []string             `dynamodbav:"actors"`
    // String Set type
    Directors []string          `dynamodbav:"directors"`
    Video Video                 `dynamodbav:"video"`
}

type Season struct {
    SeasonNumber uint64        `dynamodbav:"number"`
    // List Type
    Episodes []Episode         `dynamodbav:"episodes"`
}


//NOTE: UUID is primary key, generated using google/uuid
type Movie struct {
    // google/uuid
    UUID string                 `dynamodbav:"id"`
    Title string                `dynamodbav:"title"`
    Description string          `dynamodbav:"description"`
    // String Set type
    Genres []string             `dynamodbav:"genres"`
    // String Set type
    Actors []string             `dynamodbav:"actors"`
    // String Set type
    Directors []string          `dynamodbav:"directors"`
    Video Video                 `dynamodbav:"video"`
}

//NOTE: UUID is primary key, generated using google/uuid
type Show struct {
    // google/uuid
    UUID string                 `dynamodbav:"id"`
    Title string                `dynamodbav:"title"`
    Description string          `dynamodbav:"description"`
    // String Set type
    Genres []string             `dynamodbav:"genres"`
    // String Set type
    Actors []string             `dynamodbav:"actors"`
    // String Set type
    Directors []string          `dynamodbav:"directors"`
    // List Type
    Seasons []Season            `dynamodbav:"seasons"`
}
