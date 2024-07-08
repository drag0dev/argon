package main

import (
	"common"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)


var cognitoUserPoolId string
var cognitoClient *cognitoidentityprovider.Client
var jwksUrl string


func handleRequest(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
    token := strings.TrimPrefix(event.AuthorizationToken, "Bearer ")

    // Fetch the JWKS
    set, err := jwk.Fetch(ctx, jwksUrl)
    if (err != nil) {
        log.Println(err)
        return denyResponse("user", event.MethodArn), nil
    }

    // parse token
    parsedToken, err := jwt.Parse(
        []byte(token),
        jwt.WithKeySet(set),
        jwt.WithValidate(true),
    )
    if (err != nil) {
        log.Println(err)
        return denyResponse("user", event.MethodArn), nil
    }

    sub, ok := parsedToken.Get("sub")
    if !ok {
        log.Println("sub claim not found in token")
        return denyResponse("user", event.MethodArn), nil
    }

    principalID, ok := sub.(string)
    if !ok {
        log.Println("sub claim is not string")
        return denyResponse("user", event.MethodArn), nil
    }

    usernameTemp, ok := parsedToken.Get("cognito:username")
    if !ok {
        log.Println("username claim not found in token")
        return denyResponse("user", event.MethodArn), nil
    }

    username, ok := usernameTemp.(string)
    if !ok {
        log.Println("username claim is not string")
        return denyResponse("user", event.MethodArn), nil
    }

    log.Printf("Sub: %s", sub)
    log.Printf("Principal: %s", principalID)
    log.Printf("Username: %s\n", username)

    userGroups, err := cognitoClient.AdminListGroupsForUser(ctx, &cognitoidentityprovider.AdminListGroupsForUserInput{
        UserPoolId: aws.String(cognitoUserPoolId),
        Username:   aws.String(username),
    })

    if err != nil {
        log.Printf("Error getting user info: %v\n", err)
        return denyResponse(principalID, event.MethodArn), nil
    }

    isAdmin := false
    fmt.Println(len(userGroups.Groups))
    for _, group := range userGroups.Groups {
        if (*group.GroupName == common.AdminGroupName) {
            isAdmin = true
            break
        }
    }

    if (!isAdmin) {
        log.Printf("user is not admin")
        return denyResponse(principalID, event.MethodArn), nil
    }

    // generate policy
    authResponse := events.APIGatewayCustomAuthorizerResponse{ PrincipalID: principalID }
    authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
        Version: "2012-10-17",
        Statement: []events.IAMPolicyStatement{
            {
                Action:   []string{"execute-api:Invoke"},
                Effect:   "Allow",
                Resource: []string{event.MethodArn},
            },
        },
    }
    return authResponse, nil
}

func denyResponse(principalId string, resource string) events.APIGatewayCustomAuthorizerResponse {
    authResponse := events.APIGatewayCustomAuthorizerResponse{ PrincipalID: principalId }
    authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
        Version: "2012-10-17",
        Statement: []events.IAMPolicyStatement{
            {
                Action:   []string{"execute-api:Invoke"},
                Effect:   "Deny",
                Resource: []string{resource},
            },
        },
    }
    return authResponse
}

func main() {
    sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
    if err != nil {
        log.Fatal("Cannot load in default config")
    }
    cognitoUserPoolId = os.Getenv("COGNITO_USER_POOL_ID")
    jwksUrl = fmt.Sprintf("https://cognito-idp.eu-central-1.amazonaws.com/%s/.well-known/jwks.json", cognitoUserPoolId)
    cognitoClient = cognitoidentityprovider.NewFromConfig(sdkConfig)
    lambda.Start(handleRequest)
}
