package main

import (
	"context"
	"errors"
	"net/http"
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
	case "/api/projects":
		return ProjectsRouter(req, dbClient)
	default:
		return nil, errors.New(ErrUnexpected)
	}
}

func UsersRouter(req events.APIGatewayProxyRequest, dbClient *dynamodb.Client) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case http.MethodPost:
		return handlers.CreateUser(req, dbClient)
	default:
		return nil, errors.New(ErrUnexpected)
	}
}

func ProjectsRouter(req events.APIGatewayProxyRequest, dbClient *dynamodb.Client) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case http.MethodPost:
		return handlers.CreateProject(req, dbClient)
	case http.MethodGet:
		return handlers.FindAllProjects(req)
	case http.MethodPut:
		return handlers.UpdateProject(req)
	case http.MethodDelete:
		return handlers.DeleteProject(req)
	default:
		return nil, errors.New(ErrUnexpected)
	}
}
