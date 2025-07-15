package main

import (
	"bufio"
	"fmt"
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
				fmt.Println("Unknown command")
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
	}
}

// create cache to use during program
var cache = pokeCache.NewCache(12 * time.Second)

// hash table for access of commands
var registryCommands map[string]cliCommand

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
		return fmt.Errorf("did not pass Location to explore: %v", args)
	}

	fmt.Printf("Exploring ...%s\n", args[0])

	// get location url
	url := "https://pokeapi.co/api/v2/location-area/" + args[0]

	var searched utils.LocationSearched

	err := utils.GetPokemonsOfLocation(url, cache, &searched)
	if err != nil {
		return err
	}

	fmt.Println("Pokemons Found")
	fmt.Println("------------------------")
	for _, val := range searched.Pokemons_found {
		fmt.Printf("- %s\n", val.Pokemon.Name)
	}

	return nil
}
