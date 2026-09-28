package pokeAPI

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	pokecache "github.com/sableHoujuu/cliPokedex/internal/pokecache"
)

// A client used for communication with pokeAPI
type Client struct {
	httpClient http.Client
	cache      *pokecache.Cache
}

// Creates and returns a new client. Note that the arguments are converted to seconds with * time.Second.
func NewClient(timeout time.Duration, cacheInterval time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		cache: pokecache.NewCache(cacheInterval * time.Second),
	}

}

func (c Client) GetLocations(locationUrl string) (locations []string, nextUrl *string, prevUrl *string, e error) {
	if locationUrl == "" {
		locationUrl = "https://pokeapi.co/api/v2/location-area/"
	}

	data, exists := c.cache.Get(locationUrl)
	if !exists {
		resp, err := http.Get(locationUrl)
		if err != nil {
			return []string{}, nil, nil, err
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return []string{}, nil, nil, err
		}

		// We cache the data for later
		c.cache.Add(locationUrl, data)
	}
	var results LocationList
	if err := json.Unmarshal(data, &results); err != nil {
		return []string{}, nil, nil, err
	}

	if results.Results == nil {
		return []string{}, nil, nil, errors.New("No results found.")
	}
	var locationList []string
	for _, area := range results.Results {
		locationList = append(locationList, area.Name)
	}

	return locationList, results.Next, results.Previous, nil
}

func (c Client) ExploreLocation(locationName string) (pokemon []string, e error) {
	if locationName == "" {
		return []string{}, nil
	}

	fullUrl := "https://pokeapi.co/api/v2/location-area/" + locationName

	data, exists := c.cache.Get(fullUrl)
	if !exists {
		resp, err := http.Get(fullUrl)
		if err != nil {
			return []string{}, err
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return []string{}, err
		} else if len(data) <= 0 {
			return []string{}, nil
		}
		// We cache the data for later
		c.cache.Add(locationName, data)
	}

	var results LocationData
	if err := json.Unmarshal(data, &results); err != nil {
		return []string{}, err
	}

	if results.PokemonEncounters == nil {
		return []string{}, nil
	}

	var pokemonList []string
	for _, val := range results.PokemonEncounters {
		pokemonList = append(pokemonList, val.Pokemon.Name)
	}

	return pokemonList, nil
}

func (c Client) GetPokemonData(name string) (PokemonData, error) {
	if name == "" {
		return PokemonData{}, errors.New("Please provide a Pokemon name.")
	}

	fullUrl := "https://pokeapi.co/api/v2/pokemon/" + name

	data, exists := c.cache.Get(fullUrl)
	if !exists {
		resp, err := http.Get(fullUrl)
		if err != nil {
			return PokemonData{}, err
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return PokemonData{}, err
		} else if len(data) <= 0 {
			return PokemonData{}, nil
		}
		// We cache the data for later
		c.cache.Add(name, data)
	}

	var results PokemonData
	if err := json.Unmarshal(data, &results); err != nil {
		return PokemonData{}, err
	}

	if results.ID == 0 {
		// No pokemon has an ID of 0, this means it wasn't found
		return PokemonData{}, errors.New("No Pokemon found.")
	}

	return results, nil
}
