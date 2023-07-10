package pipeline

import (
	"ddb-stream-latency/app"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
)

type PipelineMainStageProps struct {
	awscdk.StageProps
}

type PipelineManualApprovalStageProps struct {
	awscdk.StageProps
}

func NewPipelineMainStage(scope constructs.Construct, id string, props *PipelineMainStageProps) awscdk.Stage {
	var sprops awscdk.StageProps

	if props != nil {
		sprops = props.StageProps
	}

	stage := awscdk.NewStage(scope, &id, &sprops)

	app.NewAppMainStack(stage, "DDBStreamLatency-MainStack", nil)

	return stage
}
