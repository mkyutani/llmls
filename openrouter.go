package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"golang.org/x/term"
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

// GetTerminalWidth returns the terminal width, or a default value if unavailable
func GetTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// Default to 120 if terminal size cannot be determined
		return 120
	}
	return width
}

// CalculateDescriptionWidth calculates the available width for description
func CalculateDescriptionWidth(termWidth, modelWidth, providerWidth int) int {
	// Column layout: modelID (2 spaces) provider (2 spaces) date (2 spaces) description
	// Date is always 10 characters (YYYY-MM-DD)
	dateWidth := 10
	spacingWidth := 6 // 3 separators * 2 spaces each
	usedWidth := modelWidth + providerWidth + dateWidth + spacingWidth

	descWidth := termWidth - usedWidth

	// Minimum description width
	if descWidth < 30 {
		descWidth = 30
	}

	return descWidth
}

// DisplayModels prints models in formatted output with dynamic column widths
func DisplayModels(models []Model) {
	if len(models) == 0 {
		return
	}

	// Get terminal width
	termWidth := GetTerminalWidth()

	// Calculate maximum widths for model and provider columns
	maxModelWidth := 0
	maxProviderWidth := 0

	for _, model := range models {
		if len(model.ID) > maxModelWidth {
			maxModelWidth = len(model.ID)
		}
		provider := ExtractProvider(model.ID)
		if len(provider) > maxProviderWidth {
			maxProviderWidth = len(provider)
		}
	}

	// Calculate available width for description
	descWidth := CalculateDescriptionWidth(termWidth, maxModelWidth, maxProviderWidth)

	// Display each model with dynamic column widths
	for _, model := range models {
		provider := ExtractProvider(model.ID)
		date := FormatDate(model.Created)
		desc := TruncateDescription(model.Description, descWidth)

		// Format with dynamic widths: model_id | provider | date | description
		fmt.Printf("%-*s  %-*s  %s  %s\n",
			maxModelWidth, model.ID,
			maxProviderWidth, provider,
			date, desc)
	}
}
