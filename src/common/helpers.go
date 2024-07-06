package common

import (
	"encoding/json"
	"log"
	"github.com/aws/aws-lambda-go/events"
)

type ErrorMessage struct {
    Message string `json:"message"`
}

func EmptyErrorResponse(statusCode int) events.APIGatewayProxyResponse {
    return events.APIGatewayProxyResponse{
        StatusCode: statusCode,
        Headers: map[string]string{
            "Access-Control-Allow-Origin": "*",
            "Access-Control-Allow-Credentials": "true",
            "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
            "Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
        },
    }
}

func ErrorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
    msg := ErrorMessage{Message: message}
    msgMarashled, err := json.Marshal(msg)
    if err != nil { log.Println("MARSHALING ERROR MESSAGE ERRORED???") }
    return events.APIGatewayProxyResponse{
        StatusCode: statusCode,
        Headers: map[string]string{
            "Content-Type": "application/json",
            "Access-Control-Allow-Origin": "*",
            "Access-Control-Allow-Credentials": "true",
            "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
            "Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
        },
        Body: string(msgMarashled),
    }
}
