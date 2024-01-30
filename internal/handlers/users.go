package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/raphael251/go-aws-serverless-task-management/internal/database"
	"github.com/raphael251/go-aws-serverless-task-management/internal/entity"
	"github.com/raphael251/go-aws-serverless-task-management/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

var userPKPrepend = "user#"

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
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	_, err = dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"pk":       &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", userPKPrepend, user.Username)},
			"sk":       &types.AttributeValueMemberS{Value: "info"},
			"username": &types.AttributeValueMemberS{Value: user.Username},
			"name":     &types.AttributeValueMemberS{Value: user.Name},
			"email":    &types.AttributeValueMemberS{Value: user.Email},
			"password": &types.AttributeValueMemberS{Value: string(hashedPassword)},
		},
		TableName: &database.AppTableName,
	})

	if err != nil {
		return nil, err
	}

	return utils.HttpResponseCreated()
}

type LogUserInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LogUserInOutput struct {
	AccessToken string `json:"access_token"`
}

var errInvalidUsernameOrPassword = "invalid username or password"

func LogUserIn(req events.APIGatewayProxyRequest, dbClient *dynamodb.Client) (*events.APIGatewayProxyResponse, error) {
	var input *LogUserInInput
	err := json.Unmarshal([]byte(req.Body), &input)
	if err != nil {
		return nil, err
	}

	dbItem, err := dbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s%s", userPKPrepend, input.Username),
			},
			"sk": &types.AttributeValueMemberS{
				Value: "info",
			},
		},
		TableName: &database.AppTableName,
	})
	if err != nil {
		return nil, err
	}

	foundUser := entity.User{}
	attributevalue.UnmarshalMap(dbItem.Item, &foundUser)

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(input.Password))
	if err != nil {
		return utils.HttpResponseBadRequest(errInvalidUsernameOrPassword)
	}

	secret := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": foundUser.Username,
	})
	stringToken, err := token.SignedString(secret)
	if err != nil {
		fmt.Println("error creating JWT token", err)
		return utils.HttpResponseInternalServerError("")
	}

	return utils.HttpResponseOK(LogUserInOutput{AccessToken: stringToken})
}
