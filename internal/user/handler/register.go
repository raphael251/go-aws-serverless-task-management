package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/raphael251/go-aws-serverless-task-management/internal/database"
	"github.com/raphael251/go-aws-serverless-task-management/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserInput struct {
	Username string `json:"username" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type DynamoDBClient interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

func RegisterUser(req events.APIGatewayProxyRequest, dbClient DynamoDBClient) (*events.APIGatewayProxyResponse, error) {
	var input *CreateUserInput
	err := json.Unmarshal([]byte(req.Body), &input)
	if err != nil {
		return utils.HttpResponseBadRequest("", nil)
	}

	if validationErrors := utils.ValidateRequestInput(input); validationErrors != nil {
		errs := make([]string, 0)
		for _, err := range validationErrors {
			errs = append(errs, err.Error())
		}

		return utils.HttpResponseBadRequest("", errs)
	}

	dbItem, err := dbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s%s", database.UserPKPrepend, input.Username),
			},
			"sk": &types.AttributeValueMemberS{
				Value: "info",
			},
		},
		TableName: &database.AppTableName,
	})
	if err != nil {
		fmt.Println("error performing a GetItem to register user", err)
		return utils.HttpResponseInternalServerError("")
	}
	if dbItem.Item != nil {
		return utils.HttpResponseBadRequest("username already in use. Please choose another one.", nil)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("error hashing the user password for registering user", err)
		return utils.HttpResponseInternalServerError("")
	}

	_, err = dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"pk":       &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", database.UserPKPrepend, input.Username)},
			"sk":       &types.AttributeValueMemberS{Value: "info"},
			"username": &types.AttributeValueMemberS{Value: input.Username},
			"name":     &types.AttributeValueMemberS{Value: input.Name},
			"email":    &types.AttributeValueMemberS{Value: input.Email},
			"password": &types.AttributeValueMemberS{Value: string(hashedPassword)},
		},
		TableName: &database.AppTableName,
	})
	if err != nil {
		fmt.Println("error performing a PutItem for registering user", err)
		return utils.HttpResponseInternalServerError("")
	}

	return utils.HttpResponseCreated()
}
