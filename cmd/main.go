package main

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/raphael251/go-aws-serverless-task-management/internal/routers"
	userRouter "github.com/raphael251/go-aws-serverless-task-management/internal/user/router"
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

	switch {
	case strings.Contains(req.Path, "/api/users"):
		return userRouter.Route(req, dbClient)
	case strings.Contains(req.Path, "/api/projects"):
		return routers.ProjectsRouter(req, dbClient)
	default:
		return utils.HttpResponseNotFound()
	}
}
