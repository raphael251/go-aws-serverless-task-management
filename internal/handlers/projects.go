package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/raphael251/go-aws-serverless-task-management/internal/database"
)

var ErrUnexpected = "unexpected error"

type CreateProjectInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func CreateProject(req events.APIGatewayProxyRequest, dbClient *dynamodb.Client) (*events.APIGatewayProxyResponse, error) {
	var project *CreateProjectInput

	err := json.Unmarshal([]byte(req.Body), &project)

	if err != nil {
		return nil, err
	}

	output, err := dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"my-partition-key": &types.AttributeValueMemberS{Value: uuid.New().String()},
			"title":            &types.AttributeValueMemberS{Value: project.Title},
			"description":      &types.AttributeValueMemberS{Value: project.Description},
		},
		TableName: &database.AppTableName,
	})

	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully inserted an item in dynamoDB", output)

	type SuccessResponse struct {
		Message string `json:"message"`
	}

	response, err := json.Marshal(&SuccessResponse{Message: "success"})

	if err != nil {
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(response),
	}, nil
}

func FindAllProjects(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return nil, errors.New(ErrUnexpected)
}

func UpdateProject(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return nil, errors.New(ErrUnexpected)
}

func DeleteProject(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return nil, errors.New(ErrUnexpected)
}
