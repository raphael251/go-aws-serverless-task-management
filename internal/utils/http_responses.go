package utils

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

var ErrUnexpected = "unexpected error"

func HttpResponseCreated() (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{StatusCode: 201}, nil
}

func HttpResponseMethodNotAllowed() (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{StatusCode: 405}, nil
}

func HttpResponseOK(body interface{}) (*events.APIGatewayProxyResponse, error) {
	type OKResponse struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}

	response, err := json.Marshal(&OKResponse{Message: "success", Data: body})

	if err != nil {
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(response),
	}, nil
}

func HttpResponseNotFound() (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{StatusCode: 404}, nil
}

func HttpResponseBadRequest(message string, data interface{}) (*events.APIGatewayProxyResponse, error) {
	type BadRequestResponse struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}

	responseMessage := "bad request"
	if message != "" {
		responseMessage = message
	}

	response, err := json.Marshal(&BadRequestResponse{Message: responseMessage, Data: data})

	if err != nil {
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 400,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(response),
	}, nil
}

func HttpResponseInternalServerError(message string) (*events.APIGatewayProxyResponse, error) {
	type InternalServerErrorResponse struct {
		Message string `json:"message"`
	}

	responseMessage := "internal server error"
	if message != "" {
		responseMessage = message
	}

	response, err := json.Marshal(&InternalServerErrorResponse{Message: responseMessage})

	if err != nil {
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 500,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(response),
	}, nil
}
