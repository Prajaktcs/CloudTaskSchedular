package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var dynoClient *dynamodb.DynamoDB = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1").WithCredentials(creds))

func createTable(tableName string) {
	params := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(16),
			WriteCapacityUnits: aws.Int64(16),
		},
	}
	_, err := dynoClient.CreateTable(params)
	if err != nil {
		//fmt.Println(err.Error())
	}
	fmt.Println("Was called")
}

func deleteTable(tableName string) {
	params := &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	}
	_, err := dynoClient.DeleteTable(params)
	if err != nil {
		//fmt.Println(err.Error())
	}
}
