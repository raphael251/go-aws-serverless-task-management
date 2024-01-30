package routers

import (
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/raphael251/go-aws-serverless-task-management/internal/handlers"
	"github.com/raphael251/go-aws-serverless-task-management/internal/utils"
)

func UsersRouter(req events.APIGatewayProxyRequest, dbClient *dynamodb.Client) (*events.APIGatewayProxyResponse, error) {
	switch req.Path {
	case "/api/users/register":
		switch req.HTTPMethod {
		case http.MethodPost:
			return handlers.RegisterUser(req, dbClient)
		default:
			return nil, errors.New(utils.ErrUnexpected)
		}
	case "/api/users/login":
		switch req.HTTPMethod {
		case http.MethodPost:
			return handlers.LogUserIn(req, dbClient)
		default:
			return nil, errors.New(utils.ErrUnexpected)
		}
	default:
		return nil, errors.New(utils.ErrUnexpected)
	}

}
