package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Joshua-SV/pokedexCLI/internal/pokeCache"
)

// struct for api decoding
type LocationResponse struct {
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Count    int     `json:"count"`
	Results  []Area  `json:"results"`
}

type Area struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func GetMapPokeAPI(url string, cache *pokeCache.Cache) (LocationResponse, error) {
	body, okay := cache.Get(url)

	// if cache does not have the data fetch from https
	if !okay {
		// get https request and response
		res, err := http.Get(url)
		if err != nil {
			return LocationResponse{}, err
		}

		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return LocationResponse{}, err
		}

		// add the []byte data into the cache
		cache.Add(url, body)
	}

	// parse the json into the locations struct
	var locations LocationResponse
	err := json.Unmarshal(body, &locations)
	if err != nil {
		return LocationResponse{}, err
	}

	return locations, nil
}
