// lambda-function/main.go
package lambdafunction

import (
    "context"
    "github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context) (string, error) {
    return "Hello, Amazon Q!", nil
}

func main() {
    lambda.Start(HandleRequest)
}
