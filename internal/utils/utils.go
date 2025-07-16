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

// structs for getting the pokemons found with explore command
type LocationSearched struct {
	Name           string      `json:"name"`
	ID             int         `json:"id"`
	Pokemons_found []Encounter `json:"pokemon_encounters"`
}

type Encounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// structs for catching and details of a specific pokemon
type PokemonFull struct {
	Name           string `json:"name"`
	ID             int    `json:"id"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"` // height of this Pokémon in decimetres
	Weight         int    `json:"weight"` // weight of this Pokémon in hectograms
	Order          int    `json:"order"`  // used to sort the pokemon

}

func GetMapPokeAPI(url string, cache *pokeCache.Cache, locations *LocationResponse) error {
	body, err := GetBody(url, cache)
	if err != nil {
		return err
	}

	// parse the json into the locations struct
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return err
	}

	return nil
}

func GetPokemonsOfLocation(url string, cache *pokeCache.Cache, search *LocationSearched) error {
	body, err := GetBody(url, cache)
	if err != nil {
		return err
	}

	// parse the json into the locations struct
	err = json.Unmarshal(body, &search)
	if err != nil {
		return err
	}

	return nil
}

func GetPokemon(url string, cache *pokeCache.Cache, pokemon *PokemonFull) error {
	body, err := GetBody(url, cache)
	if err != nil {
		return err
	}

	// parse the json into the locations struct
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return err
	}

	return nil
}

// get the cached body or https the body
func GetBody(url string, cache *pokeCache.Cache) ([]byte, error) {
	body, okay := cache.Get(url)

	// if cache does not have the data fetch from https
	if !okay {
		// get https request and response
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}

		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
	}

	// add the []byte data into the cache
	cache.Add(url, body)

	return body, nil
}
