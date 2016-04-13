package main

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)
	
var svc *SQS

func init(queueName,a) {
	svc = sqs.New(session.New())	

}

func createQueue(queueName) {
	response, err := svc.CreateQueue()
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	fmt.Println(response)
}

func deleteQueue(queueUrl){
	params := &sqs.DeleteQueueInput{
		QueueUrl: aws.String(queueUrl), 
	}
	response, err := svc.DeleteQueue(params)

	if err != nil{
		fmt.Println(err.Error())
		return 
	}
	fmt.Println(response)
}