package handler_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/raphael251/go-aws-serverless-task-management/internal/user/handler"
)

type mockDBClient struct {
	err error
}

func (c mockDBClient) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return nil, c.err
}

func (c mockDBClient) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return nil, c.err
}

func TestRegisterUser(t *testing.T) {
	testTable := []struct {
		name             string
		req              events.APIGatewayProxyRequest
		dbClient         func(t *testing.T) handler.DynamoDBClient
		expectedResponse *events.APIGatewayProxyResponse
		expectedErr      error
	}{
		{
			name: "should return bad request if an invalid json is sent",
			req: events.APIGatewayProxyRequest{
				Body: "invalid json string body",
			},
			dbClient: func(t *testing.T) handler.DynamoDBClient { return mockDBClient{} },
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
			dbClient: func(t *testing.T) handler.DynamoDBClient { return mockDBClient{} },
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
			dbClient: func(t *testing.T) handler.DynamoDBClient {
				return mockDBClient{
					err: errors.New("{\"message\":\"internal server error\",\"data\":[\"invalid field: Name\"]}"),
				}
			},
			expectedResponse: &events.APIGatewayProxyResponse{
				StatusCode: 500,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       "{\"message\":\"internal server error\"}",
			},
			expectedErr: nil,
		},
		// should return a bad request if get item return some register from database
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := handler.RegisterUser(tc.req, tc.dbClient(t))
			if !reflect.DeepEqual(resp, tc.expectedResponse) {
				t.Errorf("expected (%v), got (%v)", tc.expectedResponse, resp)
			}
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected (%v), got (%v)", tc.expectedErr, err)
			}
		})
	}
}
