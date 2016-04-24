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

/**
Start and ge the queue url
*/
func initialize(queueName string) {
	params := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}
	resp, err := svc.GetQueueUrl(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	queueUrl = *resp.QueueUrl
}

/**
Recieve message from the passed queue url
*/
func receiveMessage() (bool, string) {
	params := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(queueUrl),
	}
	resp, err := svc.ReceiveMessage(params)

	if err != nil {
		fmt.Println(err.Error())
		return false, ""
	}
	//fetching one message per call
	if len(resp.Messages) == 0 {
		fmt.Println("Empty Queue")
		return false, ""
	}
	message := *resp.Messages[0].Body
	receiptHandle := *resp.Messages[0].ReceiptHandle
	paramsDelete := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueUrl),
		ReceiptHandle: aws.String(receiptHandle),
	}
	_, e := svc.DeleteMessage(paramsDelete)
	if e != nil {
		fmt.Println(e.Error())
		return false, ""
	}
	return true, message
}

/**
function to send message to the specified queue
*/
func sendMessage(outputQueueUrl string) bool {
	params := &sqs.SendMessageInput{
		MessageBody:  aws.String("1"),
		QueueUrl:     aws.String(outputQueueUrl),
		DelaySeconds: aws.Int64(1),
	}
	_, err := svc.SendMessage(params)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	//fmt.Println(resp)
	return true
}
