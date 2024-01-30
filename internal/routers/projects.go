package routers

import (
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/raphael251/go-aws-serverless-task-management/internal/handlers"
	"github.com/raphael251/go-aws-serverless-task-management/internal/utils"
)

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
		return nil, errors.New(utils.ErrUnexpected)
	}
}
