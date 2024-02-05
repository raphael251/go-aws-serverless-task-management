package handler_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/raphael251/go-aws-serverless-task-management/internal/user/handler"
)

type mockDBClient struct {
	getItemEmpty bool
	getItemErr   error
	putItemErr   error
}

func (c mockDBClient) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	if c.getItemErr != nil {
		return nil, c.getItemErr
	}

	if c.getItemEmpty {
		return &dynamodb.GetItemOutput{
			Item: nil,
		}, nil
	}

	return &dynamodb.GetItemOutput{
		Item: map[string]types.AttributeValue{
			"name": &types.AttributeValueMemberS{Value: "Any Name"},
		},
	}, nil
}

func (c mockDBClient) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	if c.putItemErr != nil {
		return nil, c.putItemErr
	}

	return &dynamodb.PutItemOutput{}, nil
}

type mockHashGenerator struct {
	err error
}

func (g mockHashGenerator) GenerateFromPassword(password string) (string, error) {
	if g.err != nil {
		return "", g.err
	}
	return "any-hash", nil
}

func TestRegisterUser(t *testing.T) {
	testTable := []struct {
		name             string
		req              events.APIGatewayProxyRequest
		dbClient         handler.DynamoDBClient
		hashGenerator    handler.HashGenerator
		expectedResponse *events.APIGatewayProxyResponse
		expectedErr      error
	}{
		{
			name: "should return bad request if an invalid json is sent",
			req: events.APIGatewayProxyRequest{
				Body: "invalid json string body",
			},
			dbClient:      mockDBClient{},
			hashGenerator: mockHashGenerator{},
			expectedResponse: &events.APIGatewayProxyResponse{
				StatusCode: 400,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       "{\"message\":\"bad request\"}",
			},
			expectedErr: nil,
		},
		{
			name: "should return invalid fields if an invalid json is sent",
			req: events.APIGatewayProxyRequest{
				Body: "{\"username\": \"raphael251\",\"email\": \"raphael251@hotmail.com\", \"password\": \"123456\"}",
			},
			dbClient:      mockDBClient{},
			hashGenerator: mockHashGenerator{},
			expectedResponse: &events.APIGatewayProxyResponse{
				StatusCode: 400,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       "{\"message\":\"bad request\",\"data\":[\"invalid field: Name\"]}",
			},
			expectedErr: nil,
		},
		{
			name: "should return an internal server error if get item returns an error",
			req: events.APIGatewayProxyRequest{
				Body: "{\"username\": \"raphael251\",\"name\": \"Raphael Passos\",\"email\": \"raphael251@hotmail.com\", \"password\": \"123456\"}",
			},
			dbClient: mockDBClient{
				getItemErr: errors.New("unexpected error occurred"),
			},
			hashGenerator: mockHashGenerator{},
			expectedResponse: &events.APIGatewayProxyResponse{
				StatusCode: 500,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       "{\"message\":\"internal server error\"}",
			},
			expectedErr: nil,
		},
		{
			name: "should return a bad request if get item return some register from database, meaning that a user with this username already exists",
			req: events.APIGatewayProxyRequest{
				Body: "{\"username\": \"raphael251\",\"name\": \"Raphael Passos\",\"email\": \"raphael251@hotmail.com\", \"password\": \"123456\"}",
			},
			dbClient: mockDBClient{
				getItemErr:   nil,
				getItemEmpty: false,
			},
			hashGenerator: mockHashGenerator{},
			expectedResponse: &events.APIGatewayProxyResponse{
				StatusCode: 400,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       "{\"message\":\"username already in use. Please choose another one.\"}",
			},
			expectedErr: nil,
		},
		{
			name: "should return an internal server error if the hash generator returns some error",
			req: events.APIGatewayProxyRequest{
				Body: "{\"username\": \"raphael251\",\"name\": \"Raphael Passos\",\"email\": \"raphael251@hotmail.com\", \"password\": \"123456\"}",
			},
			dbClient: mockDBClient{
				getItemErr:   nil,
				getItemEmpty: true,
			},
			hashGenerator: mockHashGenerator{
				err: errors.New("unexpected error occurred"),
			},
			expectedResponse: &events.APIGatewayProxyResponse{
				StatusCode: 500,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       "{\"message\":\"internal server error\"}",
			},
			expectedErr: nil,
		},
		{
			name: "should return an internal server error if the put item db func returns some error",
			req: events.APIGatewayProxyRequest{
				Body: "{\"username\": \"raphael251\",\"name\": \"Raphael Passos\",\"email\": \"raphael251@hotmail.com\", \"password\": \"123456\"}",
			},
			dbClient: mockDBClient{
				getItemErr:   nil,
				getItemEmpty: true,
				putItemErr:   errors.New("unexpected error occurred"),
			},
			hashGenerator: mockHashGenerator{
				err: nil,
			},
			expectedResponse: &events.APIGatewayProxyResponse{
				StatusCode: 500,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       "{\"message\":\"internal server error\"}",
			},
			expectedErr: nil,
		},
		{
			name: "should return the created status code (201) if everything worked",
			req: events.APIGatewayProxyRequest{
				Body: "{\"username\": \"raphael251\",\"name\": \"Raphael Passos\",\"email\": \"raphael251@hotmail.com\", \"password\": \"123456\"}",
			},
			dbClient: mockDBClient{
				getItemErr:   nil,
				getItemEmpty: true,
				putItemErr:   nil,
			},
			hashGenerator: mockHashGenerator{
				err: nil,
			},
			expectedResponse: &events.APIGatewayProxyResponse{
				StatusCode: 201,
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := handler.RegisterUser(tc.req, tc.dbClient, tc.hashGenerator)
			if !reflect.DeepEqual(resp, tc.expectedResponse) {
				t.Errorf("expected (%v), got (%v)", tc.expectedResponse, resp)
			}
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected (%v), got (%v)", tc.expectedErr, err)
			}
		})
	}
}
