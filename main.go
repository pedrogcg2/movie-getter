package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var BEARER_TMBD string

func main() {
	initEnv()

	page := 1

	client := &http.Client{}
	pageMax := 1

	for {
		fmt.Println(fmt.Sprintf("Current Page: %d", page))
		movieResponse, error := getMovies(client, page)

		if error != nil {
			log.Fatal("error: %s", error.Error())
		}

		error = sendMoviesToApi(client, &movieResponse.Movies)

		if error != nil {
			log.Fatal("error: %s", error.Error())
		}
		page = page + 1
		pageMax = movieResponse.TotalPages
		if page > pageMax {
			fmt.Println("Successfully get all pages")
			return
		}
	}
}

func initEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(fmt.Sprintf("Couldnt load env file: %s ", err.Error()))
	}

	BEARER_TMBD = os.Getenv("TMDB_AUTH")
}

func getMovies(client *http.Client, page int) (*MovieResponse, error) {
	req, error := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", BEARER_TMBD)

	fmt.Println(req.URL.String())

	response, error := client.Do(req)
	if error != nil {
		return nil, error
	}
	defer response.Body.Close()

	movieResponse := &MovieResponse{}
	error = json.NewDecoder(response.Body).Decode(movieResponse)

	if error != nil {
		return nil, error
	}

	return movieResponse, nil
}

func sendMoviesToApi(client *http.Client, movies *[]Movie) error {
	bodyRequest := []MovieRequest{}

	for _, movie := range *movies {
		newMovie := MovieRequest{
			ExternalId:  movie.Id,
			Name:        movie.Title,
			Image:       movie.Image,
			Description: movie.Description,
		}
		bodyRequest = append(bodyRequest, newMovie)
	}

	json, _ := json.Marshal(bodyRequest)
	bodyBuffer := bytes.NewBuffer(json)

	postReq, error := http.NewRequest("POST", "http://localhost:8000/movies", bodyBuffer)

	_, error = client.Do(postReq)

	if error != nil {
		return error
	}

	return nil
}

type MovieResponse struct {
	Movies       []Movie `json:"results"`
	Page         int
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type Movie struct {
	Title       string
	Image       string `json:"poster_path"`
	Description string `json:"overview"`
	Id          int
}

type MovieRequest struct {
	Name        string
	Image       string
	Description string
	ExternalId  int
}
