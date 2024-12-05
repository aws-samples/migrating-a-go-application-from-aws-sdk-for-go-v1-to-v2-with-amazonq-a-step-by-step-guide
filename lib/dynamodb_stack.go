package lib

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type DynamoDBStackProps struct {
	awscdk.StackProps
}

func NewDynamoDBStack(scope constructs.Construct, id string, props *DynamoDBStackProps) (awscdk.Stack, awslambda.Function, awslambda.Function) {
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
	lambdaPath := filepath.Join(path, "updatePlayer.zip")

	// Create a new DynamoDB table
	table := awsdynamodb.NewTable(stack, jsii.String("PlayerHits"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("player_id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	getHitsLambdaPath := filepath.Join(path, "hitCounter.zip")
	getHitsLambda := awslambda.NewFunction(stack, jsii.String("GetHitsFunction"), &awslambda.FunctionProps{
		Code:         awslambda.AssetCode_FromAsset(&getHitsLambdaPath, nil),
		Handler:      jsii.String("main"),
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Architecture: awslambda.Architecture_ARM_64(),
	})

	//Set dynamoDB table as an environment variable for the lambda function
	getHitsLambda.AddEnvironment(jsii.String("DYNAMODB_TABLE_NAME"), table.TableName(), nil)
	getHitsLambda.Role().AddManagedPolicy(awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonDynamoDBFullAccess")))

	// Create a new Lambda function to store hits by a player
	hitsLambda := awslambda.NewFunction(stack, jsii.String("PlayerHitsFunction"), &awslambda.FunctionProps{
		Code:         awslambda.AssetCode_FromAsset(&lambdaPath, nil),
		Handler:      jsii.String("main"),
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Architecture: awslambda.Architecture_ARM_64(),
	})

	//Set dynamoDB table as an environment variable for the lambda function
	hitsLambda.AddEnvironment(jsii.String("DYNAMODB_TABLE_NAME"), table.TableName(), nil)
	fmt.Println("policyDoc")
	policyDoc := map[string]interface{}{
		"Version": jsii.String("2012-10-17"),
		"Statement": []interface{}{
			map[string]interface{}{
				"Effect":   jsii.String("Allow"),
				"Action":   []*string{jsii.String("dynamodb:PutItem")},
				"Resource": []*string{jsii.String("*")},
			},
		},
	}
	customPolicyDocument := awsiam.PolicyDocument_FromJson(policyDoc)
	// Grant the lambda function permission to access the DynamoDB table
	// hitsLambda.Role().AddManagedPolicy(awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonDynamoDBFullAccess")))
	newPolicy := awsiam.NewPolicy(stack, jsii.String("PutItemPolicy"), &awsiam.PolicyProps{
		Document: customPolicyDocument,
	})
	hitsLambda.Role().AttachInlinePolicy(newPolicy)

	// Add read write privileges for the lambda function
	table.GrantReadWriteData(hitsLambda)
	table.GrantReadWriteData(getHitsLambda)

	return stack, hitsLambda, getHitsLambda
}
