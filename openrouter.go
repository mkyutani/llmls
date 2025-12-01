package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"golang.org/x/term"
)

const openRouterModelsURL = "https://openrouter.ai/api/v1/models"

// Model represents an OpenRouter model
type Model struct {
	ID             string       `json:"id"`
	Name           string       `json:"name"`
	Created        int64        `json:"created"`
	Description    string       `json:"description"`
	ContextLength  int          `json:"context_length"`
	Architecture   Architecture `json:"architecture"`
	Pricing        Pricing      `json:"pricing"`
	TopProvider    TopProvider  `json:"top_provider"`
	OllamaDetails  *OllamaDetails `json:"-"` // Ollama-specific details (not from JSON)
}

// Architecture represents model architecture details
type Architecture struct {
	Modality        string   `json:"modality"`
	InputModalities []string `json:"input_modalities"`
	OutputModalities []string `json:"output_modalities"`
	Tokenizer       string   `json:"tokenizer"`
}

// Pricing represents model pricing information
type Pricing struct {
	Prompt            string `json:"prompt"`
	Completion        string `json:"completion"`
	Request           string `json:"request"`
	Image             string `json:"image"`
	WebSearch         string `json:"web_search"`
	InternalReasoning string `json:"internal_reasoning"`
}

// TopProvider represents top provider details
type TopProvider struct {
	ContextLength        int  `json:"context_length"`
	MaxCompletionTokens  int  `json:"max_completion_tokens"`
	IsModerated          bool `json:"is_moderated"`
}

// OllamaDetails represents Ollama-specific model details
type OllamaDetails struct {
	Size              int64
	Format            string
	Family            string
	ParameterSize     string
	QuantizationLevel string
}

// ModelsResponse represents the API response structure
type ModelsResponse struct {
	Data []Model `json:"data"`
}

// globMatch performs case-insensitive glob pattern matching
// Supports * (any sequence including /) and ? (single character including /)
func globMatch(pattern, str string) bool {
	// Convert pattern to regex
	// Escape special regex characters except * and ?
	regexPattern := regexp.QuoteMeta(pattern)
	// Replace escaped glob wildcards with regex equivalents
	regexPattern = strings.ReplaceAll(regexPattern, "\\*", ".*")
	regexPattern = strings.ReplaceAll(regexPattern, "\\?", ".")
	// Anchor pattern to match entire string
	regexPattern = "^" + regexPattern + "$"

	// Case-insensitive match
	re, err := regexp.Compile("(?i)" + regexPattern)
	if err != nil {
		// If pattern is invalid, fall back to case-insensitive exact match
		return strings.EqualFold(pattern, str)
	}
	return re.MatchString(str)
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

// FilterModels filters models by model ID using glob patterns
// Supports * (any sequence) and ? (single character) in patterns
// Also supports exact match against provider name (case-insensitive)
// If pattern is empty, returns all models
func FilterModels(models []Model, pattern string) []Model {
	// If no pattern, return all models
	if pattern == "" {
		return models
	}

	var filtered []Model
	for _, model := range models {
		provider := ExtractProvider(model.ID)

		// Match against:
		// 1. Model ID (glob pattern)
		// 2. Model name (glob pattern)
		// 3. Provider name (exact match, case-insensitive)
		if globMatch(pattern, model.ID) ||
		   globMatch(pattern, model.Name) ||
		   strings.EqualFold(pattern, provider) {
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

// DisplayProviders prints unique provider names
func DisplayProviders(models []Model) {
	if len(models) == 0 {
		return
	}

	// Extract unique providers
	providerSet := make(map[string]bool)
	for _, model := range models {
		provider := ExtractProvider(model.ID)
		providerSet[provider] = true
	}

	// Convert to slice
	var providers []string
	for provider := range providerSet {
		providers = append(providers, provider)
	}

	// Sort alphabetically
	sort.Strings(providers)

	// Display providers (one per line)
	for _, provider := range providers {
		fmt.Println(provider)
	}
}

// DisplayModelsDetailed prints comprehensive information for each model
func DisplayModelsDetailed(models []Model) {
	if len(models) == 0 {
		return
	}

	for i, model := range models {
		if i > 0 {
			fmt.Println() // Blank line between models
		}

		provider := ExtractProvider(model.ID)
		date := FormatDate(model.Created)

		// Print separator line
		fmt.Println(strings.Repeat("=", 80))

		// Basic information
		fmt.Printf("Model ID:          %s\n", model.ID)
		fmt.Printf("Name:              %s\n", model.Name)
		fmt.Printf("Provider:          %s\n", provider)
		fmt.Printf("Created:           %s\n", date)

		// Technical details
		if model.ContextLength > 0 {
			fmt.Printf("Context Length:    %s tokens\n", FormatNumber(model.ContextLength))
		}
		if model.TopProvider.MaxCompletionTokens > 0 {
			fmt.Printf("Max Completion:    %s tokens\n", FormatNumber(model.TopProvider.MaxCompletionTokens))
		}

		// Architecture
		if model.Architecture.Modality != "" {
			fmt.Printf("Modality:          %s\n", model.Architecture.Modality)
		}

		// Pricing information
		if model.Pricing.Prompt != "" && model.Pricing.Prompt != "0" {
			promptPrice := FormatPrice(model.Pricing.Prompt)
			completionPrice := FormatPrice(model.Pricing.Completion)
			fmt.Printf("Pricing:           $%s / 1K prompt tokens, $%s / 1K completion tokens\n",
				promptPrice, completionPrice)
		}

		// Moderation
		if model.TopProvider.IsModerated {
			fmt.Println("Moderation:        Enabled")
		}

		// Ollama-specific details
		if model.OllamaDetails != nil {
			if model.OllamaDetails.Family != "" {
				fmt.Printf("Model Family:      %s\n", model.OllamaDetails.Family)
			}
			if model.OllamaDetails.ParameterSize != "" {
				fmt.Printf("Parameter Size:    %s\n", model.OllamaDetails.ParameterSize)
			}
			if model.OllamaDetails.QuantizationLevel != "" {
				fmt.Printf("Quantization:      %s\n", model.OllamaDetails.QuantizationLevel)
			}
			if model.OllamaDetails.Format != "" {
				fmt.Printf("Format:            %s\n", model.OllamaDetails.Format)
			}
			if model.OllamaDetails.Size > 0 {
				sizeGB := float64(model.OllamaDetails.Size) / (1024 * 1024 * 1024)
				fmt.Printf("Model Size:        %.2f GB\n", sizeGB)
			}
		}

		// Description (full, not truncated)
		if model.Description != "" {
			fmt.Println("Description:")
			// Wrap description text at terminal width
			termWidth := GetTerminalWidth()
			descWidth := termWidth - 4 // Leave margin
			if descWidth < 40 {
				descWidth = 40
			}
			wrappedDesc := WrapText(model.Description, descWidth)
			for _, line := range wrappedDesc {
				fmt.Printf("  %s\n", line)
			}
		}

		fmt.Println(strings.Repeat("=", 80))
	}
}

// FormatNumber formats a number with thousand separators
func FormatNumber(n int) string {
	str := fmt.Sprintf("%d", n)
	result := ""
	for i, c := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += ","
		}
		result += string(c)
	}
	return result
}

// FormatPrice formats a price string to a readable format
func FormatPrice(price string) string {
	// Convert from per-token to per-1K-tokens
	if price == "" || price == "0" {
		return "0"
	}

	// Parse as float
	p := 0.0
	fmt.Sscanf(price, "%f", &p)

	// Multiply by 1000 to get per-1K price
	p1k := p * 1000

	// Format with appropriate precision
	if p1k >= 1 {
		return fmt.Sprintf("%.3f", p1k)
	} else if p1k >= 0.01 {
		return fmt.Sprintf("%.4f", p1k)
	} else {
		return fmt.Sprintf("%.6f", p1k)
	}
}

// WrapText wraps text to specified width, breaking at word boundaries
func WrapText(text string, width int) []string {
	// Replace newline characters with spaces
	text = strings.ReplaceAll(text, "\r\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")
	text = strings.ReplaceAll(text, "\n", " ")

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}

	var lines []string
	currentLine := words[0]

	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	lines = append(lines, currentLine)

	return lines
}

