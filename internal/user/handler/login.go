package handler

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
	"github.com/raphael251/go-aws-serverless-task-management/internal/user/entity"
	"github.com/raphael251/go-aws-serverless-task-management/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

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
				Value: fmt.Sprintf("%s%s", database.UserPKPrepend, input.Username),
			},
			"sk": &types.AttributeValueMemberS{
				Value: "info",
			},
		},
		TableName: &database.AppTableName,
	})
	if err != nil {
		fmt.Println("error performing a GetItem to log user in", err)
		return utils.HttpResponseInternalServerError("")
	}

	foundUser := entity.User{}
	err = attributevalue.UnmarshalMap(dbItem.Item, &foundUser)
	if err != nil {
		fmt.Println("error trying to unmarshal the user from the db in login", err)
		return utils.HttpResponseInternalServerError("")
	}

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
