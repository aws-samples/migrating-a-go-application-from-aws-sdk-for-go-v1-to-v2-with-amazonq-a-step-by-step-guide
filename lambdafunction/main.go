package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	lambdaHandler "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambdaService "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Player struct {
	PlayerID           string `json:"player_id"`
	LastName           string `json:"lastName"`
	FirstName          string `json:"firstname"`
	DOB                string `json:"dob"`
	Plays              string `json:"plays"`
	CountryOfBirth     string `json:"countryOfBirth"`
	CountryOfResidence string `json:"countryOfResidence"`
}

type playerWithHits struct {
	Player Player
	Hits   int
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Print the incoming request
	fmt.Printf("Received request: %v\n", request)

	bucketName := os.Getenv("BUCKET_NAME")
	fmt.Printf("Bucket name: %s\n", bucketName)

	hitsLambda := os.Getenv("HITS_LAMBDA")
	fmt.Printf("Hits Lambda: %s\n", hitsLambda)

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
			lastName := fields[1]
			firstName := fields[2]

			// Filter for specific first name or last name
			if strings.EqualFold(firstName, request.QueryStringParameters["firstName"]) || strings.EqualFold(lastName, request.QueryStringParameters["lastName"]) {
				fmt.Printf("Found a match for the player %v\n", line)
				player := Player{
					PlayerID:           fields[0],
					LastName:           fields[1],
					FirstName:          fields[2],
					DOB:                fields[3],
					Plays:              fields[4],
					CountryOfBirth:     fields[5],
					CountryOfResidence: fields[6],
				}
				// return this line as a response
				filteredPlayers = append(filteredPlayers, player)
				fmt.Printf("Calling hitslambda %v\n", line)
				// update struct with hits
				playerWithHits := playerWithHits{
					Player: player,
					Hits:   1,
				}
				fmt.Printf("Player with hits %v\n", playerWithHits)
				// Convert the filtered players to JSON
				jsonData2, err := json.Marshal(player)
				if err != nil {
					fmt.Printf("Error marshaling JSON: %v", err)
					ApiResponse.Body = "Error marshaling JSON"
					ApiResponse.StatusCode = 400
					return ApiResponse, nil
				}

				svc := lambdaService.New(sess)

				input := &lambdaService.InvokeInput{
					FunctionName:   &hitsLambda,
					InvocationType: aws.String("Event"),
					Payload:        jsonData2, // Replace with your desired payload,
				}

				_, err = svc.Invoke(input)
				if err != nil {
					fmt.Println("Error invoking Lambda function:", err)
					return ApiResponse, err
				}
				ApiResponse = events.APIGatewayProxyResponse{
					StatusCode: 200,
					Body:       string(jsonData2),
				}

			}
		}

	}

	// Response
	return ApiResponse, err

}

func main() {
	lambdaHandler.Start(HandleRequest)
}
