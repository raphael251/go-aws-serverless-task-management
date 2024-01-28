package main

import (
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	handlers "github.com/raphael251/go-aws-serverless-task-management/internal"
)

func main() {
	lambda.Start(Handler)
}

var ErrUnexpected = "unexpected error"

func Handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "POST":
		return handlers.CreateProject(req)
	case "GET":
		return handlers.FindAllProjects(req)
	case "PUT":
		return handlers.UpdateProject(req)
	case "DELETE":
		return handlers.DeleteProject(req)
	}
	return nil, errors.New(ErrUnexpected)
}
