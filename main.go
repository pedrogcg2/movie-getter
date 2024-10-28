package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	page := 1
	language := "pt-BR"
	bearer := ""

	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/popular?page=%d&language=%s",
		page, language)

	client := &http.Client{}
	req, error := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", bearer)

	response, error := client.Do(req)
	if error != nil {
		fmt.Println("Deu ruim pae")
		fmt.Println(error.Error())
		return
	}
	defer response.Body.Close()

	var movieResponse MovieResponse
	error = json.NewDecoder(response.Body).Decode(&movieResponse)

	if error != nil {
		fmt.Println("Deu ruim pae 3")
		fmt.Println(error.Error())
		return
	}

	bodyRequest := []MovieRequest{}

	for _, movie := range movieResponse.Movies {
		newMovie := MovieRequest{
			ExternalId:  movie.Id,
			Name:        movie.Title,
			Image:       movie.Image,
			Description: movie.Description,
		}
		bodyRequest = append(bodyRequest, newMovie)
	}

	bodyJson, _ := json.Marshal(bodyRequest)
	body := bytes.NewBuffer(bodyJson)

	fmt.Println(string(bodyJson))
	postReq, error := http.NewRequest("POST", "http://localhost:8000/movies", body)

	res, error := client.Do(postReq)

	if error != nil {
		fmt.Println("Deu ruim pae 4")
		fmt.Println(error.Error())
		return
	}
	defer res.Body.Close()

	fmt.Println(res.Status)
	bs, _ := io.ReadAll(res.Body)
	fmt.Println(string(bs))
	fmt.Println("deu bom")
}

type MovieResponse struct {
	Movies       []Movie `json:"results"`
	Page         int
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type Movie struct {
	Id          int
	Title       string
	Image       string `json:"poster_path"`
	Description string `json:"overview"`
}

type MovieRequest struct {
	ExternalId  int
	Name        string
	Image       string
	Description string
}
