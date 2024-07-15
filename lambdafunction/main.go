// lambda-function/main.go
package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Print the incoming request
	fmt.Printf("Received request: %v\n", request)

	ApiResponse := events.APIGatewayProxyResponse{}
	// Switch for identifying the HTTP request
	switch request.HTTPMethod {
	case "GET":
		// Obtain the QueryStringParameter
		instanceId := request.QueryStringParameters["instanceId"]
		region := request.QueryStringParameters["region"]

		// Check if the name parameter is present
		if instanceId != "" {
			ApiResponse = events.APIGatewayProxyResponse{Body: "instanceId=" + instanceId + " Region= " + region, StatusCode: 200}
		} else {
			ApiResponse = events.APIGatewayProxyResponse{Body: "Error: Query Parameter name missing", StatusCode: 500}
		}
	}
	// Response
	return ApiResponse, nil

}

func main() {
	lambda.Start(HandleRequest)
}
