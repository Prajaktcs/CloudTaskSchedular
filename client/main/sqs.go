package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"strings"
)

var creds *credentials.Credentials = credentials.NewEnvCredentials()
var svc *sqs.SQS = sqs.New(session.New(), aws.NewConfig().WithRegion("us-east-1").WithCredentials(creds))
var queueUrl string

func createQueue(queueName string) string {
	paramsCreate := &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
	}
	_, err := svc.CreateQueue(paramsCreate)
	if err != nil {
		//fmt.Println(err.Error())
		return ""
	}
	params := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}
	resp, _ := svc.GetQueueUrl(params)
	if strings.Compare(queueName, "output") != 0 {
		queueUrl = *resp.QueueUrl
		return ""
	} else {
		return *resp.QueueUrl
	}
}

func deleteQueue() {
	params := &sqs.DeleteQueueInput{
		QueueUrl: aws.String(queueUrl),
	}
	response, err := svc.DeleteQueue(params)

	if err != nil {
		//fmt.Println(err.Error())
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
	_, err := svc.SendMessage(params)
	if err != nil {
		//fmt.Println(err.Error())
		return false
	}

	return true
}

/**
Recieve message from the passed queue url
*/
func receiveMessage(outputQueueUrl string) (bool, string) {
	params := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(outputQueueUrl),
	}
	resp, err := svc.ReceiveMessage(params)

	if err != nil {
		//fmt.Println(err.Error())
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
		QueueUrl:      aws.String(outputQueueUrl),
		ReceiptHandle: aws.String(receiptHandle),
	}
	_, e := svc.DeleteMessage(paramsDelete)
	if e != nil {
		//fmt.Println(e.Error())
		return false, ""
	}
	return true, message
}
