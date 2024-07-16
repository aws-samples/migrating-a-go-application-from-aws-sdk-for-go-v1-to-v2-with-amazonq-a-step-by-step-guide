package main

import (
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	// "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/jsii-runtime-go"
)

// example tests. To run these tests, uncomment this file along with the
// example resource in go_sdk_amazon_q_test.go
func TestGoSdkAmazonQStack(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewGoSdkWithAmazonQDemoStack(app, "MyStack", nil)

	// THEN
	template := assertions.Template_FromStack(stack, nil)

	template.HasResourceProperties(jsii.String("AWS::Lambda::Function"), map[string]any{
		"Runtime": "provided.al2023",
	})

	// Check for api gw
	template.HasResourceProperties(jsii.String("AWS::ApiGateway::RestApi"), map[string]any{
		"Name": "Endpoint",
	})

	// check for S3 bucket
	template.HasResourceProperties(jsii.String("AWS::S3::Bucket"), map[string]any{
		"BucketName": "my-bucket-20240716",
	})

}
