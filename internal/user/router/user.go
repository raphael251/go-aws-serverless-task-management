package router

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/raphael251/go-aws-serverless-task-management/internal/user/handler"
	"github.com/raphael251/go-aws-serverless-task-management/internal/utils"
)

func Route(req events.APIGatewayProxyRequest, dbClient *dynamodb.Client) (*events.APIGatewayProxyResponse, error) {
	switch req.Path {
	case "/api/users/register":
		switch req.HTTPMethod {
		case http.MethodPost:
			return handler.RegisterUser(req, dbClient)
		default:
			return utils.HttpResponseMethodNotAllowed()
		}
	case "/api/users/login":
		switch req.HTTPMethod {
		case http.MethodPost:
			return handler.LogUserIn(req, dbClient)
		default:
			return utils.HttpResponseMethodNotAllowed()
		}
	default:
		return utils.HttpResponseNotFound()
	}

}
