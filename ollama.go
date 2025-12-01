package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// OllamaModel represents a model from Ollama API
type OllamaModel struct {
	Name       string    `json:"name"`
	ModifiedAt time.Time `json:"modified_at"`
	Size       int64     `json:"size"`
	Digest     string    `json:"digest"`
	Details    struct {
		Format            string   `json:"format"`
		Family            string   `json:"family"`
		Families          []string `json:"families"`
		ParameterSize     string   `json:"parameter_size"`
		QuantizationLevel string   `json:"quantization_level"`
	} `json:"details"`
}

// OllamaModelsResponse represents the Ollama API response
type OllamaModelsResponse struct {
	Models []OllamaModel `json:"models"`
}

// GetOllamaHost returns the Ollama host URL from flag, env var, or default
func GetOllamaHost(flagHost string) string {
	// Priority: 1. Flag, 2. Env var, 3. Default
	if flagHost != "" {
		return flagHost
	}
	if envHost := os.Getenv("OLLAMA_HOST"); envHost != "" {
		return envHost
	}
	return "http://localhost:11434"
}

// FetchOllamaModels retrieves models from Ollama API
// Returns empty slice if server is unavailable (silent fail)
func FetchOllamaModels(host string) []Model {
	url := host + "/api/tags"

	client := &http.Client{
		Timeout: 3 * time.Second, // Short timeout for local server
	}

	resp, err := client.Get(url)
	if err != nil {
		// Silent fail - server not available
		return []Model{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Silent fail - API error
		return []Model{}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Model{}
	}

	var ollamaResp OllamaModelsResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return []Model{}
	}

	// Convert Ollama models to unified Model format
	models := make([]Model, 0, len(ollamaResp.Models))
	for _, om := range ollamaResp.Models {
		model := Model{
			ID:          "ollama/" + om.Name,
			Name:        om.Name,
			Created:     om.ModifiedAt.Unix(),
			Description: buildOllamaDescription(om),
			// Store Ollama-specific data for detailed view
			OllamaDetails: &OllamaDetails{
				Size:              om.Size,
				Format:            om.Details.Format,
				Family:            om.Details.Family,
				ParameterSize:     om.Details.ParameterSize,
				QuantizationLevel: om.Details.QuantizationLevel,
			},
		}
		models = append(models, model)
	}

	return models
}

// buildOllamaDescription creates a description from Ollama model details
func buildOllamaDescription(om OllamaModel) string {
	desc := ""

	if om.Details.Family != "" {
		desc = om.Details.Family
	}

	if om.Details.ParameterSize != "" {
		if desc != "" {
			desc += " "
		}
		desc += om.Details.ParameterSize
	}

	if om.Details.QuantizationLevel != "" {
		if desc != "" {
			desc += " "
		}
		desc += "(" + om.Details.QuantizationLevel + ")"
	}

	// Add size information
	sizeGB := float64(om.Size) / (1024 * 1024 * 1024)
	if sizeGB > 0 {
		if desc != "" {
			desc += " - "
		}
		desc += fmt.Sprintf("%.1f GB", sizeGB)
	}

	if desc == "" {
		desc = "Ollama local model"
	}

	return desc
}
