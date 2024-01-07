package main

import (
	"bytes"
  "unicode"
  "unicode/utf8"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// Config represents the configuration structure
type Config struct {
	Model       string `json:"model,omitempty"`
	OllamaURL   string `json:"ollama_url,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	Prompt      string `json:"prompt,omitempty"`
  Stream      bool `json:"stream,omitempty"` 
  Columns     int  `json:"columns,omitempty"`
}

// Request represents the request structure
type Request struct {
	Model       string `json:"model"`
	Prompt      string `json:"prompt"`
  Stream      bool `json:"stream"` 
	Temperature float64 `json:"temperature,omitempty"`
}

// Message represents the JSON response structure expected from the API
type Message struct {
	Response string `json:"response"`
}


// TODO: Word wrap response at lineWidth
func wordWrap(text string, lineWidth int) string {
    wrap := make([]byte, 0, len(text)+2*len(text)/lineWidth)
    eoLine := lineWidth
    inWord := false
    for i, j := 0, 0; ; {
        r, size := utf8.DecodeRuneInString(text[i:])
        if size == 0 && r == utf8.RuneError {
            r = ' '
        }
        if unicode.IsSpace(r) {
            if inWord {
                if i >= eoLine {
                    wrap = append(wrap, '\n')
                    eoLine = len(wrap) + lineWidth
                } else if len(wrap) > 0 {
                    wrap = append(wrap, ' ')
                }
                wrap = append(wrap, text[j:i]...)
            }
            inWord = false
        } else if !inWord {
            inWord = true
            j = i
        }
        if size == 0 && r == ' ' {
            break
        }
        i += size
    }
    return string(wrap)
}

func main() {
  args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Usage: gat <filename> [-c] [-e] [-h ollama_url] [-t temperature] [-m model] [-p prompt]")
		return
	}

	filename := args[0]
	colorize := false
	stream := false
  explain := false
	var model string
	var temperature float64
	var format bool = false
	var prompt string
	var host string
  var fileContent string
  var ollamaUrl string
  var basePrompt string
  var language string

	flagSet := flag.NewFlagSet("gat", flag.ExitOnError)
	flagSet.BoolVar(&colorize, "c", false, "Colorize the output")
  // TODO: Handle streamed responses
	// flagSet.BoolVar(&stream, "s", false, "Request the response to be streamed")
	flagSet.BoolVar(&explain, "e", false, "Ask for explaination instead of summary")
	flagSet.BoolVar(&format, "f", false, "Format the output as markdown)")
	flagSet.Float64Var(&temperature, "t", 0.5, "Specify the temperature (default 0.5)")
	flagSet.StringVar(&language, "l", "", "Ask model to output in the defined language")
	flagSet.StringVar(&host, "h", "", "Specify the ollama host endpoint (overrides config file)")
	flagSet.StringVar(&model, "m", "", "Specify the model (overrides config file)")
	flagSet.StringVar(&prompt, "p", "", "Specify a prompt")
	flagSet.Parse(args[1:])

	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %s\n", err)
		return
	}

	if model == "" {
		model = config.Model
  }

	if host == "" {
		ollamaUrl = config.OllamaURL
  } else {
		ollamaUrl = host
  }

	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		return
	}

	var output string

  fileContent = string(file)


  if explain {
    basePrompt = "explain in details, with usage instruction if it's some code, the following content :"
  } else  {
	  if prompt == "" {
	  	basePrompt = config.Prompt
	  } else {
      basePrompt = prompt
    }
  }

  var langPrompt string

  if language != "" {
    langPrompt = "Your response must be in " + language
  } else {
    langPrompt = ""
  }

  var finalPrompt string

  if format {
    finalPrompt = "Your response must be in markdown format." +
      langPrompt + basePrompt + "\n" + fileContent
  } else {
    finalPrompt = langPrompt + basePrompt + "\n" + fileContent
    colorize = false
  }

  request := Request{
    Model: model,
    Prompt: finalPrompt,
    Stream: stream,
    Temperature: temperature,
  }

	requestJSON, err := json.Marshal(request)

	if err != nil {
		fmt.Printf("Error encoding JSON: %s\n", err)
		return
	}

	resp, err := http.Post(ollamaUrl,
                         "application/json",
                         bytes.NewBuffer(requestJSON))

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
		output = colorizeMarkdown(message.Response)
	} else {
		output = message.Response
	}
  fmt.Println(output)
  // FIXME: handle correct word wrapping
  // fmt.Printf("%s\n\n", wordWrap(output, 80))
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
      Model:     "neural-chat:7b",
			OllamaURL: "http://127.0.0.1:11434/api/generate",
      Temperature: 0.5,
      Stream: false,
      Prompt: "summarize the following content :",
      Columns: 80,
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

