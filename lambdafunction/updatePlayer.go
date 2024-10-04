package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type playerWithHitsBody struct {
	PlayerID           string `json:"player_id"`
	LastName           string `json:"lastName"`
	FirstName          string `json:"firstName"`
	DOB                string `json:"dob"`
	Plays              string `json:"plays"`
	CountryOfBirth     string `json:"countryOfBirth"`
	CountryOfResidence string `json:"countryOfResidence"`
	Hits               int    `json:"hits"`
}

// HandlePlayerRequest Create the handler function and put and update player
func HandlePlayerRequest(request playerWithHitsBody) (string, error) {

	// Print the incoming request
	fmt.Printf("Received request: %v\n", request)
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")

	result, err := UpdateHits(request, tableName)

	if err != nil {
		return "", err
	}
	fmt.Printf("The result is %v\n", result)

	return result, err
}

func UpdateHits(requestBody playerWithHitsBody, tableName string) (string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	// Player to update is
	fmt.Printf("The player to update is %v\n", requestBody.PlayerID)

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"player_id": {
				S: aws.String(requestBody.PlayerID),
			},
		},
		UpdateExpression: aws.String("SET lastName = :l, firstName = :f, dob = :dob, plays = :plays, countryOfBirth = :cob, countryOfResidence = :cor, hits = if_not_exists(hits, :zero) + :incr"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":l": {
				S: aws.String(requestBody.LastName),
			},
			":f": {
				S: aws.String(requestBody.FirstName),
			},
			":dob": {
				S: aws.String(requestBody.DOB),
			},
			":plays": {
				S: aws.String(requestBody.Plays),
			},
			":cob": {
				S: aws.String(requestBody.CountryOfBirth),
			},
			":cor": {
				S: aws.String(requestBody.CountryOfResidence),
			},
			":incr": {
				N: aws.String("1"),
			},
			":zero": {
				N: aws.String("0"),
			},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
	}
	fmt.Printf("The input is %v\n", input)

	output, err := svc.UpdateItem(input)
	if err != nil {
		return "", err
	}
	fmt.Printf("The output is %v\n", output)

	// Convert the updated item attributes to JSON
	updatedItem := map[string]*dynamodb.AttributeValue{}

	for k, v := range output.Attributes {
		updatedItem[k] = v
	}
	fmt.Printf("The updated item uncoded is %v\n", updatedItem)
	jsonBytes, err := json.Marshal(updatedItem)
	if err != nil {
		return "", err
	}
	fmt.Printf("The updated item is %v\n", string(jsonBytes))

	return string(jsonBytes), nil

}

func main() {
	lambda.Start(HandlePlayerRequest)
}
