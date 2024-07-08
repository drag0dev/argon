package common

import (
	_ "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	_ "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Video struct {
	// path on S3
	FileName string `dynamodbav:"fileName" json:"fileName"`
	FileType string `dynamodbav:"fileType" json:"fileType"`
	// size in bytes
	FileSize            uint64 `dynamodbav:"fileSize" json:"fileSize"`
	CreationTimestamp   int64  `dynamodbav:"creationTimestamp" json:"creationTimestamp"`
	LastChangeTimestamp int64  `dynamodbav:"lastChangeTimestamp" json:"lastChangeTimestamp"`
	// has video been processed and ready to be watched by the user
	Ready bool `dynamodbav:"ready" json:"ready"`
}

type Episode struct {
	EpisodeNumber uint64 `dynamodbav:"episodeNumber" json:"episodeNumber"`
	Title         string `dynamodbav:"title" json:"title"`
	Description   string `dynamodbav:"description" json:"description"`
	// String Set type
	Actors []string `dynamodbav:"actors" json:"actors"`
	// String Set type
	Directors []string `dynamodbav:"directors" json:"directors"`
	Video     Video    `dynamodbav:"video" json:"video"`
}

type Season struct {
	SeasonNumber uint64 `dynamodbav:"seasonNumber" json:"seasonNumber"`
	// List Type
	Episodes []Episode `dynamodbav:"episodes" json:"episodes"`
}

// NOTE: UUID is primary key, generated using google/uuid
type Movie struct {
	// google/uuid
	UUID        string `dynamodbav:"id" json:"id"`
	Title       string `dynamodbav:"title" json:"title"`
	Description string `dynamodbav:"description" json:"description"`
	// String Set type
	Genres []string `dynamodbav:"genres" json:"genres"`
	// String Set type
	Actors []string `dynamodbav:"actors" json:"actors"`
	// String Set type
	Directors []string `dynamodbav:"directors" json:"directors"`
	Video     Video    `dynamodbav:"video" json:"video"`
}

// NOTE: UUID is primary key, generated using google/uuid
type Show struct {
	// google/uuid
	UUID        string `dynamodbav:"id" json:"id"`
	Title       string `dynamodbav:"title" json:"title"`
	Description string `dynamodbav:"description" json:"description"`
	// String Set type
	Genres []string `dynamodbav:"genres" json:"genres"`
	// String Set type
	Actors []string `dynamodbav:"actors" json:"actors"`
	// String Set type
	Directors []string `dynamodbav:"directors" json:"directors"`
	// List Type
	Seasons []Season `dynamodbav:"seasons" json:"seasons"`
}

type SubscriptionType uint8

const (
	Actor SubscriptionType = iota
	Director
	Genre
)

// Subscription NOTE: UUID is primary key, generated using google/uuid
type Subscription struct {
	// google/uuid
	UUID     string           `dynamodbav:"id" json:"id"`
	UserUUID string           `dynamodbav:"userId" json:"userId"`
	Type     SubscriptionType `dynamodbav:"type" json:"type"`
	// The thing the user is subscribed to
	Target string `dynamodbav:"target" json:"target"`
	// For GSI partition key
	UserUUIDType string `dynamodbav:"userIdType"`
}

// Review NOTE: UUID is primary key, generated using google/uuid
type Review struct {
	// google/uuid
	UUID     string `dynamodbav:"id" json:"id"`
	UserUUID string `dynamodbav:"userId" json:"userId"`
	// The movie/show the user is reviewing
	TargetUUID string `dynamodbav:"targetId" json:"targetId"`
	Grade      uint8  `dynamodbav:"grade" json:"grade"`
	Comment    string `dynamodbav:"comment" json:"comment"`
}

type EditMetadataRequest struct {
	// The thing whose metadata is being edited
	TargetUUID    string   `json:"targetId"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Genres        []string `json:"genres"`
	Actors        []string `json:"actors"`
	Directors     []string `json:"directors"`
	SeasonNumber  *uint64  `json:"seasonNumber"`
	EpisodeNumber *uint64  `json:"episodeNumber"`
}

const VideoBucketName = "argon-videos-bucket"
const MovieTableName = "movie"
const ShowTableName = "show"
const SubscriptionTableName = "subscription"
const ReviewTableName = "review"

const SubscriptionTableSecondaryIndex = "subscription-secondary-index"

const SubscriptionQueueName = "subscription-queue"
const UnsubscriptionQueueName = "unsubscription-queue"
const ReviewQueueName = "review-queue"
const EditMetadataRequestQueueName = "edit-metadata-request-queue"

// email consts
var SenderEmail = "senderargon@gmail.com"
var ReceiverEmail1 = "receiverargon1@gmail.com"
var ReceiverEmail2 = "receiverargon2@gmail.com"

const Resolution1 = "1920:1080"
const Resolution2 = "1280:720"
const Resolution3 = "800:600"
const OriginalSuffix = "_original"
const UpdateSuffix = "_update"

// auth
const AdminGroupName = "argon-admin-group"
