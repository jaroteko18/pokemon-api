package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type PokemonService interface {
	GetPokemon(nameOrID string) (*PokemonResponse, error)
}

type pokemonService struct {
	client  *http.Client
	baseURL string
}

type PokemonResponse struct {
	Found bool         `json:"found"`
	Data  *PokemonData `json:"data,omitempty"`
}

type PokemonData struct {
	ID        int          `json:"id"`
	Name      string       `json:"name"`
	Types     string       `json:"types"`
	Abilities string       `json:"abilities"`
	Stats     PokemonStats `json:"stats"`
	Height    string       `json:"height"`
	Weight    string       `json:"weight"`
	Sprite    string       `json:"sprite"`
}

type PokemonStats struct {
	HP        int `json:"hp"`
	Attack    int `json:"attack"`
	Defense   int `json:"defense"`
	SpAttack  int `json:"spAttack"`
	SpDefense int `json:"spDefense"`
	Speed     int `json:"speed"`
}

func NewPokemonService() PokemonService {
	return &pokemonService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://pokeapi.co/api/v2",
	}
}

func (s *pokemonService) GetPokemon(nameOrID string) (*PokemonResponse, error) {
	url := fmt.Sprintf("%s/pokemon/%s", s.baseURL, strings.ToLower(nameOrID))

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return &PokemonResponse{Found: false}, nil
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("PokeAPI returned status %d", resp.StatusCode)
	}

	var rawData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawData); err != nil {
		return nil, err
	}

	return &PokemonResponse{
		Found: true,
		Data:  s.transformData(rawData),
	}, nil
}

func (s *pokemonService) transformData(raw map[string]interface{}) *PokemonData {
	data := &PokemonData{
		ID:   int(raw["id"].(float64)),
		Name: capitalize(raw["name"].(string)),
	}

	// Extract types
	types := raw["types"].([]interface{})
	typeNames := make([]string, len(types))
	for i, t := range types {
		typeMap := t.(map[string]interface{})
		typeName := typeMap["type"].(map[string]interface{})["name"].(string)
		typeNames[i] = capitalize(typeName)
	}
	data.Types = strings.Join(typeNames, ", ")

	// Extract abilities
	abilities := raw["abilities"].([]interface{})
	abilityNames := make([]string, len(abilities))
	for i, a := range abilities {
		abilityMap := a.(map[string]interface{})
		abilityName := abilityMap["ability"].(map[string]interface{})["name"].(string)
		abilityNames[i] = capitalize(strings.ReplaceAll(abilityName, "-", " "))
	}
	data.Abilities = strings.Join(abilityNames, ", ")

	// Extract stats
	stats := raw["stats"].([]interface{})
	data.Stats = PokemonStats{
		HP:        int(stats[0].(map[string]interface{})["base_stat"].(float64)),
		Attack:    int(stats[1].(map[string]interface{})["base_stat"].(float64)),
		Defense:   int(stats[2].(map[string]interface{})["base_stat"].(float64)),
		SpAttack:  int(stats[3].(map[string]interface{})["base_stat"].(float64)),
		SpDefense: int(stats[4].(map[string]interface{})["base_stat"].(float64)),
		Speed:     int(stats[5].(map[string]interface{})["base_stat"].(float64)),
	}

	// Height and weight
	height := raw["height"].(float64)
	weight := raw["weight"].(float64)
	data.Height = fmt.Sprintf("%.1f", height/10)
	data.Weight = fmt.Sprintf("%.1f", weight/10)

	// Sprite
	sprites := raw["sprites"].(map[string]interface{})
	if other, ok := sprites["other"].(map[string]interface{}); ok {
		if artwork, ok := other["official-artwork"].(map[string]interface{}); ok {
			if front, ok := artwork["front_default"].(string); ok {
				data.Sprite = front
			}
		}
	}
	if data.Sprite == "" {
		if front, ok := sprites["front_default"].(string); ok {
			data.Sprite = front
		}
	}

	return data
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
