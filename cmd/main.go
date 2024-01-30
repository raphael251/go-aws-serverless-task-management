package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/raphael251/go-aws-serverless-task-management/internal/handlers"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {

	lambda.Start(Handler)
}

var ErrUnexpected = "unexpected error"

func Handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	region := os.Getenv("AWS_REGION")
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))

	if err != nil {
		return nil, err
	}

	dbClient := dynamodb.NewFromConfig(cfg)

	switch req.Path {
	case "/api/users":
		return UsersRouter(req, dbClient)
	}

	switch req.HTTPMethod {
	case "POST":
		return handlers.CreateProject(req, dbClient)
	case "GET":
		return handlers.FindAllProjects(req)
	case "PUT":
		return handlers.UpdateProject(req)
	case "DELETE":
		return handlers.DeleteProject(req)
	}
	return nil, errors.New(ErrUnexpected)
}

func UsersRouter(req events.APIGatewayProxyRequest, dbClient *dynamodb.Client) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "POST":
		return handlers.CreateUser(req, dbClient)
	default:
		return nil, errors.New(ErrUnexpected)
	}
}
