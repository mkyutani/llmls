package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

const openRouterModelsURL = "https://openrouter.ai/api/v1/models"

// Model represents an OpenRouter model
type Model struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Created     int64  `json:"created"`
	Description string `json:"description"`
}

// ModelsResponse represents the API response structure
type ModelsResponse struct {
	Data []Model `json:"data"`
}

// FetchModels retrieves models from OpenRouter API
func FetchModels() ([]Model, error) {
	resp, err := http.Get(openRouterModelsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var modelsResp ModelsResponse
	if err := json.Unmarshal(body, &modelsResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return modelsResp.Data, nil
}

// FilterModels filters models by provider and model name (case-insensitive partial match)
func FilterModels(models []Model, providerFilter, modelFilter string) []Model {
	if providerFilter == "" && modelFilter == "" {
		return models
	}

	var filtered []Model
	providerLower := strings.ToLower(providerFilter)
	modelLower := strings.ToLower(modelFilter)

	for _, model := range models {
		// Extract provider from ID (format: "provider/model-name")
		provider := ""
		if idx := strings.Index(model.ID, "/"); idx > 0 {
			provider = model.ID[:idx]
		}

		providerMatch := providerFilter == "" || strings.Contains(strings.ToLower(provider), providerLower)
		modelMatch := modelFilter == "" || strings.Contains(strings.ToLower(model.ID), modelLower) || strings.Contains(strings.ToLower(model.Name), modelLower)

		if providerMatch && modelMatch {
			filtered = append(filtered, model)
		}
	}

	return filtered
}

// SortModelsByCreatedDesc sorts models by creation date in descending order
func SortModelsByCreatedDesc(models []Model) {
	sort.Slice(models, func(i, j int) bool {
		return models[i].Created > models[j].Created
	})
}

// FormatDate converts Unix timestamp to YYYY-MM-DD format in local timezone
func FormatDate(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02")
}

// TruncateDescription truncates description to maxLen characters, adding ".." if truncated
func TruncateDescription(desc string, maxLen int) string {
	if len(desc) <= maxLen {
		return desc
	}
	return desc[:maxLen] + ".."
}

// ExtractProvider extracts provider name from model ID
func ExtractProvider(modelID string) string {
	if idx := strings.Index(modelID, "/"); idx > 0 {
		return modelID[:idx]
	}
	return "Unknown"
}

// DisplayModels prints models in formatted output
func DisplayModels(models []Model) {
	for _, model := range models {
		provider := ExtractProvider(model.ID)
		date := FormatDate(model.Created)
		desc := TruncateDescription(model.Description, 98)

		// Format: model_name (padded) | provider (padded) | date | description
		fmt.Printf("%-30s %-20s %s  %s\n", model.ID, provider, date, desc)
	}
}
