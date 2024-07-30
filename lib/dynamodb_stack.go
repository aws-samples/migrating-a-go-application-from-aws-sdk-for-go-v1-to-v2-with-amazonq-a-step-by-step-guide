package lib

import (
    "github.com/aws/aws-cdk-go/awscdk/v2"
    "github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
    "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
    "github.com/aws/constructs-go/constructs/v10"
    "github.com/aws/jsii-runtime-go"
	"log"
	"os" 
	"path/filepath"
)

type DynamoDBStackProps struct {
	awscdk.StackProps
}

func NewDynamoDBStack(scope constructs.Construct, id string, props *DynamoDBStackProps) (awscdk.Stack, awslambda.Function) {
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
	lambdaPath := filepath.Join(path, "hitsLambda.zip")

    // Create a new DynamoDB table
    table := awsdynamodb.NewTable(stack, jsii.String("PlayerHits"), &awsdynamodb.TableProps{
        PartitionKey: &awsdynamodb.Attribute{
            Name: jsii.String("player_id"),
            Type: awsdynamodb.AttributeType_STRING,
        },
        RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
    })

    // Create a new Lambda function to store hits by a player
	hitsLambda := awslambda.NewFunction(stack, jsii.String("PlayerHitsFunction"), &awslambda.FunctionProps{
		Code:         awslambda.AssetCode_FromAsset(&lambdaPath, nil),
		Handler:      jsii.String("main"),
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Architecture: awslambda.Architecture_ARM_64(),
	})

    //Set dynamoDB table as an environment variable for the lambda function
    hitsLambda.AddEnvironment(jsii.String("DYNAMODB_TABLE_NAME"), table.TableName(), nil);
	
    // Add read write privileges for the lambda function
    table.GrantReadWriteData(hitsLambda);

    return stack, hitsLambda
}
