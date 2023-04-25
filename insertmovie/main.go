package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"example.com/internal/movies"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

func insert(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var movie movies.Movie
	err := json.Unmarshal([]byte(request.Body), &movie)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	svc := dynamodb.NewFromConfig(cfg)
	_, err = svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Item: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{
				Value: *aws.String(fmt.Sprint(movie.ID)),
			},
			"Name": &types.AttributeValueMemberS{
				Value: *aws.String(movie.Name),
			},
		},
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       "{}",
	}, nil
}

func main() {
	lambda.Start(insert)
}
