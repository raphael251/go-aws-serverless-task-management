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
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterUser(req events.APIGatewayProxyRequest, dbClient *dynamodb.Client) (*events.APIGatewayProxyResponse, error) {
	var user *CreateUserInput
	err := json.Unmarshal([]byte(req.Body), &user)
	if err != nil {
		return utils.HttpResponseBadRequest("")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("error hashing the user password for registering user", err)
		return utils.HttpResponseInternalServerError("")
	}

	_, err = dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"pk":       &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", database.UserPKPrepend, user.Username)},
			"sk":       &types.AttributeValueMemberS{Value: "info"},
			"username": &types.AttributeValueMemberS{Value: user.Username},
			"name":     &types.AttributeValueMemberS{Value: user.Name},
			"email":    &types.AttributeValueMemberS{Value: user.Email},
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
