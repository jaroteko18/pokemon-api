package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/pokemon-chatbot-api/internal/services"
)

type PokemonHandler struct {
	service services.PokemonService
}

func NewPokemonHandler(service services.PokemonService) *PokemonHandler {
	return &PokemonHandler{service: service}
}

func (h *PokemonHandler) GetPokemon(c *gin.Context) {
	name := c.Param("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"found": false,
			"error": "Pokemon name is required",
		})
		return
	}

	result, err := h.service.GetPokemon(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"found": false,
			"error": "Failed to fetch Pokemon data",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *PokemonHandler) SearchPokemon(c *gin.Context) {
	query := c.Param("query")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"found": false,
			"error": "Search query is required",
		})
		return
	}

	// For now, search is same as get
	// You can enhance this later with fuzzy search
	result, err := h.service.GetPokemon(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"found": false,
			"error": "Failed to search Pokemon",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *PokemonHandler) GetSearchStats(c *gin.Context) {
	stats, err := h.service.GetSearchStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"stats":   stats,
	})
}
