package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"strings"
)

type Player struct {
	LastName           string
	FirstName          string
	DOB                string
	Plays              string
	CountryOfBirth     string
	CountryOfResidence string
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Print the incoming request
	fmt.Printf("Received request: %v\n", request)

	bucketName := os.Getenv("BUCKET_NAME")
	fmt.Printf("Bucket name: %s\n", bucketName)

	ApiResponse := events.APIGatewayProxyResponse{}
	// Switch for identifying the HTTP request
	switch request.HTTPMethod {
	case "GET":
		// Obtain the QueryStringParameter
		playerFirstName := request.QueryStringParameters["firstName"]

		// Check if the name parameter is present
		if playerFirstName == "" {
			ApiResponse = events.APIGatewayProxyResponse{Body: "Error: Please enter firstName or firsName and LastName to lookup Players", StatusCode: 400}
		}
	}

	// Query S3 bucket to read a CSV file from prefix "allplayers"
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	downloader := s3manager.NewDownloader(sess)
	buffer := aws.NewWriteAtBuffer([]byte{})

	_, err := downloader.Download(buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("activeplayers/players.csv"), // Replace with the actual file name
	})

	if err != nil {
		fmt.Printf("Error downloading file: %v", err)
		ApiResponse.Body = "Error downloading file"
		ApiResponse.StatusCode = 500
		return ApiResponse, nil
	}
	var filteredPlayers []Player

	// Process the CSV data
	csvData := strings.Split(string(buffer.Bytes()), "\n")
	for _, line := range csvData {
		// Yamamoto,Sakura,1991-04-10,Right,Japan,Japan
		fmt.Println(line)
		fields := strings.Split(line, ",")
		// filter for first name or last name and return the data

		if len(fields) >= 6 { // Ensure we have at least six fields
			lastName := fields[0]
			firstName := fields[1]

			// Filter for specific first name or last name
			if strings.EqualFold(firstName, request.QueryStringParameters["firstName"]) || strings.EqualFold(lastName, request.QueryStringParameters["lastName"]) {
				fmt.Println(line)
				// return this line as a response
				filteredPlayers = append(filteredPlayers, Player{
					LastName:           fields[0],
					FirstName:          fields[1],
					DOB:                fields[2],
					Plays:              fields[3],
					CountryOfBirth:     fields[4],
					CountryOfResidence: fields[5],
				})
			}
		}

		// Convert the filtered players to JSON
		jsonData, err := json.Marshal(filteredPlayers)
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

	// Response
	return ApiResponse, nil

}

func main() {
	lambda.Start(HandleRequest)
}
