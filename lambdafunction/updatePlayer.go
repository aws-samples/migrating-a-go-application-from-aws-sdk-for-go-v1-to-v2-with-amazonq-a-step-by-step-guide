package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
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
func HandleRequest(playerWithHits playerWithHits) (string, error) {

	// Print the incoming request
	fmt.Printf("Received request: %v\n", playerWithHits)
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")

	// Update the hits for the player
	err := UpdateHits(playerWithHits.PlayerID, tableName)
	if err != nil {
		return "", err
	}

	// return the player who was updated
	return "Player updated successfully", nil
}

func UpdateHits(playerID string, tableName string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	// Player to update is
	fmt.Println("The player to update is %v", playerID)

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"player_id": {
				S: aws.String(playerID),
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
