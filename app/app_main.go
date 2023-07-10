package app

import (
	"ddb-stream-latency/config"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AppMainStackProps struct {
	awscdk.StackProps
}

func NewAppMainStack(scope constructs.Construct, id string, props *AppMainStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	awsdynamodb.NewTable(stack, jsii.String("timer"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("timestamp"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode:   awsdynamodb.BillingMode_PROVISIONED,
		TableName:     jsii.String("timer"),
		ReadCapacity:  jsii.Number(config.TimerTableReadCapacity),
		WriteCapacity: jsii.Number(config.TimerTableWriteCapacity),
		Stream:        awsdynamodb.StreamViewType_NEW_IMAGE,
	})

	awslambda.NewFunction(stack, jsii.String("ddbreader"), &awslambda.FunctionProps{
		Runtime:    awslambda.Runtime_GO_1_X(),
		Handler:    jsii.String("bin/reader"),
		MemorySize: jsii.Number(128),
		Timeout:    awscdk.Duration_Seconds(jsii.Number(10)),
		Code:       awslambda.Code_FromAsset(jsii.String("app/ddbreader"), nil),
		CurrentVersionOptions: &awslambda.VersionOptions{
			RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
			RetryAttempts: jsii.Number(1),
		},
	})

	awsdynamodb.NewTable(stack, jsii.String("dashboard"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("timestamp"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode:   awsdynamodb.BillingMode_PROVISIONED,
		TableName:     jsii.String("dashboard"),
		ReadCapacity:  jsii.Number(config.DashboardTableReadCapacity),
		WriteCapacity: jsii.Number(config.DashboardTableWriteCapacity),
		Stream:        awsdynamodb.StreamViewType_NEW_IMAGE,
	})

	return stack
}
