package pipeline

import (
	"ddb-stream-latency/config"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/pipelines"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func NewPipelineStack(scope constructs.Construct, id string, props *awscdk.StackProps) awscdk.Stack {

	stack := awscdk.NewStack(scope, &id, props)

	pipeline := pipelines.NewCodePipeline(stack, jsii.String(config.PipelineName), &pipelines.CodePipelineProps{
		Synth: pipelines.NewCodeBuildStep(jsii.String("SynthStep"), &pipelines.CodeBuildStepProps{
			Input: pipelines.CodePipelineSource_Connection(jsii.String(config.RepoName), jsii.String(config.RepoBranch), &pipelines.ConnectionSourceOptions{
				ConnectionArn:        jsii.String(config.ConnectionArn),
				CodeBuildCloneOutput: jsii.Bool(true),
				TriggerOnPush:        jsii.Bool(true),
			}),
			Commands: jsii.Strings(
				"npx npm install -g aws-cdk",
				"goenv install 1.18.3",
				"npx cdk synth"),
		}),
	})

	deploy := NewPipelineMainStage(stack, "DeployStage", nil)

	pipeline.AddStage(deploy, nil)

	return stack
}
