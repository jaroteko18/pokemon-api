package repository

import (
	"encoding/json"
	"fmt"
)

type PokemonSearch struct {
	ID          int    `json:"id,omitempty"`
	PokemonName string `json:"pokemon_name"`
	PokemonID   *int   `json:"pokemon_id,omitempty"`
	Found       bool   `json:"found"`
	SearchedAt  string `json:"searched_at,omitempty"`
}

type SearchStats struct {
	TotalSearches   int                `json:"total_searches"`
	FoundSearches   int                `json:"found_searches"`
	NotFoundSearches int               `json:"not_found_searches"`
	TopSearched     []TopSearchedItem  `json:"top_searched"`
	RecentSearches  []PokemonSearch    `json:"recent_searches"`
}

type TopSearchedItem struct {
	PokemonName string `json:"pokemon_name"`
	Count       int    `json:"count"`
}

type SearchRepository interface {
	LogSearch(pokemonName string, pokemonID *int, found bool) error
	GetStats() (*SearchStats, error)
}

type searchRepository struct {
	client *SupabaseClient
}

func NewSearchRepository(supabaseURL, supabaseKey string) SearchRepository {
	return &searchRepository{
		client: NewSupabaseClient(supabaseURL, supabaseKey),
	}
}

func (r *searchRepository) LogSearch(pokemonName string, pokemonID *int, found bool) error {
	search := PokemonSearch{
		PokemonName: pokemonName,
		PokemonID:   pokemonID,
		Found:       found,
	}

	_, err := r.client.Insert("pokemon_searches", search)
	return err
}

func (r *searchRepository) GetStats() (*SearchStats, error) {
	// Get total counts (ordered by searched_at)
	allBody, err := r.client.SelectAllOrdered("pokemon_searches", "searched_at.desc")
	if err != nil {
		return nil, fmt.Errorf("failed to get searches: %w", err)
	}

	var allSearches []PokemonSearch
	if err := json.Unmarshal(allBody, &allSearches); err != nil {
		return nil, fmt.Errorf("failed to parse searches: %w", err)
	}

	// Calculate stats
	stats := &SearchStats{
		TotalSearches: len(allSearches),
	}

	// Count found/not found
	countMap := make(map[string]int)
	for _, s := range allSearches {
		if s.Found {
			stats.FoundSearches++
		} else {
			stats.NotFoundSearches++
		}
		countMap[s.PokemonName]++
	}

	// Get top searched (sort by count)
	type kv struct {
		Key   string
		Value int
	}
	var sorted []kv
	for k, v := range countMap {
		sorted = append(sorted, kv{k, v})
	}
	// Simple bubble sort for top 10
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].Value > sorted[i].Value {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// Take top 10
	limit := 10
	if len(sorted) < limit {
		limit = len(sorted)
	}
	for i := 0; i < limit; i++ {
		stats.TopSearched = append(stats.TopSearched, TopSearchedItem{
			PokemonName: sorted[i].Key,
			Count:       sorted[i].Value,
		})
	}

	// Recent searches (first 10 from sorted by date desc)
	limit = 10
	if len(allSearches) < limit {
		limit = len(allSearches)
	}
	stats.RecentSearches = allSearches[:limit]

	return stats, nil
}
