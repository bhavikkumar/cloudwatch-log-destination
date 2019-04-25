package logs_test

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/bhavikkumar/cloudwatch-log-destination/cloudwatch/logs"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNoDestinationArn(t *testing.T) {
	_, err := logs.DestinationArn()
	assert.Error(t, err)
}

func TestDestinationArn(t *testing.T) {
	os.Setenv("DESTINATION_ARN", "arn:test")
	actual, err := logs.DestinationArn()
	assert.Equal(t, actual, "arn:test")
	assert.Nil(t, err)
}

func TestNewCloudWatchLogFromEventWithNoRequestDetail(t *testing.T) {
	eventJson := `{"version": "0","id": "560646ad-c1f3-a8fb-ea6d-730c1b2bfd63","detail-type": "AWS API Call via CloudTrail","source": "aws.logs","account": "336840772780","time": "2019-04-10T19:18:47Z","region": "us-east-1","resources": []}`
	var event events.CloudWatchEvent
	json.Unmarshal([]byte(eventJson), &event)
	actual := logs.NewFromEvent(event)
	assert.Empty(t, actual.LogGroupName)
	assert.Empty(t, actual.DestinationArn)
}

func TestNewCloudWatchLogFromEventNoRequestParameters(t *testing.T) {
	eventJson := `{"version": "0","id": "560646ad-c1f3-a8fb-ea6d-730c1b2bfd63","detail-type": "AWS API Call via CloudTrail","source": "aws.logs","account": "336840772780","time": "2019-04-10T19:18:47Z","region": "us-east-1","resources": [], "detail": { "requestParameters":{}}}`
	var event events.CloudWatchEvent
	json.Unmarshal([]byte(eventJson), &event)
	actual := logs.NewFromEvent(event)
	assert.Empty(t, actual.LogGroupName)
	assert.Empty(t, actual.DestinationArn)
}

func TestNewCloudWatchLogFromMalformedEvent(t *testing.T) {
	eventJson := `{"version": "0","id": "560646ad-c1f3-a8fb-ea6d-730c1b2bfd63","detail-type": "AWS API Call via CloudTrail","source": "aws.logs","account": "336840772780","time": "2019-04-10T19:18:47Z","region": "us-east-1","resources": [], "detail": { "requestParameters": { "logGroupName": 1}}}`
	var event events.CloudWatchEvent
	json.Unmarshal([]byte(eventJson), &event)
	actual := logs.NewFromEvent(event)
	assert.Empty(t, actual.LogGroupName)
	assert.Empty(t, actual.DestinationArn)
}

func TestNewCloudWatchLogFromEvent(t *testing.T) {
	eventJson := `{"version": "0","id": "560646ad-c1f3-a8fb-ea6d-730c1b2bfd63","detail-type": "AWS API Call via CloudTrail","source": "aws.logs","account": "336840772780","time": "2019-04-10T19:18:47Z","region": "us-east-1","resources": [], "detail": { "requestParameters": { "logGroupName": "test", "destinationArn": "arn:test"}}}`
	var event events.CloudWatchEvent
	json.Unmarshal([]byte(eventJson), &event)
	actual := logs.NewFromEvent(event)
	expected := logs.CloudWatchLog{LogGroupName: "test", DestinationArn: "arn:test"}
	assert.Equal(t, expected, actual)
}

func TestUpdateSubscriptionFilter(t *testing.T) {
	cloudWatchLog := logs.CloudWatchLog{LogGroupName: "test", DestinationArn: "arn:test"}
	err := cloudWatchLog.UpdateSubscriptionFilter(mockCloudWatchLogs{}, "arn:test")
	assert.NoError(t, err)
}

func TestDeleteSubscriptionFilter(t *testing.T) {
	cloudWatchLog := logs.CloudWatchLog{LogGroupName: "test", DestinationArn: "arn:test"}
	err := cloudWatchLog.UpdateSubscriptionFilter(mockCloudWatchLogs{}, "arn:test")
	assert.NoError(t, err)
}

type mockCloudWatchLogs struct {
	cloudwatchlogsiface.CloudWatchLogsAPI
	PutSubResp    cloudwatchlogs.PutSubscriptionFilterOutput
	DeleteSubResp cloudwatchlogs.DeleteSubscriptionFilterOutput
}

func (m mockCloudWatchLogs) PutSubscriptionFilter(in *cloudwatchlogs.PutSubscriptionFilterInput) (*cloudwatchlogs.PutSubscriptionFilterOutput, error) {
	return &m.PutSubResp, nil
}

func (m mockCloudWatchLogs) DeleteSubscriptionFilter(in *cloudwatchlogs.DeleteSubscriptionFilterInput) (*cloudwatchlogs.DeleteSubscriptionFilterOutput, error) {
	return &m.DeleteSubResp, nil
}
