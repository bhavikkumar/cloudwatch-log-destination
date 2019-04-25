package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/bhavikkumar/cloudwatch-log-destination/cloudwatch/logs"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(event events.CloudWatchEvent) error {
	cwLog := logs.NewFromEvent(event)
	defaultDestinationArn, err := logs.DestinationArn()
	if err != nil {
		return err
	}

	session := session.Must(session.NewSession())
	if cwLog.DestinationArn != defaultDestinationArn && cwLog.FilterName != "" {
		err = cwLog.DeleteSubscriptionFilter(cloudwatchlogs.New(session), cwLog.FilterName)
		if err != nil {
			log.Warn("Could not delete subscription filter", err)
			return err
		}
	}

	if cwLog.DestinationArn != defaultDestinationArn {
		return cwLog.UpdateSubscriptionFilter(cloudwatchlogs.New(session), defaultDestinationArn)
	}
	return nil
}
