package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error loading AWS SDK configuration:", err)
		return
	}

	client := dynamodb.NewFromConfig(cfg)

	tableName := "timer"

	// Get the write interval and duration from command-line arguments
	writeInterval, duration, routine := getInput()

	endTime := time.Now().Add(duration)

	log := logrus.New()

	for i := 0; i < routine; i++ {
		go func(goroutineID int) {
			for time.Now().Before(endTime) {
				currentTime := time.Now().Format("2006-01-02 15:04:05.000")

				item := map[string]types.AttributeValue{
					"timestamp":    &types.AttributeValueMemberS{Value: currentTime},
					"goroutine_id": &types.AttributeValueMemberN{Value: strconv.Itoa(goroutineID)},
				}

				input := &dynamodb.PutItemInput{
					Item:      item,
					TableName: &tableName,
				}

				_, err := client.PutItem(context.TODO(), input)

				if err != nil {
					log.Error("Error writing to DynamoDB:", err)
				} else {
					log.Info("GoroutineID: ", goroutineID, "\tTimestamp written to DynamoDB: ", currentTime)
				}

				time.Sleep(time.Duration(writeInterval) * time.Second)
			}
		}(i)
	}

	time.Sleep(duration + time.Second)

	log.Info("All goroutines completed execution.")
}

func getInput() (int, time.Duration, int) {
	defaultInterval := 1      // Default interval in seconds
	defaultDuration := 60     // Default duration in minute
	defaultRoutineNumber := 1 //Default number of goroutines

	intervalPtr := flag.Int("t", defaultInterval, "Write interval in seconds")
	durationPtr := flag.Int("d", defaultDuration, "Duration in minutes")
	routinePtr := flag.Int("r", defaultRoutineNumber, "Number of goroutines")

	flag.Parse()

	writeInterval := *intervalPtr

	log := logrus.New()

	if writeInterval <= 0 {
		log.Error("Invalid write interval provided. Using default interval of", defaultInterval, "second(s).")
		writeInterval = defaultInterval
	}

	duration := time.Duration(*durationPtr) * time.Minute

	if duration <= 0 {
		log.Error("Invalid duration provided. Using default duration of", defaultDuration, "minute(s).")
		duration = time.Duration(defaultDuration) * time.Minute
	}

	routine := *routinePtr

	if routine <= 0 {
		log.Error("Invalid number of routines provided. Using default number of", defaultRoutineNumber, "goroutine(s).")
		routine = defaultRoutineNumber
	}

	return writeInterval, duration, routine
}
