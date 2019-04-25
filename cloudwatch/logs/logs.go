package logs

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	log "github.com/sirupsen/logrus"
	"os"
)

type requestParameters struct {
	CloudWatchLog CloudWatchLog `json:"requestParameters"`
}

type CloudWatchLog struct {
	LogGroupName   string `json:"logGroupName"`
	FilterName     string `json:"filterName"`
	DestinationArn string `json:"destinationArn"`
}

func NewFromEvent(event events.CloudWatchEvent) (cwLog CloudWatchLog) {
	cwLog = CloudWatchLog{}
	cwLog.parseCloudWatchEvent(event)
	return
}

func DestinationArn() (destinationArn string, err error) {
	destinationArn = os.Getenv("DESTINATION_ARN")
	if destinationArn == "" {
		err = errors.New("DESTINATION_ARN environment variable not set")
	}
	return
}

func (cwLog *CloudWatchLog) parseCloudWatchEvent(event events.CloudWatchEvent) {
	if len(event.Detail) <= 0 {
		log.WithFields(log.Fields{"id": event.Version, "detailType": event.DetailType, "source": event.Source}).Warn("CloudWatch Event missing detail section")
		return
	}

	var requestParameters requestParameters
	err := json.Unmarshal(event.Detail, &requestParameters)
	if err != nil {
		log.WithFields(log.Fields{"detail": event.Detail}).Warn("Could not parse CloudWatch Event details", err)
		return
	}
	cwLog.LogGroupName = requestParameters.CloudWatchLog.LogGroupName
	cwLog.DestinationArn = requestParameters.CloudWatchLog.DestinationArn
}

func (cwLog *CloudWatchLog) UpdateSubscriptionFilter(client cloudwatchlogsiface.CloudWatchLogsAPI, destinationArn string) (err error) {
	var filterName = "DefaultLogDestination"
	var filterPattern = ""
	var distribution = cloudwatchlogs.DistributionByLogStream
	input := &cloudwatchlogs.PutSubscriptionFilterInput{
		LogGroupName:   &cwLog.LogGroupName,
		DestinationArn: &destinationArn,
		FilterPattern:  &filterPattern,
		FilterName:     &filterName,
		Distribution:   &distribution,
	}
	_, err = client.PutSubscriptionFilter(input)
	return
}

func (cwLog *CloudWatchLog) DeleteSubscriptionFilter(client cloudwatchlogsiface.CloudWatchLogsAPI, filterName string) (err error) {
	input := &cloudwatchlogs.DeleteSubscriptionFilterInput{
		LogGroupName: &cwLog.LogGroupName,
		FilterName:   &filterName,
	}
	_, err = client.DeleteSubscriptionFilter(input)
	return
}
