package main

import (
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sirupsen/logrus"
)

func handleReqeust(ctx context.Context, event events.DynamoDBEvent) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		logrus.Errorf("failed to load AWS config: %v", err)
		return err
	}

	svc := dynamodb.NewFromConfig(cfg)

	for _, record := range event.Records {
		if record.EventName == "INSERT" || record.EventName == "MODIFY" {
			newImage := record.Change.NewImage

			writertimestamp := newImage["timestamp"].String()
			writergoroutineID := newImage["goroutine_id"].String()

			currentTime := time.Now()

			writerTime, err := time.Parse("2006-01-02 15:04:05.000", writertimestamp)

			if err != nil {
				logrus.Errorf("failed to parse writertimestamp: %v", err)
				return err
			}

			timeDiff := currentTime.Sub(writerTime).Milliseconds()

			item := map[string]types.AttributeValue{
				"timestamp":          &types.AttributeValueMemberS{Value: currentTime.Format("2006-01-02 15:04:05.000")},
				"writertimestamp":    &types.AttributeValueMemberS{Value: writertimestamp},
				"writergoroutine_id": &types.AttributeValueMemberS{Value: writergoroutineID},
				"time_diff_ms":       &types.AttributeValueMemberN{Value: strconv.FormatInt(timeDiff, 10)},
			}

			_, err = svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
				TableName: aws.String("dashboard"),
				Item:      item,
			})

			if err != nil {
				logrus.Errorf("failed to put item into dashboard table: %v", err)
				return err
			}
		}
	}
	return nil
}
func main() {
	lambda.Start(handleReqeust)
}
