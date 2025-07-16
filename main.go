package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Joshua-SV/pokedexCLI/internal/pokeCache"
	"github.com/Joshua-SV/pokedexCLI/internal/utils"
)

func main() {
	// scanner for getting user input
	scanner := bufio.NewScanner(os.Stdin)
	var confi Config

	for true {
		// prompt user
		fmt.Print("Pokedex > ")
		// capture user input
		scanner.Scan()
		// get the input into string form
		txt := scanner.Text()
		// format user input
		words := cleanInput(txt)
		if len(words) != 0 {
			command, ok := registryCommands[words[0]]
			// if command is valid execute it
			if ok == true {
				err := command.callback(&confi, words[1:]...)
				if err != nil {
					fmt.Printf("Error Command: %v\n", err)
				}
			} else {
				fmt.Printf("Unknown command: %s\n", words[0])
			}
		}
	}

}

// use init function which is a special function called automatically by Go runtime
func init() {
	registryCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "displays the names of next 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "displays the names of previous 20 location areas in the Pokemon world",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "displays all pokemons found in the location specified",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catching Pokemon adds them to the user's Pokedex",
			callback:    commandCatch,
		},
	}
	// initialize an empty pokedex
	pokedex = make(map[string]utils.PokemonFull)
}

// global create cache to use during program
var cache = pokeCache.NewCache(12 * time.Second)

// global hash table for access of commands
var registryCommands map[string]cliCommand

// global hashtable for pokedex of captured Pokemons by the user
var pokedex map[string]utils.PokemonFull

// struct for managing commands
type cliCommand struct {
	name        string
	description string
	callback    func(cfg *Config, args ...string) error
}

// struct for pagination
type Config struct {
	Next     *string
	Previous *string
}

func cleanInput(txt string) []string {
	lowTxt := strings.ToLower(txt)
	lst := strings.Fields(lowTxt)

	return lst
}

// used as callback functions
func commandExit(cfg *Config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage Commands")
	fmt.Println("------------------------")
	for _, command := range registryCommands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(cfg *Config, args ...string) error {
	url := "https://pokeapi.co/api/v2/location-area/"

	// check if config.next exists for pagination url
	if cfg.Next != nil {
		url = *cfg.Next
	}

	// parse the json into the locations struct
	var locations utils.LocationResponse

	// use pokeAPI
	err := utils.GetMapPokeAPI(url, cache, &locations)
	if err != nil {
		return err
	}

	// update config for pagination
	cfg.Next = locations.Next
	cfg.Previous = locations.Previous
	// print all 20 locations
	for _, area := range locations.Results {
		fmt.Printf("%s\n", area.Name)
	}
	return nil
}

func commandMapBack(cfg *Config, args ...string) error {
	url := "https://pokeapi.co/api/v2/location-area/"

	// check if config.Previous exists for pagination url
	if cfg.Previous != nil {
		url = *cfg.Previous
	}

	// parse the json into the locations struct
	var locations utils.LocationResponse

	// use pokeAPI
	err := utils.GetMapPokeAPI(url, cache, &locations)
	if err != nil {
		return err
	}

	// update config for pagination
	cfg.Next = locations.Next
	cfg.Previous = locations.Previous
	// print all 20 locations
	for _, area := range locations.Results {
		fmt.Printf("%s\n", area.Name)
	}

	return nil
}

func commandExplore(cfg *Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("did not pass Location to explore: %v\n", args)
	}

	fmt.Printf("Exploring ...%s\n", args[0])

	// get location url
	url := "https://pokeapi.co/api/v2/location-area/" + args[0]

	var searched utils.LocationSearched

	err := utils.GetPokemonsOfLocation(url, cache, &searched)
	if err != nil {
		return fmt.Errorf("invalid location name: %v\n", err)
	}

	fmt.Println("Pokemons Found")
	fmt.Println("------------------------")
	for _, val := range searched.Pokemons_found {
		fmt.Printf("- %s\n", val.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("did not pass pokemon to catch: %v\n", args)
	}

	url := "https://pokeapi.co/api/v2/pokemon/" + args[0]

	fmt.Printf("Throwing a Pokeball at %s...\n", args[0])

	var pokemon utils.PokemonFull

	err := utils.GetPokemon(url, cache, &pokemon)
	if err != nil {
		return err
	}

	// Create and seed a new Rand instance (recommended for Go 1.20+)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// calculate the catch using modulus
	caught := r.Intn(pokemon.BaseExperience) % pokemon.BaseExperience
	// check catch
	if caught >= 0 && caught < (pokemon.BaseExperience/4)+3 {
		pokedex[pokemon.Name] = pokemon
		fmt.Printf("%s was caught!\n", pokemon.Name)
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func commandInspect(cfg *Config, args ...string) error {

	return nil
}
