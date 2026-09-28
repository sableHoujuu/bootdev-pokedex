package main

import (
	"bufio"
	"fmt"
	"os"

	pokeAPI "github.com/sableHoujuu/cliPokedex/internal"
)

func main() {
	cfg := &config{
		commandRegistry: getCommands(),
		pokeAPIClient:   pokeAPI.NewClient(30, 60),
		pokedex:         make(map[string]pokeAPI.PokemonData),
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		var input string
		if scanner.Scan() {
			input = scanner.Text()
		} else if scanner.Err() != nil {
			fmt.Printf("Couldn't grab input: %v\n", scanner.Err())
		} else {
			fmt.Println("Please provide an input.")
			continue
		}

		strs := cleanInput(input)
		if len(strs) < 1 {
			fmt.Println("Please input a command.")
			continue
		}
		command, ok := cfg.commandRegistry[strs[0]]
		if !ok || command.callback == nil {
			fmt.Println("Unknown command")
			continue
		}

		if err := command.callback(cfg, strs[1:]); err != nil {
			fmt.Printf("error encountered during callback: %v\n", err)
		}
	}
}
