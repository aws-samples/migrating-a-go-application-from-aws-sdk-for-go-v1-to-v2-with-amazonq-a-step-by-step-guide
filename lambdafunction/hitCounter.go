package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type GetPlayer struct {
	PlayerID           int    `json:"player_id"`
	LastName           string `json:"lastName"`
	FirstName          string `json:"firstName"`
	DOB                string `json:"dob"`
	Plays              string `json:"plays"`
	CountryOfBirth     string `json:"countryOfBirth"`
	CountryOfResidence string `json:"countryOfResidence"`
	Hits               int    `json:"hits"`
}
type PlayerResult struct {
	PlayerID           string `json:"player_id"`
	LastName           string `json:"lastName"`
	FirstName          string `json:"firstName"`
	DOB                string `json:"dob"`
	Plays              string `json:"plays"`
	CountryOfBirth     string `json:"countryOfBirth"`
	CountryOfResidence string `json:"countryOfResidence"`
	Hits               int    `json:"hits"`
}

func HitCounter(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	tableName := os.Getenv("DYNAMODB_TABLE_NAME")

	ApiResponse := events.APIGatewayProxyResponse{}

	switch request.HTTPMethod {
	case "GET":
		var requestBody GetPlayer
		fmt.Printf("Received request: %v\n", request.QueryStringParameters["player_id"])
		var err = json.Unmarshal([]byte(request.QueryStringParameters["player_id"]), &requestBody.PlayerID)
		if err != nil {
			return ApiResponse, fmt.Errorf("error unmarshaling request body: %v", err)
		}

		input := &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"player_id": {
					S: aws.String(request.QueryStringParameters["player_id"]),
				},
			},
		}
		result, err := svc.GetItem(input)
		if err != nil {
			return ApiResponse, err
		}

		if result.Item == nil {
			msg := "could not find"
			return ApiResponse, errors.New(msg)
		}

		item := PlayerResult{}

		fmt.Printf("Item: %v\n", result.Item)

		err = dynamodbattribute.UnmarshalMap(result.Item, &item)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}
		jsonData, err := json.Marshal(item)
		if err != nil {
			fmt.Printf("Error marshaling JSON: %v", err)
			ApiResponse.Body = "Error marshaling JSON"
			ApiResponse.StatusCode = 400
			return ApiResponse, nil
		}

		// return ApiResponse as json
		ApiResponse = events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(jsonData),
		}
	}

	return ApiResponse, nil

}

func main() {
	lambda.Start(HitCounter)
}
