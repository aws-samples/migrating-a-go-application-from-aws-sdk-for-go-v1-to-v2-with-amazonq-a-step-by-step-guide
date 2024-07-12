// lambda-function/main.go
package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
)


type MyEvent struct {
	Name string `json:"What is your name?"`
	Age  int    `json:"How old are you?"`
}

type MyResponse struct {
	Message string `json:"Answer"`
}

func HandleRequest(ctx context.Context, event *MyEvent) (*MyResponse, error) {
	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}

	return &MyResponse{Message: fmt.Sprintf("%s is %d years old!", event.Name, event.Age)}, nil

}

func main() {
	lambda.Start(HandleRequest)
}