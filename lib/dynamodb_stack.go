package lib

import (
    "github.com/aws/aws-cdk-go/awscdk/v2"
    "github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
    "github.com/aws/constructs-go/constructs/v10"
    "github.com/aws/jsii-runtime-go"
)

type DynamoDBStackProps struct {
	awscdk.StackProps
}

func NewDynamoDBStack(scope constructs.Construct, id string, props *DynamoDBStackProps) awscdk.Stack {
    var sprops awscdk.StackProps
    if props != nil {
        sprops = props.StackProps
    }
    stack := awscdk.NewStack(scope, &id, &sprops)

    // Create a new DynamoDB table
    awsdynamodb.NewTable(stack, jsii.String("PlayerHits"), &awsdynamodb.TableProps{
        PartitionKey: &awsdynamodb.Attribute{
            Name: jsii.String("id"),
            Type: awsdynamodb.AttributeType_STRING,
        },
        RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
    })

	

    return stack
}
