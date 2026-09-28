package main

import (
	"fmt"
	"math/rand"
	"os"

	pokeAPI "github.com/sableHoujuu/cliPokedex/internal"
)

// All of the functions and types that are directly invoked with commands

type config struct {
	commandRegistry map[string]cliCommand
	cmdMap          struct {
		nextURL     *string
		previousURL *string
	}
	pokeAPIClient pokeAPI.Client
	pokedex       map[string]pokeAPI.PokemonData
}

type parameter struct {
	signature string
	usage     string
}

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, params []string) error
	parameters  []parameter
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Goes to the next page of the location list, or the first page if the command has not been used in this session",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Goes to the previous page of the location list",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Displays a list of pokemon encounterable in the specified area",
			callback:    commandExplore,
			parameters:  []parameter{{signature: "<area_name>", usage: "explore <area_name>"}},
		},
		"catch": {
			name:        "catch",
			description: "Attempts to catch a given Pokemon",
			callback:    commandCatch,
			parameters:  []parameter{{signature: "<pokemon>", usage: "catch <pokemon>"}},
		},
		"inspect": {
			name:        "inspect",
			description: "Lists information about a pokemon, if previously caught",
			callback:    commandInspect,
			parameters:  []parameter{{signature: "<pokemon>", usage: "inspect <pokemon>"}},
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists all pokemon the user has caught in this session",
			callback:    commandPokedex,
			parameters:  []parameter{{}},
		},
	}
}

func commandExit(cfg *config, params []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, params []string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range cfg.commandRegistry {
		fmt.Printf("%v: %v\n", command.name, command.description)
	}
	return nil
}

func commandMap(cfg *config, params []string) error {
	url := ""
	if cfg.cmdMap.nextURL != nil {
		url = *cfg.cmdMap.nextURL
	}

	locations, nextUrl, previousUrl, err := cfg.pokeAPIClient.GetLocations(url)
	if err != nil {
		return err
	}

	for i := 0; i < len(locations); i++ {
		fmt.Println(locations[i])
	}
	cfg.cmdMap.nextURL = nextUrl
	cfg.cmdMap.previousURL = previousUrl
	return nil
}

func commandMapb(cfg *config, params []string) error {
	var url string
	if cfg.cmdMap.previousURL != nil {
		url = *cfg.cmdMap.previousURL
	} else {
		fmt.Println("You're on the first page, use `map` to page forwards.")
		return nil
	}

	locations, nextUrl, previousUrl, err := cfg.pokeAPIClient.GetLocations(url)
	if err != nil {
		return err
	}

	for i := 0; i < len(locations); i++ {
		fmt.Println(locations[i])
	}

	cfg.cmdMap.nextURL = nextUrl
	cfg.cmdMap.previousURL = previousUrl
	return nil
}

func commandExplore(cfg *config, params []string) error {
	if len(params) <= 0 || params[0] == "" {
		fmt.Println("Please provide an area name.")
		return nil
	}
	areaName := params[0]

	fmt.Printf("Exploring %v...\n", areaName)

	pokemon, err := cfg.pokeAPIClient.ExploreLocation(areaName)
	if err != nil {
		return err
	}

	if len(pokemon) <= 0 {
		fmt.Println("No area found.")
		return nil
	}

	fmt.Println("Found Pokemon:")
	for i := 0; i < len(pokemon); i++ {
		fmt.Printf(" - %v\n", pokemon[i])
	}

	return nil
}

func commandCatch(cfg *config, params []string) error {
	if len(params) <= 0 || params[0] == "" {
		fmt.Println("Please provide a Pokemon.")
		return nil
	}
	pokemonName := params[0]

	data, err := cfg.pokeAPIClient.GetPokemonData(pokemonName)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %v...\n", pokemonName)

	caught := (func() bool {
		n := rand.Intn(data.BaseExperience + 100)
		if n <= 100 {
			return true
		}
		return false
	})()

	if caught {
		fmt.Printf("%v was caught!\n", pokemonName)
		cfg.pokedex[data.Name] = data
	} else {
		fmt.Printf("%v escaped!\n", pokemonName)
	}

	return nil
}

func commandInspect(cfg *config, params []string) error {
	data, exists := cfg.pokedex[params[0]]
	if !exists {
		fmt.Println("You have not caught that pokemon.")
		return nil
	}

	fmt.Printf("Name: %v\nHeight: %v\nWeight: %v\n",
		data.Name,
		data.Height,
		data.Weight,
	)

	fmt.Print("Stats:\n")
	for _, statData := range data.Stats {
		fmt.Printf(" -%v: %v\n", statData.Stat.Name, statData.BaseStat)
	}

	fmt.Print("Types:\n")
	for _, typeData := range data.Types {
		fmt.Printf(" -%v\n", typeData.Type.Name)
	}

	return nil
}

func commandPokedex(cfg *config, params []string) error {
	if len(cfg.pokedex) == 0 {
		fmt.Println("Your Pokedex is empty. Use `catch` to begin your journey!")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name, _ := range cfg.pokedex {
		fmt.Printf(" - %v\n", name)
	}
	return nil
}
