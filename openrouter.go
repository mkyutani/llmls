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

// FilterModels filters models by provider, model name, description, and unified search term (case-insensitive partial match)
// If explicit filters (provider, model, description) are provided, they take precedence over searchTerm
// searchTerm performs OR matching across model ID, provider name, and description
func FilterModels(models []Model, providerFilter, modelFilter, descriptionFilter, searchTerm string) []Model {
	// If using explicit filters, ignore search term
	hasExplicitFilters := providerFilter != "" || modelFilter != "" || descriptionFilter != ""

	// If no filters at all, return all models
	if !hasExplicitFilters && searchTerm == "" {
		return models
	}

	var filtered []Model

	// Use search term if no explicit filters are provided
	if !hasExplicitFilters && searchTerm != "" {
		searchLower := strings.ToLower(searchTerm)

		for _, model := range models {
			// Extract provider from ID (format: "provider/model-name")
			provider := ExtractProvider(model.ID)

			// OR matching: search term matches any field
			modelMatch := strings.Contains(strings.ToLower(model.ID), searchLower) || strings.Contains(strings.ToLower(model.Name), searchLower)
			providerMatch := strings.Contains(strings.ToLower(provider), searchLower)
			descriptionMatch := strings.Contains(strings.ToLower(model.Description), searchLower)

			if modelMatch || providerMatch || descriptionMatch {
				filtered = append(filtered, model)
			}
		}

		return filtered
	}

	// Use explicit filters (AND matching)
	providerLower := strings.ToLower(providerFilter)
	modelLower := strings.ToLower(modelFilter)
	descriptionLower := strings.ToLower(descriptionFilter)

	for _, model := range models {
		// Extract provider from ID (format: "provider/model-name")
		provider := ExtractProvider(model.ID)

		providerMatch := providerFilter == "" || strings.Contains(strings.ToLower(provider), providerLower)
		modelMatch := modelFilter == "" || strings.Contains(strings.ToLower(model.ID), modelLower) || strings.Contains(strings.ToLower(model.Name), modelLower)
		descriptionMatch := descriptionFilter == "" || strings.Contains(strings.ToLower(model.Description), descriptionLower)

		if providerMatch && modelMatch && descriptionMatch {
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
// Also replaces newline characters with spaces
func TruncateDescription(desc string, maxLen int) string {
	// Replace newline characters with spaces
	desc = strings.ReplaceAll(desc, "\r\n", " ")
	desc = strings.ReplaceAll(desc, "\r", " ")
	desc = strings.ReplaceAll(desc, "\n", " ")

	// Convert to rune slice to handle multi-byte characters correctly
	runes := []rune(desc)
	if len(runes) <= maxLen {
		return desc
	}
	return string(runes[:maxLen]) + ".."
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
	// Column layout: modelID (1 space) provider (1 space) date (1 space) description
	// Date is always 10 characters (YYYY-MM-DD)
	dateWidth := 10
	spacingWidth := 3 // 3 separators * 1 space each
	safetyMargin := 5 // Safety margin to prevent line wrapping
	usedWidth := modelWidth + providerWidth + dateWidth + spacingWidth + safetyMargin

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
		fmt.Printf("%-*s %-*s %s %s\n",
			maxModelWidth, model.ID,
			maxProviderWidth, provider,
			date, desc)
	}
}

// DisplayProviders prints unique provider names, optionally filtered
func DisplayProviders(models []Model, filter string) {
	if len(models) == 0 {
		return
	}

	// Extract unique providers
	providerSet := make(map[string]bool)
	for _, model := range models {
		provider := ExtractProvider(model.ID)
		providerSet[provider] = true
	}

	// Convert to slice and filter
	var providers []string
	filterLower := strings.ToLower(filter)
	for provider := range providerSet {
		if filter == "" || strings.Contains(strings.ToLower(provider), filterLower) {
			providers = append(providers, provider)
		}
	}

	// Sort alphabetically
	sort.Strings(providers)

	// Calculate max width
	maxWidth := 0
	for _, provider := range providers {
		if len(provider) > maxWidth {
			maxWidth = len(provider)
		}
	}

	// Display providers
	for _, provider := range providers {
		fmt.Printf("%-*s\n", maxWidth, provider)
	}
}

// DisplayModelsFiltered prints models with provider names, filtered by model name
func DisplayModelsFiltered(models []Model, filter string) {
	if len(models) == 0 {
		return
	}

	// Filter models by name (case-insensitive)
	var filtered []Model
	filterLower := strings.ToLower(filter)
	for _, model := range models {
		if filter == "" || strings.Contains(strings.ToLower(model.ID), filterLower) || strings.Contains(strings.ToLower(model.Name), filterLower) {
			filtered = append(filtered, model)
		}
	}

	if len(filtered) == 0 {
		return
	}

	// Sort by creation date descending
	SortModelsByCreatedDesc(filtered)

	// Get terminal width
	termWidth := GetTerminalWidth()

	// Calculate maximum widths for model and provider columns
	maxModelWidth := 0
	maxProviderWidth := 0

	for _, model := range filtered {
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
	for _, model := range filtered {
		provider := ExtractProvider(model.ID)
		date := FormatDate(model.Created)
		desc := TruncateDescription(model.Description, descWidth)

		// Format with dynamic widths: model_id | provider | date | description
		fmt.Printf("%-*s %-*s %s %s\n",
			maxModelWidth, model.ID,
			maxProviderWidth, provider,
			date, desc)
	}
}
