package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/whiterthanwhite/aws_movies_practice/internal/movies"
)

func readMovies(fileName string) ([]movies.Movie, error) {
	movies := make([]movies.Movie, 0)

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return movies, err
	}

	err = json.Unmarshal(data, &movies)
	if err != nil {
		return movies, err
	}

	return movies, err
}

func insertMovie(cfg aws.Config, movie movies.Movie) error {
	type m struct {
		ID   string
		Name string
	}
	item, err := attributevalue.MarshalMap(&m{
		ID:   fmt.Sprint(movie.ID),
		Name: movie.Name,
	})
	if err != nil {
		return err
	}

	svc := dynamodb.NewFromConfig(cfg)

	_, err = svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("movies"),
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-north-1"))
	if err != nil {
		log.Fatal()
	}

	movies, err := readMovies("movies.json")
	if err != nil {
		log.Fatal(err)
	}

	for _, movie := range movies {
		fmt.Println("Inserting:", movie.Name)
		err = insertMovie(cfg, movie)
		if err != nil {
			log.Fatal(err)
		}
	}
}
