# Welcome to your CDK Go project!

This is a project for CDK development with Go.

The `cdk.json` file tells the CDK toolkit how to execute your app.

## Useful commands

 * `GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap lambdafunction/main.go`                           Build the Lambda

 * `zip myFunction.zip bootstrap`   To zip the GO lambda
 * `cdk deploy`                     deploy this stack to your default AWS account/region , set AWS_REGION 
 * `cdk diff`                       compare deployed stack with current state
 * `cdk synth`                      emits the synthesized CloudFormation template
 * `cdk bootstrap`                  For the first time bootstrap in the account/region
 * `go test`                        run unit tests
