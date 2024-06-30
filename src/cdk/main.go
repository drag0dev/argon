package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func NewArgonStack(scope constructs.Construct, id string, props *awscdk.StackProps) awscdk.Stack {
    stack := awscdk.NewStack(scope, &id, props)

    // Create a Cognito User Pool
    userPool := awscognito.NewUserPool(stack, jsii.String("ArgonUserPool"), &awscognito.UserPoolProps{
        UserPoolName: jsii.String("argon-user-pool"),
        SelfSignUpEnabled: jsii.Bool(true),
        SignInAliases: &awscognito.SignInAliases{
            Email: jsii.Bool(true),
        },
        PasswordPolicy: &awscognito.PasswordPolicy{
            MinLength: jsii.Number(8),
            RequireLowercase: jsii.Bool(true),
            RequireDigits: jsii.Bool(true),
            RequireSymbols: jsii.Bool(true),
            RequireUppercase: jsii.Bool(true),
        },
        StandardAttributes: &awscognito.StandardAttributes{
            Email: &awscognito.StandardAttribute{
                Required: jsii.Bool(true),
                Mutable: jsii.Bool(false),
            },
        },
    })
    awscdk.NewCfnOutput(stack, jsii.String("Argon User Pool"), &awscdk.CfnOutputProps{
        Value:       userPool.UserPoolId(),
        Description: jsii.String("Argon User Pool"),
    })

    // Create a Cognito User Pool Client
    userPoolClient := awscognito.NewUserPoolClient(stack, jsii.String("ArgonFrontend"), &awscognito.UserPoolClientProps{
        UserPool: userPool,
        GenerateSecret: jsii.Bool(false),
    })
    awscdk.NewCfnOutput(stack, jsii.String("Argon Frontend"), &awscdk.CfnOutputProps{
        Value:       userPoolClient.UserPoolClientId(),
        Description: jsii.String("ArgonFrontend"),
    })



    // Video bucket
    videoBucket := awss3.NewBucket(stack, jsii.String("argon-videos-bucket"), &awss3.BucketProps{
        BucketName: jsii.String("argon-videos-bucket"),
    })
    awscdk.NewCfnOutput(stack, jsii.String("argon videos bucket"), &awscdk.CfnOutputProps{
        Value:       jsii.String("argon-videos-bucket"),
        Description: jsii.String("argon-videos-bucket"),
    })



    // Movie and show tables
    movieTable := awsdynamodb.NewTable(stack, jsii.String("movie-table"), &awsdynamodb.TableProps{
        TableName:   jsii.String("movie"),
        PartitionKey: &awsdynamodb.Attribute{
            Name: jsii.String("id"),
            Type: awsdynamodb.AttributeType_STRING,
        },
        BillingMode: awsdynamodb.BillingMode_PROVISIONED,
        ReadCapacity: jsii.Number(1),
        WriteCapacity: jsii.Number(1),
    })
    awscdk.NewCfnOutput(stack, jsii.String("movie table"), &awscdk.CfnOutputProps{
        Value:       movieTable.TableName(),
        Description: jsii.String("movie-table"),
    })

    showTable := awsdynamodb.NewTable(stack, jsii.String("show-table"), &awsdynamodb.TableProps{
        TableName:   jsii.String("show"),
        PartitionKey: &awsdynamodb.Attribute{
            Name: jsii.String("id"),
            Type: awsdynamodb.AttributeType_STRING,
        },
        BillingMode: awsdynamodb.BillingMode_PROVISIONED,
        ReadCapacity: jsii.Number(1),
        WriteCapacity: jsii.Number(1),
    })
    awscdk.NewCfnOutput(stack, jsii.String("show table"), &awscdk.CfnOutputProps{
        Value:       showTable.TableName(),
        Description: jsii.String("show-table"),
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
    videoBucket.GrantRead(getMovieLambda, jsii.String("*"));
    videoBucket.GrantPut(postMovieLambda, jsii.String("*"));
    movieTable.GrantReadData(getMovieLambda);
    movieTable.GrantWriteData(postMovieLambda);


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
    videoBucket.GrantRead(getShowLambda, jsii.String("*"));
    videoBucket.GrantPut(postShowLambda, jsii.String("*"));
    showTable.GrantReadData(getShowLambda);
    showTable.GrantWriteData(postShowLambda);



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

    // TODO: enable when frontend is done
    // Create a Cognito Authorizer
    // authorizer := awsapigateway.NewCognitoUserPoolsAuthorizer(stack, jsii.String("ArgonCognitoAuthorizer"), &awsapigateway.CognitoUserPoolsAuthorizerProps{
    //     CognitoUserPools: &[]awscognito.IUserPool{userPool},
    // })



    // Api Gateway movie resource
    movieApiResource := api.Root().AddResource(jsii.String("movie"), nil)
    movieApiResource.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(getMovieLambda, nil), &awsapigateway.MethodOptions{
        // TODO: enable when frontend is done
        // AuthorizationType: awsapigateway.AuthorizationType_COGNITO,
        // Authorizer:        authorizer,
    })
    movieApiResource.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(postMovieLambda, nil), &awsapigateway.MethodOptions{
        // TODO: enable when frontend is done
        // AuthorizationType: awsapigateway.AuthorizationType_COGNITO,
        // Authorizer:        authorizer,
    })



    // Api GateWay tv show resource
    tvShowApiResource := api.Root().AddResource(jsii.String("tvShow"), nil)
    tvShowApiResource.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(getShowLambda, nil), &awsapigateway.MethodOptions{
        // TODO: enable when frontend is done
        // AuthorizationType: awsapigateway.AuthorizationType_COGNITO,
        // Authorizer:        authorizer,
    })
    tvShowApiResource.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(postShowLambda, nil), &awsapigateway.MethodOptions{
        // TODO: enable when frontend is done
        // AuthorizationType: awsapigateway.AuthorizationType_COGNITO,
        // Authorizer:        authorizer,
    })

    return stack
}

func main() {
    app := awscdk.NewApp(nil)
    NewArgonStack(app, "ArgonStack", &awscdk.StackProps{ Env: &awscdk.Environment{
        Account: jsii.String(os.Getenv("ACCOUNT_ID")),
        Region: jsii.String("eu-central-1"),
    }, })
    app.Synth(nil)
}
