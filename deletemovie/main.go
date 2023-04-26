package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/whiterthanwhite/aws_movies_practice/internal/movies"
)

func deleteMovie(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var movie movies.Movie
	err := json.Unmarshal([]byte(request.Body), &movie)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid payload",
		}, nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while retrieving credentials",
		}, nil
	}

	svc := dynamodb.NewFromConfig(cfg)
	_, err = svc.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{
				Value: *aws.String(fmt.Sprint(movie.ID)),
			},
		},
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while deleting movie from DynamoDB",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       "{}",
	}, nil
}

func main() {
	lambda.Start(deleteMovie)
}
