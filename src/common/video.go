package common

import (
    _ "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
    _ "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Video struct {
    // path on S3
    FileName           string `dynamodbav:"fileName" json:"fileName"`
    FileType           string `dynamodbav:"fileType" json:"fileType"`
    // size in bytes
    FileSize           uint64 `dynamodbav:"fileSize" json:"fileSize"`
    CreationTimestamp  int64  `dynamodbav:"creationTimestamp" json:"creationTimestamp"`
    LastChangeTimestamp int64 `dynamodbav:"lastChangeTimestamp" json:"lastChangeTimestamp"`
}

type Episode struct {
    EpisodeNumber uint64  `dynamodbav:"episodeNumber" json:"episodeNumber"`
    Title         string  `dynamodbav:"title" json:"title"`
    Description   string  `dynamodbav:"description" json:"description"`
    // String Set type
    Actors        []string `dynamodbav:"actors" json:"actors"`
    // String Set type
    Directors     []string `dynamodbav:"directors" json:"directors"`
    Video         Video    `dynamodbav:"video" json:"video"`
}

type Season struct {
    SeasonNumber uint64    `dynamodbav:"number" json:"seasonNumber"`
    // List Type
    Episodes     []Episode `dynamodbav:"episodes" json:"episodes"`
}

//NOTE: UUID is primary key, generated using google/uuid
type Movie struct {
    // google/uuid
    UUID         string    `dynamodbav:"id" json:"id"`
    Title        string    `dynamodbav:"title" json:"title"`
    Description  string    `dynamodbav:"description" json:"description"`
    // String Set type
    Genres       []string  `dynamodbav:"genres" json:"genres"`
    // String Set type
    Actors       []string  `dynamodbav:"actors" json:"actors"`
    // String Set type
    Directors    []string  `dynamodbav:"directors" json:"directors"`
    Video        Video     `dynamodbav:"video" json:"video"`
}

//NOTE: UUID is primary key, generated using google/uuid
type Show struct {
    // google/uuid
    UUID         string    `dynamodbav:"id" json:"id"`
    Title        string    `dynamodbav:"title" json:"title"`
    Description  string    `dynamodbav:"description" json:"description"`
    // String Set type
    Genres       []string  `dynamodbav:"genres" json:"genres"`
    // String Set type
    Actors       []string  `dynamodbav:"actors" json:"actors"`
    // String Set type
    Directors    []string  `dynamodbav:"directors" json:"directors"`
    // List Type
    Seasons      []Season  `dynamodbav:"seasons" json:"seasons"`
}

const VideoBucketName = "argon-videos-bucket"
const MovieTableName = "movie"
