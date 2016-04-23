package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var creds *credentials.Credentials = credentials.NewEnvCredentials()
var svc *sqs.SQS = sqs.New(session.New(), aws.NewConfig().WithRegion("us-east-1").WithCredentials(creds))
var queueUrl string

func createQueue(queueName string) {
	paramsCreate := &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
	}
	_, err := svc.CreateQueue(paramsCreate)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	params := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}
	resp, _ := svc.GetQueueUrl(params)
	queueUrl = *resp.QueueUrl
}

func deleteQueue() {
	params := &sqs.DeleteQueueInput{
		QueueUrl: aws.String(queueUrl),
	}
	response, err := svc.DeleteQueue(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(response)
}

/**
function to send message to the specified queue
*/
func sendMessage(message string) bool {
	params := &sqs.SendMessageInput{
		MessageBody:  aws.String(message),
		QueueUrl:     aws.String(queueUrl),
		DelaySeconds: aws.Int64(1),
	}
	resp, err := svc.SendMessage(params)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	fmt.Println(resp)
	return true
}
