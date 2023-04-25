package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
)

func findAll(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	size, err := strconv.Atoi(request.Headers["Count"])
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Count header should be a nubmer",
		}, nil
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-north-1"))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while retrieving AWS credentials",
		}, err
	}

	svc := dynamodb.NewFromConfig(cfg)
	scanInput, err := svc.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Limit:     aws.Int32(int32(size)),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while scanning DynamoDB",
		}, nil
	}

	resp, err := json.Marshal(scanInput.Items)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while decoding to string value",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(resp),
	}, nil
}

func main() {
	lambda.Start(findAll)
}
