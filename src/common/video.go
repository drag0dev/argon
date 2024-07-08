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

const VideoBucketName = "argon-videos-bucket"
const MovieTableName = "movie"
const ShowTableName = "show"
const SubscriptionTableName = "subscription"

const SubscriptionTableSecondaryIndex = "subscription-secondary-index"

const SubscriptionQueueName = "subscription-queue"
const UnsubscriptionQueueName = "unsubscription-queue"

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

// feed
const PreferenceUpdateQueue = "preference-update-queue"
const UpdateFeedQueue = "update-feed-queue"

const UpdateCounterThreshhold = 35

const GetChangeWeight = 0.1
const GetUpdateWeight = 0.5

const SubscribeChangeWeight = 5
const SubcribeUpdateWeight = 5

const ReviewChangeWeight1 = -3
const ReviewChangeWeight2 = -2
const ReviewChangeWeight3 = -1
const ReviewChangeWeight4 = 2
const ReviewChangeWeight5 = 3
const ReviewUpdateWeight = 5

const UserPreferenceTableName = "user-preference-table"
