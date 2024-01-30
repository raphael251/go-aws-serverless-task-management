package utils

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

var ErrUnexpected = "unexpected error"

func HttpResponseCreated() (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func HttpResponseOK(body interface{}) (*events.APIGatewayProxyResponse, error) {
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
