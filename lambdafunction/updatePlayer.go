package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type playerWithHits struct {
	PlayerID           string `json:"player_id"`
	LastName           string `json:"lastName"`
	FirstName          string `json:"firstName"`
	DOB                string `json:"dob"`
	Plays              string `json:"plays"`
	CountryOfBirth     string `json:"countryOfBirth"`
	CountryOfResidence string `json:"countryOfResidence"`
	Hits               int    `json:"hits"`
}

// Create the handler function and put and update player
func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Print the incoming request
	fmt.Printf("Received request: %v\n", request)
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")

	ApiResponse := events.APIGatewayProxyResponse{}

	switch request.HTTPMethod {
	case "POST":
		var requestBody playerWithHits
		var err = json.Unmarshal([]byte(request.Body), &requestBody)
		if err != nil {
			return ApiResponse, fmt.Errorf("Error unmarshaling request body: %v", err)
		}

		fmt.Printf("Received request: %v\n", requestBody.PlayerID)
		// playerId := requestBody
		err = UpdateHits(requestBody, tableName)
		if err != nil {
			return ApiResponse, err
		}
	}
	ApiResponse = events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Player updated successfully",
	}
	return ApiResponse, nil
}

func UpdateHits(requestBody playerWithHits, tableName string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	// Player to update is
	fmt.Println("The player to update is %v", requestBody.PlayerID)

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"player_id": {
				S: aws.String(requestBody.PlayerID),
			},
		},
		UpdateExpression: aws.String("ADD hits :incr"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":incr": {
				N: aws.String("1"),
			},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
	}

	_, err := svc.UpdateItem(input)
	return err
}

func main() {
	lambda.Start(HandleRequest)
}
