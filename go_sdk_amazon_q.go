package main

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"log"
	"os"
	"path/filepath"
	"go_sdk_amazon_q/lib"
)

type GoSdkAmazonQStackProps struct {
	awscdk.StackProps
}

func NewGoSdkWithAmazonQDemoStack(scope constructs.Construct, id string, props *GoSdkAmazonQStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	lambdaPath := filepath.Join(path, "myFunction.zip")

	// Create the Lambda function
	lambdaFunction := awslambda.NewFunction(stack, jsii.String("MyLambdaFunction"), &awslambda.FunctionProps{
		Code:         awslambda.AssetCode_FromAsset(&lambdaPath, nil),
		Handler:      jsii.String("main"),
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Architecture: awslambda.Architecture_ARM_64(),
	})

	if lambdaFunction.FunctionArn() != nil{
		fmt.Println("Lambda function ARN:", *lambdaFunction.FunctionArn())
	}

	// Create the API Gateway
	awsapigateway.NewLambdaRestApi(stack, jsii.String("Endpoint"), &awsapigateway.LambdaRestApiProps{
		Handler: lambdaFunction,
	})

	// create a s3 bucket
	// Format the date as a string in the desired format
	bucketName := fmt.Sprintf("my-bucket-20240716")

	// Print the bucket name
	fmt.Println("Bucket Name:", bucketName)

	// Create an S3 bucket
	bucket := awss3.NewBucket(stack, jsii.String("jrtestaccessbucket"), &awss3.BucketProps{
		BucketName:    jsii.String(bucketName), // Convert bucketName to *string
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	// Set a bucket name environment variable in lambda function and add permissions to lambda to read from S3 bucket
	lambdaFunction.AddEnvironment(jsii.String("BUCKET_NAME"), bucket.BucketName(), nil)

	bucket.GrantRead(lambdaFunction, nil)

	// Output the bucket name
	awscdk.NewCfnOutput(stack, jsii.String("BucketName"), &awscdk.CfnOutputProps{
		Value: bucket.BucketName(),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewGoSdkWithAmazonQDemoStack(app, "GoSdkAmazonQStack", &GoSdkAmazonQStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	lib.NewDynamoDBStack(app, "DynamoDBStack", &lib.DynamoDBStackProps{
        awscdk.StackProps{
            Env: env(),
        },
    })

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
