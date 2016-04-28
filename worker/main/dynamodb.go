package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var dynoClient *dynamodb.DynamoDB = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1").WithCredentials(creds))

func checkItem(tableName string, id string) bool {
	params := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}
	resp, err := dynoClient.GetItem(params)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if len(resp.Item) == 0 {
		return false
	}
	return true
}

func getItem(tableName string, id string) string{
	params := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}
	resp, err := dynoClient.GetItem(params)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	if len(resp.Item) == 0 {
		return ""
	}
	return *resp.Item["value"].S
}

func writeItem(tableName string, id string, value string) bool {
	params := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
			"value": {
				S: aws.String(value),
			},
		},
	}
	resp, err := dynoClient.PutItem(params)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	fmt.Println(resp)
	return true
}
