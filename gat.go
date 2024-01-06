package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// Config represents the configuration structure
type Config struct {
	Model     string `json:"model"`
	OllamaURL string `json:"ollama_url"`
}

// Message represents the JSON response structure expected from the API
type Message struct {
	Result string `json:"content"`
}

func main() {
  args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Usage: gat <filename> [-c] [-m model]")
		return
	}

	filename := args[0]
	colorize := false
	var model string

	flagSet := flag.NewFlagSet("gat", flag.ExitOnError)
	flagSet.BoolVar(&colorize, "c", false, "Colorize the output")
	flagSet.StringVar(&model, "m", "", "Specify the model (overrides config file)")
	flagSet.Parse(args[1:])

	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %s\n", err)
		return
	}

	if model != "" {
		config.Model = model
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		return
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}

	// Parse Markdown content
	fileContent := buf.String()
	var output string

	jsonData, err := json.Marshal(map[string]string{"content": fileContent}) // Send the original content to API
	if err != nil {
		fmt.Printf("Error encoding JSON: %s\n", err)
		return
	}

	resp, err := http.Post(config.OllamaURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error sending data to API: %s\n", err)
		return
	}
	defer resp.Body.Close()

	var message Message
	err = json.NewDecoder(resp.Body).Decode(&message)
	if err != nil {
		fmt.Printf("Error decoding API response: %s\n", err)
		return
	}
  
  if colorize {
	  fmt.Println("Colorized output")
		output = colorizeMarkdown(message.Result)
	} else {
	  fmt.Println("Raw output")
		output = message.Result
	}
	fmt.Println(output)
}

// colorizeMarkdown colorizes the Markdown content using terminal colors
func colorizeMarkdown(content string) string {
	lines := strings.Split(content, "\n")
	var coloredLines []string

	for _, line := range lines {
		coloredLines = append(coloredLines, colorizeLine(line))
	}

	return strings.Join(coloredLines, "\n")
}

// colorizeLine colorizes a single line of Markdown content
func colorizeLine(line string) string {
	// You can define various colors and styles for different Markdown elements here
	// For example:
	// Heading - Magenta, Bold
	// Emphasis (Italic) - Cyan
	// Strong (Bold) - Green
	// Code - Yellow

	// This is just an example. You can modify these colors/styles based on your preferences.
	line = strings.ReplaceAll(line, "#", color.MagentaString("#"))
	line = strings.ReplaceAll(line, "_", color.CyanString("_"))
	line = strings.ReplaceAll(line, "**", color.GreenString("**"))
	line = strings.ReplaceAll(line, "`", color.YellowString("`"))

	return line
}

// loadConfig loads the configuration from config.json in the .config/gat/ directory
func loadConfig() (Config, error) {
	var config Config

	usr, err := user.Current()
	if err != nil {
		return config, err
	}

	configDir := filepath.Join(usr.HomeDir, ".config", "gat")
	configPath := filepath.Join(configDir, "config.json")

	// Create .config/gat directory if it doesn't exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, 0755)
		if err != nil {
			return config, err
		}

		// Create and write default config.json if it doesn't exist
		defaultConfig := Config{
			Model:     "mistral",
			OllamaURL: "http://127.0.0.1:11434",
		}
		defaultJSON, err := json.MarshalIndent(defaultConfig, "", "    ")
		if err != nil {
			return config, err
		}

		err = os.WriteFile(configPath, defaultJSON, 0644)
		if err != nil {
			return config, err
		}

		// Load default config
		return defaultConfig, nil
	}

	// Read config.json
	file, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

