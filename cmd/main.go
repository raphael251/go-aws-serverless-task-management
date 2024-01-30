package main

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/raphael251/go-aws-serverless-task-management/internal/handlers"
	"github.com/raphael251/go-aws-serverless-task-management/internal/routers"
	"github.com/raphael251/go-aws-serverless-task-management/internal/utils"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {

	lambda.Start(Handler)
}

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
		return routers.ProjectsRouter(req, dbClient)
	default:
		return nil, errors.New(utils.ErrUnexpected)
	}
}

func UsersRouter(req events.APIGatewayProxyRequest, dbClient *dynamodb.Client) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case http.MethodPost:
		return handlers.CreateUser(req, dbClient)
	default:
		return nil, errors.New(utils.ErrUnexpected)
	}
}
