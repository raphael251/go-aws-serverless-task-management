package handlers

import (
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var ErrUnexpected = "unexpected error"

type CreateProjectInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func CreateProject(req events.APIGatewayProxyRequest, dbClient *dynamodb.Client) (*events.APIGatewayProxyResponse, error) {
	return nil, errors.New(ErrUnexpected)
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
