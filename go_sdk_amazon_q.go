package main

import (
	"fmt"
	"go_sdk_amazon_q/lib"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	s3deploy "github.com/aws/aws-cdk-go/awscdk/v2/awss3deployment" // Added for deployment
	assets "github.com/aws/aws-cdk-go/awscdk/v2/awss3assets" // Added for deployment
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdklabs/cdk-nag-go/cdknag/v2"
)

type GoSdkAmazonQStackProps struct {
	awscdk.StackProps
	HitsLambda    awslambda.Function
	getHitsLambda awslambda.Function
}

func NewGoSdkWithAmazonQDemoStack(scope constructs.Construct, id string, props *GoSdkAmazonQStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	hitsLambda := props.HitsLambda

	// The code that defines your stack goes here
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	lambdaPath := filepath.Join(path, "main.zip")

	// Create the Lambda function
	lambdaFunction := awslambda.NewFunction(stack, jsii.String("MyLambdaFunction"), &awslambda.FunctionProps{
		Code:         awslambda.AssetCode_FromAsset(&lambdaPath, nil),
		Handler:      jsii.String("main"),
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Architecture: awslambda.Architecture_ARM_64(),
	})

	lambdaRole := lambdaFunction.Role()

	invokePermissionStatement := awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect:    awsiam.Effect_ALLOW,
		Actions:   jsii.Strings("lambda:InvokeFunction"),
		Resources: jsii.Strings(*hitsLambda.FunctionArn()),
	})

	lambdaRole.AddToPrincipalPolicy(invokePermissionStatement)

	// Create the API Gateway

	api := awsapigateway.NewRestApi(stack, jsii.String("HitsAPI"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String("HitsAPI"),
		Description: jsii.String("This is my HitsAPI"),
	})

	getPlayers := api.Root().AddResource(jsii.String("getPlayers"), nil)
	getHits := api.Root().AddResource(jsii.String("getHits"), nil)

	getPlayers.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(lambdaFunction, &awsapigateway.LambdaIntegrationOptions{}), &awsapigateway.MethodOptions{
		MethodResponses: &[]*awsapigateway.MethodResponse{
			{StatusCode: jsii.String("200")},
		},
	})
	getHits.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(props.getHitsLambda, &awsapigateway.LambdaIntegrationOptions{}), &awsapigateway.MethodOptions{
		MethodResponses: &[]*awsapigateway.MethodResponse{
			{StatusCode: jsii.String("200")},
		},
	})

	// Create the S3 bucket
	bucket := awss3.NewBucket(stack, jsii.String("amazonqgosdk"), &awss3.BucketProps{
		//BucketName:    jsii.String(bucketName), // Convert bucketName to *string
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		AutoDeleteObjects: jsii.Bool(true),
	})
	
     // Specify asset options
	 assetOptions := &assets.AssetOptions{
        // Example: exclude certain files from the asset bundle
        Exclude: jsii.Strings("*.tmp", "*.log"),
	 }

    // Deploy files to the S3 bucket
    s3deploy.NewBucketDeployment(stack, jsii.String("DeployFiles"), &s3deploy.BucketDeploymentProps{
        Sources: &[]s3deploy.ISource{s3deploy.Source_Asset(jsii.String("./activeplayers"), assetOptions)},
        DestinationBucket: bucket,
        DestinationKeyPrefix: jsii.String("activeplayers/"), // Optional: folder in the bucket
    })
	
					
	fmt.Println("Bucket Name:", *bucket.BucketName())

	// Set a bucket name environment variable in lambda function and add permissions to lambda to read from S3 bucket
	lambdaFunction.AddEnvironment(jsii.String("BUCKET_NAME"), bucket.BucketName(), nil)
	lambdaFunction.AddEnvironment(jsii.String("HITS_LAMBDA"), hitsLambda.FunctionName(), nil)

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

	_, hitsLambda, getHitsLambda := lib.NewDynamoDBStack(app, "DynamoDBStack", &lib.DynamoDBStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	stack := NewGoSdkWithAmazonQDemoStack(app, "GoSdkAmazonQStack", &GoSdkAmazonQStackProps{
		awscdk.StackProps{
			Env: env(),
		},
		hitsLambda,
		getHitsLambda,
	})

	// Add AWS Solutions checks
	cdknag.NagSuppressions_AddStackSuppressions(stack, &[]*cdknag.NagPackSuppression{
		{
			Id:     jsii.String("AwsSolutions-IAM4"),
			Reason: jsii.String("Using AWS managed policies is acceptable for this demo"),
		},
		{
			Id:     jsii.String("AwsSolutions-IAM5"),
			Reason: jsii.String("Wildcard permissions are acceptable for this demo"),
		},
		{
			Id:     jsii.String("AwsSolutions-APIG1"),
			Reason: jsii.String("API Gateway logging not required for this demo"),
		},
		{
			Id:     jsii.String("AwsSolutions-APIG2"),
			Reason: jsii.String("Request validation not needed for this demo"),
		},
		{
			Id:     jsii.String("AwsSolutions-APIG3"),
			Reason: jsii.String("WAF not required for this demo"),
		},
		{
			Id:     jsii.String("AwsSolutions-APIG4"),
			Reason: jsii.String("Authorization will be implemented later"),
		},
		{
			Id:     jsii.String("AwsSolutions-APIG6"),
			Reason: jsii.String("CloudWatch logging for all methods not required for this demo"),
		},
		{
			Id:     jsii.String("AwsSolutions-COG4"),
			Reason: jsii.String("Cognito user pool will be implemented later"),
		},
		{
			Id:     jsii.String("AwsSolutions-S1"),
			Reason: jsii.String("S3 server access logging not required for this demo"),
		},
		{
			Id:     jsii.String("AwsSolutions-S10"),
			Reason: jsii.String("SSL requirement will be implemented later"),
		},
		// Add suppression for AwsSolutions-L1: The non-container Lambda function is not configured to use the latest runtime version.
		{
			Id:     jsii.String("AwsSolutions-L1"),
			Reason: jsii.String("The non-container Lambda function is not configured to use the latest runtime version."),
		},
		
	}, jsii.Bool(true))

	// Add the AWS Solutions Checks
	awscdk.Aspects_Of(stack).Add(cdknag.NewAwsSolutionsChecks(nil))

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
