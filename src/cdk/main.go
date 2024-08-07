package main

import (
	"common"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambdaeventsources"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3notifications"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsses"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssns"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssnssubscriptions"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func generateLambdaIntegrationOptions() *awsapigateway.LambdaIntegrationOptions {
	return &awsapigateway.LambdaIntegrationOptions{
		Proxy: jsii.Bool(true),
	}
}
func generateMethodResponses() *[]*awsapigateway.MethodResponse {
	return &[]*awsapigateway.MethodResponse{
		{
			StatusCode: jsii.String("200"),
			ResponseParameters: &map[string]*bool{
				"method.response.header.Access-Control-Allow-Origin":  jsii.Bool(true),
				"method.response.header.Access-Control-Allow-Headers": jsii.Bool(true),
				"method.response.header.Access-Control-Allow-Methods": jsii.Bool(true),
			},
		},
		{
			StatusCode: jsii.String("400"),
			ResponseParameters: &map[string]*bool{
				"method.response.header.Access-Control-Allow-Origin":  jsii.Bool(true),
				"method.response.header.Access-Control-Allow-Headers": jsii.Bool(true),
				"method.response.header.Access-Control-Allow-Methods": jsii.Bool(true),
			},
		},
		{
			StatusCode: jsii.String("404"),
			ResponseParameters: &map[string]*bool{
				"method.response.header.Access-Control-Allow-Origin":  jsii.Bool(true),
				"method.response.header.Access-Control-Allow-Headers": jsii.Bool(true),
				"method.response.header.Access-Control-Allow-Methods": jsii.Bool(true),
			},
		},
		{
			StatusCode: jsii.String("500"),
			ResponseParameters: &map[string]*bool{
				"method.response.header.Access-Control-Allow-Origin":  jsii.Bool(true),
				"method.response.header.Access-Control-Allow-Headers": jsii.Bool(true),
				"method.response.header.Access-Control-Allow-Methods": jsii.Bool(true),
			},
		},
	}
}

func NewArgonStack(scope constructs.Construct, id string, props *awscdk.StackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, props)

	preferenceUpdateQueue := awssqs.NewQueue(stack, jsii.String("PreferenceUpdateQueue"), &awssqs.QueueProps{
		QueueName:         jsii.String(common.PreferenceUpdateQueue),
		VisibilityTimeout: awscdk.Duration_Minutes(jsii.Number(15)),
	})
	updateFeedQueue := awssqs.NewQueue(stack, jsii.String("UpdateFeedQueue"), &awssqs.QueueProps{
		QueueName:         jsii.String(common.UpdateFeedQueue),
		VisibilityTimeout: awscdk.Duration_Minutes(jsii.Number(15)),
	})

	userPreferenceTable := awsdynamodb.NewTable(stack, jsii.String("user-preference-table"), &awsdynamodb.TableProps{
		TableName: jsii.String(common.UserPreferenceTableName),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("userId"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode:   awsdynamodb.BillingMode_PROVISIONED,
		ReadCapacity:  jsii.Number(1),
		WriteCapacity: jsii.Number(1),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})
	awscdk.NewCfnOutput(stack, jsii.String("user preference table"), &awscdk.CfnOutputProps{
		Value:       userPreferenceTable.TableName(),
		Description: jsii.String("user-preference-table"),
	})

	userFeedTable := awsdynamodb.NewTable(stack, jsii.String("user-feed-table"), &awsdynamodb.TableProps{
		TableName: jsii.String(common.FeedTableName),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("userId"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode:   awsdynamodb.BillingMode_PROVISIONED,
		ReadCapacity:  jsii.Number(1),
		WriteCapacity: jsii.Number(1),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})
	awscdk.NewCfnOutput(stack, jsii.String("user feed table"), &awscdk.CfnOutputProps{
		Value:       userFeedTable.TableName(),
		Description: jsii.String("user-preference-table"),
	})

	// update user preference lambda
	updateUserPreferenceLambda := awslambda.NewFunction(stack, jsii.String("UpdateUserPreferenceLambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.Code_FromAsset(jsii.String("../lambda-update-user-preference/function.zip"), &awss3assets.AssetOptions{}),
		Timeout: awscdk.Duration_Minutes(jsii.Number(10)),
	})
	updateUserPreferenceLambda.AddEventSource(awslambdaeventsources.NewSqsEventSource(preferenceUpdateQueue, &awslambdaeventsources.SqsEventSourceProps{
		BatchSize: jsii.Number(10),
	}))
	preferenceUpdateQueue.GrantConsumeMessages(updateUserPreferenceLambda)
	updateFeedQueue.GrantSendMessages(updateUserPreferenceLambda)
	userPreferenceTable.GrantReadWriteData(updateUserPreferenceLambda)

	// updateFeedQueue := awssqs.NewQueue(stack, jsii.String("updateFeedQueue"), &awssqs.QueueProps{
	//     QueueName: jsii.String(common.UpdateFeedQueue),
	//     VisibilityTimeout:    awscdk.Duration_Minutes(jsii.Number(15)),
	// })

	// publishing topic
	publishingTopic := awssns.NewTopic(stack, jsii.String("PublishingTopic"), &awssns.TopicProps{
		TopicName: jsii.String("publishing-topic"),
	})

	// notification queue
	notificationQueue := awssqs.NewQueue(stack, jsii.String("NotificationQueue"), &awssqs.QueueProps{
		QueueName:         jsii.String("notification-queue"),
		VisibilityTimeout: awscdk.Duration_Minutes(jsii.Number(15)),
	})
	publishingTopic.AddSubscription(awssnssubscriptions.NewSqsSubscription(notificationQueue, nil))

	// Create a Cognito User Pool
	userPool := awscognito.NewUserPool(stack, jsii.String("ArgonUserPool"), &awscognito.UserPoolProps{
		UserPoolName:      jsii.String("argon-user-pool"),
		SelfSignUpEnabled: jsii.Bool(true),
		SignInAliases: &awscognito.SignInAliases{
			Email:    jsii.Bool(true),
			Username: jsii.Bool(true),
		},
		PasswordPolicy: &awscognito.PasswordPolicy{
			MinLength:        jsii.Number(8),
			RequireLowercase: jsii.Bool(true),
			RequireDigits:    jsii.Bool(true),
			RequireSymbols:   jsii.Bool(true),
			RequireUppercase: jsii.Bool(true),
		},
		StandardAttributes: &awscognito.StandardAttributes{
			Email: &awscognito.StandardAttribute{
				Required: jsii.Bool(true),
				Mutable:  jsii.Bool(false),
			},
		},
		AutoVerify: &awscognito.AutoVerifiedAttrs{
			Email: jsii.Bool(true),
		},
		UserVerification: &awscognito.UserVerificationConfig{
			EmailStyle: awscognito.VerificationEmailStyle_LINK,
		},
		CustomAttributes: &map[string]awscognito.ICustomAttribute{
			"firstName": awscognito.NewStringAttribute(&awscognito.StringAttributeProps{
				Mutable: jsii.Bool(true),
			}),
			"lastName": awscognito.NewStringAttribute(&awscognito.StringAttributeProps{
				Mutable: jsii.Bool(true),
			}),
			"dateOfBirth": awscognito.NewStringAttribute(&awscognito.StringAttributeProps{
				Mutable: jsii.Bool(true),
			}),
		},
	})
	userPool.AddDomain(aws.String("Verification Domain"), &awscognito.UserPoolDomainOptions{
		CognitoDomain: &awscognito.CognitoDomainOptions{
			DomainPrefix: aws.String("argon-verification-domain"),
		},
	})
	awscdk.NewCfnOutput(stack, jsii.String("Argon User Pool"), &awscdk.CfnOutputProps{
		Value:       userPool.UserPoolId(),
		Description: jsii.String("Argon User Pool"),
	})

	// Create a Cognito User Pool Client
	userPoolClient := awscognito.NewUserPoolClient(stack, jsii.String("ArgonFrontend"), &awscognito.UserPoolClientProps{
		UserPool:       userPool,
		GenerateSecret: jsii.Bool(false),
		AuthFlows: &awscognito.AuthFlow{
			UserSrp:      jsii.Bool(true),
			UserPassword: jsii.Bool(true),
		},
	})
	awscdk.NewCfnOutput(stack, jsii.String("Argon Frontend"), &awscdk.CfnOutputProps{
		Value:       userPoolClient.UserPoolClientId(),
		Description: jsii.String("ArgonFrontend"),
	})

	_ = awscognito.NewCfnUserPoolGroup(stack, jsii.String("AdminGroup"), &awscognito.CfnUserPoolGroupProps{
		GroupName:  jsii.String(common.AdminGroupName),
		UserPoolId: userPool.UserPoolId(),
	})

	adminAuthorizerLambda := awslambda.NewFunction(stack, jsii.String("AdminAuthorizerFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.Code_FromAsset(jsii.String("../lambda-admin-authorizer/function.zip"), &awss3assets.AssetOptions{}),
		Environment: &map[string]*string{
			"COGNITO_USER_POOL_ID": userPool.UserPoolId(),
		},
	})
	userPool.Grant(adminAuthorizerLambda, aws.String("cognito-idp:AdminListGroupsForUser"))

	userPoolAuthorizer := awsapigateway.NewCognitoUserPoolsAuthorizer(stack, jsii.String("userPoolAuthorizer"), &awsapigateway.CognitoUserPoolsAuthorizerProps{
		CognitoUserPools: &[]awscognito.IUserPool{userPool},
		IdentitySource:   jsii.String("method.request.header.Authorization"),
	})

	adminAuthorizer := awsapigateway.NewTokenAuthorizer(stack, jsii.String("AdminAuthorizer"), &awsapigateway.TokenAuthorizerProps{
		Handler: adminAuthorizerLambda,
	})

	// Video bucket
	videoBucket := awss3.NewBucket(stack, jsii.String("argon-videos-bucket"), &awss3.BucketProps{
		BucketName:        jsii.String("argon-videos-bucket"),
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
		AutoDeleteObjects: jsii.Bool(true),
	})
	star := "*"
	videoBucket.AddCorsRule(&awss3.CorsRule{
		AllowedMethods: &[]awss3.HttpMethods{
			awss3.HttpMethods_GET,
			awss3.HttpMethods_PUT,
			awss3.HttpMethods_POST,
		},
		AllowedOrigins: &[]*string{&star},
		AllowedHeaders: &[]*string{&star},
	})
	awscdk.NewCfnOutput(stack, jsii.String("argon videos bucket"), &awscdk.CfnOutputProps{
		Value:       jsii.String("argon-videos-bucket"),
		Description: jsii.String("argon-videos-bucket"),
	})

	// Movie and show tables
	movieTable := awsdynamodb.NewTable(stack, jsii.String("movie-table"), &awsdynamodb.TableProps{
		TableName: jsii.String("movie"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode:   awsdynamodb.BillingMode_PROVISIONED,
		ReadCapacity:  jsii.Number(1),
		WriteCapacity: jsii.Number(1),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})
	awscdk.NewCfnOutput(stack, jsii.String("movie table"), &awscdk.CfnOutputProps{
		Value:       movieTable.TableName(),
		Description: jsii.String("movie-table"),
	})

	showTable := awsdynamodb.NewTable(stack, jsii.String("show-table"), &awsdynamodb.TableProps{
		TableName: jsii.String("show"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode:   awsdynamodb.BillingMode_PROVISIONED,
		ReadCapacity:  jsii.Number(1),
		WriteCapacity: jsii.Number(1),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})
	awscdk.NewCfnOutput(stack, jsii.String("show table"), &awscdk.CfnOutputProps{
		Value:       showTable.TableName(),
		Description: jsii.String("show-table"),
	})

	// update user feed lambda
	updateUserFeedLambda := awslambda.NewFunction(stack, jsii.String("UpdateUserFeedLambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.Code_FromAsset(jsii.String("../lambda-update-user-feed/function.zip"), &awss3assets.AssetOptions{}),
		Timeout: awscdk.Duration_Minutes(jsii.Number(10)),
	})
	updateUserFeedLambda.AddEventSource(awslambdaeventsources.NewSqsEventSource(updateFeedQueue, &awslambdaeventsources.SqsEventSourceProps{
		BatchSize: jsii.Number(10),
	}))
	updateFeedQueue.GrantConsumeMessages(updateUserFeedLambda)
	userFeedTable.GrantReadWriteData(updateUserFeedLambda)
	userPreferenceTable.GrantReadData(updateUserFeedLambda)
	movieTable.GrantReadData(updateUserFeedLambda)
	showTable.GrantReadData(updateUserFeedLambda)

	// Subscription table
	subscriptionTable := awsdynamodb.NewTable(stack, jsii.String("subscription-table"), &awsdynamodb.TableProps{
		TableName: jsii.String("subscription"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode:   awsdynamodb.BillingMode_PROVISIONED,
		ReadCapacity:  jsii.Number(1),
		WriteCapacity: jsii.Number(1),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})
	subscriptionTable.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName: jsii.String("subscription-secondary-index"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("userIdType"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("target"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		ReadCapacity:  jsii.Number(1),
		WriteCapacity: jsii.Number(1),
	})
	subscriptionTable.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName: jsii.String("user-id-secondary-index"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("userId"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		ReadCapacity:  jsii.Number(1),
		WriteCapacity: jsii.Number(1),
	})
	awscdk.NewCfnOutput(stack, jsii.String("subscription table"), &awscdk.CfnOutputProps{
		Value:       subscriptionTable.TableName(),
		Description: jsii.String("subscription-table"),
	})

	// Review table
	reviewTable := awsdynamodb.NewTable(stack, jsii.String("review-table"), &awsdynamodb.TableProps{
		TableName: jsii.String("review"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode:   awsdynamodb.BillingMode_PROVISIONED,
		ReadCapacity:  jsii.Number(1),
		WriteCapacity: jsii.Number(1),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})
	awscdk.NewCfnOutput(stack, jsii.String("review table"), &awscdk.CfnOutputProps{
		Value:       reviewTable.TableName(),
		Description: jsii.String("review-table"),
	})

	// subscription queue
	subscriptionQueue := awssqs.NewQueue(stack, jsii.String("SubscriptionQueue"), &awssqs.QueueProps{
		QueueName: jsii.String("subscription-queue"),
	})

	// SES
	senderEmailIdentity := awsses.NewEmailIdentity(stack, jsii.String("ArgonSenderEmailIdentity"), &awsses.EmailIdentityProps{
		Identity: awsses.Identity_Email(jsii.String(common.SenderEmail)),
	})
	receiverEmailOneIdentity := awsses.NewEmailIdentity(stack, jsii.String("ArgonReceiverEmailIdentityOne"), &awsses.EmailIdentityProps{
		Identity: awsses.Identity_Email(jsii.String(common.ReceiverEmail1)),
	})
	receiverEmailTwoIdentity := awsses.NewEmailIdentity(stack, jsii.String("ArgonReceiverEmailIdentityTwo"), &awsses.EmailIdentityProps{
		Identity: awsses.Identity_Email(jsii.String(common.ReceiverEmail2)),
	})
	// configSet := awsses.NewConfigurationSet(stack, jsii.String("ArgonEmailConfigSet"), &awsses.ConfigurationSetProps{
	//     ConfigurationSetName: jsii.String("argon-email-config-set"),
	// })

	// create notifications lambda
	createNotificationLambda := awslambda.NewFunction(stack, jsii.String("CreateNotificationLambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.Code_FromAsset(jsii.String("../lambda-create-notification/function.zip"), &awss3assets.AssetOptions{}),
		Environment: &map[string]*string{
			"COGNITO_USER_POOL_ID": userPool.UserPoolId(),
		},
		Timeout: awscdk.Duration_Minutes(jsii.Number(10)),
	})
	createNotificationLambda.AddEventSource(awslambdaeventsources.NewSqsEventSource(notificationQueue, &awslambdaeventsources.SqsEventSourceProps{
		BatchSize: jsii.Number(1),
	}))
	notificationQueue.GrantConsumeMessages(createNotificationLambda)
	adminGetUserPermission := "cognito-idp:AdminGetUser"
	userPool.Grant(createNotificationLambda, &adminGetUserPermission)
	senderEmailIdentity.GrantSendEmail(createNotificationLambda)
	subscriptionTable.GrantReadData(createNotificationLambda)
	movieTable.GrantReadData(createNotificationLambda)
	showTable.GrantReadData(createNotificationLambda)
	receiverEmailOneIdentity.GrantSendEmail(createNotificationLambda)
	receiverEmailTwoIdentity.GrantSendEmail(createNotificationLambda)

	// Transcoding lambda
	ffmpegLayer := awslambda.NewLayerVersion(stack, jsii.String("FFmpegLayer"), &awslambda.LayerVersionProps{
		Code:        awslambda.Code_FromAsset(jsii.String("../lambda-transcoder/ffmpeg.zip"), &awss3assets.AssetOptions{}),
		Description: jsii.String("FFmpeg binary"),
		CompatibleRuntimes: &[]awslambda.Runtime{
			awslambda.Runtime_PROVIDED_AL2023(),
		},
	})
	transcoderLambda := awslambda.NewFunction(stack, jsii.String("VideoTranscoding"), &awslambda.FunctionProps{
		Runtime:    awslambda.Runtime_PROVIDED_AL2023(),
		Handler:    jsii.String("main"),
		Code:       awslambda.Code_FromAsset(jsii.String("../lambda-transcoder/function.zip"), &awss3assets.AssetOptions{}),
		Timeout:    awscdk.Duration_Minutes(jsii.Number(2)),
		MemorySize: jsii.Number(1024),
		Layers: &[]awslambda.ILayerVersion{
			ffmpegLayer,
		},
		Environment: &map[string]*string{
			"PUBLISHING_TOPIC_ARN": publishingTopic.TopicArn(),
		},
	})
	videoBucket.GrantReadWrite(transcoderLambda, jsii.String("*"))
	videoBucket.AddEventNotification(awss3.EventType_OBJECT_CREATED,
		awss3notifications.NewLambdaDestination(transcoderLambda),
		&awss3.NotificationKeyFilter{
			Suffix: jsii.String("_original"),
		},
	)
	movieTable.GrantReadWriteData(transcoderLambda)
	showTable.GrantReadWriteData(transcoderLambda)
	publishingTopic.GrantPublish(transcoderLambda)

	// unsubscription queue
	unsubscriptionQueue := awssqs.NewQueue(stack, jsii.String("UnsubscriptionQueue"), &awssqs.QueueProps{
		QueueName: jsii.String("unsubscription-queue"),
	})

	// review queue
	reviewQueue := awssqs.NewQueue(stack, jsii.String("ReviewQueue"), &awssqs.QueueProps{
		QueueName: jsii.String("review-queue"),
	})

	// edit metadata request queue
	editMetadataRequestQueue := awssqs.NewQueue(stack, jsii.String("EditMetadataRequestQueue"), &awssqs.QueueProps{
		QueueName: jsii.String("edit-metadata-request-queue"),
	})

	// Movie Lambdas
	getMovieLambda := awslambda.NewFunction(stack, jsii.String("GetMovie"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.Code_FromAsset(jsii.String("../lambda-get-movie/function.zip"), &awss3assets.AssetOptions{}),
	})
	postMovieLambda := awslambda.NewFunction(stack, jsii.String("PostMovie"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.Code_FromAsset(jsii.String("../lambda-upload-movie/function.zip"), &awss3assets.AssetOptions{}),
	})
	deleteMovieLambda := awslambda.NewFunction(stack, jsii.String("DeleteMovie"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.Code_FromAsset(jsii.String("../lambda-delete-movie/function.zip"), &awss3assets.AssetOptions{}),
	})
	updateMovieVideo := awslambda.NewFunction(stack, jsii.String("UpdateMovieVideo"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.Code_FromAsset(jsii.String("../lambda-update-video-movie/function.zip"), &awss3assets.AssetOptions{}),
	})
	videoBucket.GrantRead(getMovieLambda, jsii.String("*"))
	preferenceUpdateQueue.GrantSendMessages(getMovieLambda)
	videoBucket.GrantPut(postMovieLambda, jsii.String("*"))
	videoBucket.GrantDelete(deleteMovieLambda, jsii.String("*"))
	videoBucket.GrantReadWrite(updateMovieVideo, jsii.String("*"))
	movieTable.GrantReadData(getMovieLambda)
	movieTable.GrantWriteData(postMovieLambda)
	movieTable.GrantReadWriteData(deleteMovieLambda)
	movieTable.GrantReadWriteData(updateMovieVideo)

	// Tv Show Lambdas
	getShowLambda := awslambda.NewFunction(stack, jsii.String("GetTvShow"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.Code_FromAsset(jsii.String("../lambda-get-show/function.zip"), &awss3assets.AssetOptions{}),
	})
	postShowLambda := awslambda.NewFunction(stack, jsii.String("PostTvShow"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.Code_FromAsset(jsii.String("../lambda-upload-show/function.zip"), &awss3assets.AssetOptions{}),
	})
	deleteShowLambda := awslambda.NewFunction(stack, jsii.String("DeleteTvShow"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.Code_FromAsset(jsii.String("../lambda-delete-show/function.zip"), &awss3assets.AssetOptions{}),
	})
	updateShowVideo := awslambda.NewFunction(stack, jsii.String("UpdateShowVideo"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.Code_FromAsset(jsii.String("../lambda-update-video-show/function.zip"), &awss3assets.AssetOptions{}),
	})
	videoBucket.GrantRead(getShowLambda, jsii.String("*"))
	preferenceUpdateQueue.GrantSendMessages(getShowLambda)
	videoBucket.GrantPut(postShowLambda, jsii.String("*"))
	videoBucket.GrantDelete(deleteShowLambda, jsii.String("*"))
	videoBucket.GrantReadWrite(updateShowVideo, jsii.String("*"))
	showTable.GrantReadData(getShowLambda)
	showTable.GrantWriteData(postShowLambda)
	showTable.GrantReadWriteData(deleteShowLambda)
	showTable.GrantReadWriteData(updateShowVideo)

	// Subscription Lambdas
	queueSubscriptionLambda := awslambda.NewFunction(
		stack,
		jsii.String("QueueSubscription"),
		&awslambda.FunctionProps{
			Runtime: awslambda.Runtime_PROVIDED_AL2023(),
			Handler: jsii.String("main"),
			Code: awslambda.Code_FromAsset(
				jsii.String("../lambda-queue-subscription/function.zip"),
				&awss3assets.AssetOptions{},
			),
		},
	)
	preferenceUpdateQueue.GrantSendMessages(queueSubscriptionLambda)
	subscribeLambda := awslambda.NewFunction(stack, jsii.String("Subscribe"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code: awslambda.Code_FromAsset(
			jsii.String("../lambda-subscribe/function.zip"),
			&awss3assets.AssetOptions{},
		),
	})
	queueUnsubscriptionLambda := awslambda.NewFunction(
		stack,
		jsii.String("QueueUnsubscription"),
		&awslambda.FunctionProps{
			Runtime: awslambda.Runtime_PROVIDED_AL2023(),
			Handler: jsii.String("main"),
			Code: awslambda.Code_FromAsset(
				jsii.String("../lambda-queue-unsubscription/function.zip"),
				&awss3assets.AssetOptions{},
			),
		},
	)
	preferenceUpdateQueue.GrantSendMessages(queueUnsubscriptionLambda)
	unsubscribeLambda := awslambda.NewFunction(stack, jsii.String("Unsubscribe"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code: awslambda.Code_FromAsset(
			jsii.String("../lambda-unsubscribe/function.zip"),
			&awss3assets.AssetOptions{},
		),
	})
	subscriptionQueue.GrantSendMessages(queueSubscriptionLambda)
	subscribeLambda.AddEventSource(awslambdaeventsources.NewSqsEventSource(
		subscriptionQueue,
		&awslambdaeventsources.SqsEventSourceProps{
			BatchSize: jsii.Number(1),
		},
	))
	subscriptionQueue.GrantConsumeMessages(subscribeLambda)
	unsubscriptionQueue.GrantSendMessages(queueUnsubscriptionLambda)
	unsubscribeLambda.AddEventSource(awslambdaeventsources.NewSqsEventSource(
		unsubscriptionQueue,
		&awslambdaeventsources.SqsEventSourceProps{
			BatchSize: jsii.Number(1),
		},
	))
	unsubscriptionQueue.GrantConsumeMessages(unsubscribeLambda)
	subscriptionTable.GrantReadData(queueSubscriptionLambda)
	subscriptionTable.GrantReadData(queueUnsubscriptionLambda)
	subscriptionTable.GrantWriteData(unsubscribeLambda)
	subscriptionTable.GrantReadWriteData(subscribeLambda)

	// Review Lambdas
	queueReviewLambda := awslambda.NewFunction(stack, jsii.String("QueueReview"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code: awslambda.Code_FromAsset(
			jsii.String("../lambda-queue-review/function.zip"),
			&awss3assets.AssetOptions{},
		),
	})
	reviewLambda := awslambda.NewFunction(stack, jsii.String("Review"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code: awslambda.Code_FromAsset(
			jsii.String("../lambda-review/function.zip"),
			&awss3assets.AssetOptions{},
		),
	})
	reviewQueue.GrantSendMessages(queueReviewLambda)
	reviewLambda.AddEventSource(awslambdaeventsources.NewSqsEventSource(
		reviewQueue,
		&awslambdaeventsources.SqsEventSourceProps{
			BatchSize: jsii.Number(1),
		},
	))
	reviewQueue.GrantConsumeMessages(reviewLambda)
	movieTable.GrantReadData(queueReviewLambda)
	showTable.GrantReadData(queueReviewLambda)
	movieTable.GrantReadData(reviewLambda)
	showTable.GrantReadData(reviewLambda)
	reviewTable.GrantWriteData(reviewLambda)
	preferenceUpdateQueue.GrantSendMessages(queueReviewLambda)

	// Edit metadata Lambdas
	queueEditMetadataLambda := awslambda.NewFunction(
		stack,
		jsii.String("QueueEditMetadata"),
		&awslambda.FunctionProps{
			Runtime: awslambda.Runtime_PROVIDED_AL2023(),
			Handler: jsii.String("main"),
			Code: awslambda.Code_FromAsset(
				jsii.String("../lambda-queue-edit-metadata/function.zip"),
				&awss3assets.AssetOptions{},
			),
		},
	)
	editMetadataLambda := awslambda.NewFunction(
		stack,
		jsii.String("EditMetadata"),
		&awslambda.FunctionProps{
			Runtime: awslambda.Runtime_PROVIDED_AL2023(),
			Handler: jsii.String("main"),
			Code: awslambda.Code_FromAsset(
				jsii.String("../lambda-edit-metadata/function.zip"),
				&awss3assets.AssetOptions{},
			),
		},
	)
	editMetadataRequestQueue.GrantSendMessages(queueEditMetadataLambda)
	editMetadataLambda.AddEventSource(awslambdaeventsources.NewSqsEventSource(
		editMetadataRequestQueue,
		&awslambdaeventsources.SqsEventSourceProps{
			BatchSize: jsii.Number(1),
		},
	))
	editMetadataRequestQueue.GrantConsumeMessages(editMetadataLambda)
	movieTable.GrantReadData(queueEditMetadataLambda)
	movieTable.GrantReadWriteData(editMetadataLambda)
	showTable.GrantReadData(queueEditMetadataLambda)
	showTable.GrantReadWriteData(editMetadataLambda)

	// Create an API Gateway
	api := awsapigateway.NewRestApi(stack, jsii.String("ArgonAPI"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String("ArgonAPI"),
		Description: jsii.String("ArgonAPI"),
	})

	// Output the API URL
	awscdk.NewCfnOutput(stack, jsii.String("GatewayURL"), &awscdk.CfnOutputProps{
		Value:       api.Url(),
		Description: jsii.String("API Gateway URL"),
	})

	// Api Gateway movie resource
	movieApiResource := api.Root().AddResource(jsii.String("movie"), nil)
	movieApiResource.AddCorsPreflight(&awsapigateway.CorsOptions{
		AllowOrigins: awsapigateway.Cors_ALL_ORIGINS(),
		AllowMethods: jsii.Strings("GET", "POST", "PUT", "DELETE", "OPTIONS"),
		AllowHeaders: jsii.Strings(
			"Content-Type",
			"X-Amz-Date",
			"Authorization",
			"X-Api-Key",
			"X-Amz-Security-Token",
		),
	})
	movieApiResource.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(getMovieLambda, generateLambdaIntegrationOptions()), &awsapigateway.MethodOptions{
		AuthorizationType: awsapigateway.AuthorizationType_COGNITO,
		Authorizer:        userPoolAuthorizer,
		MethodResponses:   generateMethodResponses(),
	})
	movieApiResource.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(postMovieLambda, generateLambdaIntegrationOptions()), &awsapigateway.MethodOptions{
		Authorizer:      adminAuthorizer,
		MethodResponses: generateMethodResponses(),
	})
	movieApiResource.AddMethod(jsii.String("DELETE"), awsapigateway.NewLambdaIntegration(deleteMovieLambda, generateLambdaIntegrationOptions()), &awsapigateway.MethodOptions{
		Authorizer:      adminAuthorizer,
		MethodResponses: generateMethodResponses(),
	})
	movieApiResource.AddMethod(jsii.String("PUT"), awsapigateway.NewLambdaIntegration(updateMovieVideo, generateLambdaIntegrationOptions()), &awsapigateway.MethodOptions{
		Authorizer:      adminAuthorizer,
		MethodResponses: generateMethodResponses(),
	})

	// Api GateWay tv show resource
	tvShowApiResource := api.Root().AddResource(jsii.String("tvShow"), nil)
	tvShowApiResource.AddCorsPreflight(&awsapigateway.CorsOptions{
		AllowOrigins: awsapigateway.Cors_ALL_ORIGINS(),
		AllowMethods: jsii.Strings("GET", "POST", "PUT", "DELETE", "OPTIONS"),
		AllowHeaders: jsii.Strings(
			"Content-Type",
			"X-Amz-Date",
			"Authorization",
			"X-Api-Key",
			"X-Amz-Security-Token",
		),
	})
	tvShowApiResource.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(getShowLambda, generateLambdaIntegrationOptions()), &awsapigateway.MethodOptions{
		AuthorizationType: awsapigateway.AuthorizationType_COGNITO,
		Authorizer:        userPoolAuthorizer,
		MethodResponses:   generateMethodResponses(),
	})
	tvShowApiResource.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(postShowLambda, generateLambdaIntegrationOptions()), &awsapigateway.MethodOptions{
		Authorizer:      adminAuthorizer,
		MethodResponses: generateMethodResponses(),
	})
	tvShowApiResource.AddMethod(jsii.String("DELETE"), awsapigateway.NewLambdaIntegration(deleteShowLambda, generateLambdaIntegrationOptions()), &awsapigateway.MethodOptions{
		Authorizer:      adminAuthorizer,
		MethodResponses: generateMethodResponses(),
	})
	tvShowApiResource.AddMethod(jsii.String("PUT"), awsapigateway.NewLambdaIntegration(updateShowVideo, generateLambdaIntegrationOptions()), &awsapigateway.MethodOptions{
		Authorizer:      adminAuthorizer,
		MethodResponses: generateMethodResponses(),
	})

	// API gateway subscription resource
	subscriptionApiResource := api.Root().AddResource(jsii.String("subscription"), nil)
	subscriptionApiResource.AddCorsPreflight(&awsapigateway.CorsOptions{
		AllowOrigins: awsapigateway.Cors_ALL_ORIGINS(),
		AllowMethods: jsii.Strings("GET", "POST", "PUT", "DELETE", "OPTIONS"),
		AllowHeaders: jsii.Strings(
			"Content-Type",
			"X-Amz-Date",
			"Authorization",
			"X-Api-Key",
			"X-Amz-Security-Token",
		),
	})
	subscriptionApiResource.AddMethod(
		jsii.String("POST"),
		awsapigateway.NewLambdaIntegration(queueSubscriptionLambda, generateLambdaIntegrationOptions()),
		&awsapigateway.MethodOptions{
			AuthorizationType: awsapigateway.AuthorizationType_COGNITO,
			Authorizer:        userPoolAuthorizer,
			MethodResponses:   generateMethodResponses(),
		},
	)
	subscriptionApiResource.AddMethod(
		jsii.String("DELETE"),
		awsapigateway.NewLambdaIntegration(queueUnsubscriptionLambda, generateLambdaIntegrationOptions()),
		&awsapigateway.MethodOptions{
			AuthorizationType: awsapigateway.AuthorizationType_COGNITO,
			Authorizer:        userPoolAuthorizer,
			MethodResponses:   generateMethodResponses(),
		},
	)

	// API gateway review resource
	reviewApiResource := api.Root().AddResource(jsii.String("review"), nil)
	reviewApiResource.AddCorsPreflight(&awsapigateway.CorsOptions{
		AllowOrigins: awsapigateway.Cors_ALL_ORIGINS(),
		AllowMethods: jsii.Strings("GET", "POST", "PUT", "DELETE", "OPTIONS"),
		AllowHeaders: jsii.Strings(
			"Content-Type",
			"X-Amz-Date",
			"Authorization",
			"X-Api-Key",
			"X-Amz-Security-Token",
		),
	})
	reviewApiResource.AddMethod(
		jsii.String("POST"),
		awsapigateway.NewLambdaIntegration(queueReviewLambda, generateLambdaIntegrationOptions()),
		&awsapigateway.MethodOptions{
			AuthorizationType: awsapigateway.AuthorizationType_COGNITO,
			Authorizer:        userPoolAuthorizer,
			MethodResponses:   generateMethodResponses(),
		},
	)

	// API gateway edit metadata resource
	editMetadataApiResource := api.Root().AddResource(jsii.String("editMetadata"), nil)
	editMetadataApiResource.AddCorsPreflight(&awsapigateway.CorsOptions{
		AllowOrigins: awsapigateway.Cors_ALL_ORIGINS(),
		AllowMethods: jsii.Strings("GET", "POST", "PUT", "DELETE", "OPTIONS"),
		AllowHeaders: jsii.Strings(
			"Content-Type",
			"X-Amz-Date",
			"Authorization",
			"X-Api-Key",
			"X-Amz-Security-Token",
		),
	})
	editMetadataApiResource.AddMethod(
		jsii.String("PUT"),
		awsapigateway.NewLambdaIntegration(queueEditMetadataLambda, generateLambdaIntegrationOptions()),
		&awsapigateway.MethodOptions{
			Authorizer:      adminAuthorizer,
			MethodResponses: generateMethodResponses(),
		},
	)

	return stack
}

func main() {
	app := awscdk.NewApp(nil)
	NewArgonStack(app, "ArgonStack", &awscdk.StackProps{Env: &awscdk.Environment{
		Account: jsii.String(os.Getenv("ACCOUNT_ID")),
		Region:  jsii.String("eu-central-1"),
	}})
	app.Synth(nil)
}
