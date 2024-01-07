# gat : a cat powered by Ollama AI

Fun note : Part of this README has been generated using this tool

This program is designed to summarize the content of a text file and 
generate a markdown-formatted response using Ollama API. 
The program allows users to specify various options, such as the model to use, 
temperature, and more. It also provides colorized output if requested.

WARNING:

When using an LLM don't take its response as the absolute truth ! 

In the next months we will have to take digital news as probabilities
and facts checking will be necessary.

## Requirements

A running ollama server (ollama serve) with the model pulled (ollama pull
mistral).

You can find and install ollama from [ollama website](https://ollama.ai/)

By default gat is configured to use mistral 7b model as it performs well 
on customer grade GPU with 8GB of RAM but you can override the model in the
configuration file or at runtime.

## Building gat

```bash
go get
go mod vendor
go build
```

Then you can copy the `gat` binary in your $PATH

## Command Line Arguments and Configuration
The application accepts command-line arguments using the `flag`
package. It recognizes the following flags:

- `-c` or `--color`: colorizes the output using terminal colors (default
is off)
- `-e` or `--explain`: asks for an explanation instead of a summary
- `-f` or `--format`: formats the output as markdown (default is false)
- `-h` or `--ollama_url`: specifies the Ollama API endpoint to use
(overrides config file)
- `-m` or `--model`: specifies the GPT model to use (overrides config file)
- `-p` or `--prompt`: specifies a custom prompt for the API request
- `-s` or `--stream`: requests the response to be streamed (default is off)
- `-t` or `--temperature`: sets the temperature of the generated output
(default is 0.5)

The application loads a configuration file called `config.json` in the 
`~/.config/gat/` directory if it exists.
If it doesn't a the directory will be created and a default `config.json` will
be created.
The configuration file contains settings such as the model, Ollama URL, 
temperature, and default prompt.
If the configuration file doesn't exist or the specified flags override its 
values, they are used instead.

## Usage
```bash
gat <filename> [-c] [-e] [-f] [-h ollama_url] \
[-t temperature] [-m model] [-p prompt]
```

## Examples

### French constitution explained

```bash
gat "samples/Déclaration des Droits de l'Homme et du Citoyen de 1789.txt" \
 -h http://localhost:11434/api/generate -e
```

Generated output :
------------------

```
The text you provided is a statement of the French National Assembly's
decision to declare certain fundamental human rights, which were
considered essential for maintaining a just society and preventing
corruption in government. The declaration was made on August 26, 1789,
during the French Revolution.

Here is an explanation of each article:

Article 1: All men are born free and equal in rights. Discrimination
based solely on utility cannot exist.

This article establishes that all people are entitled to certain
fundamental rights without any discrimination based on social status,
wealth, or any other factors. It sets the foundation for the protection
of individual liberties and equality under the law.

Article 2: The purpose of every political association is to preserve
the natural and imprescriptible rights of man. These rights are liberty,
property, security, and resistance to oppression.

This article ensures that any government or organization seeking power
must respect and protect individual rights. It also outlines the specific
rights that individuals have the right to, including freedom, ownership
of property, protection from harm, and resistance to oppression.

Article 3: Sovereignty lies in the nation. No authority can be exercised
without express permission from the people.

This article emphasizes the importance of popular sovereignty, meaning
that ultimate power resides with the people of a country. It ensures that
no government or organization can act on behalf of the people without
explicit authorization from them.

Article 4: Liberty consists in the ability to do anything that does not
harm others. The limits of individual freedoms are determined by law.

This article defines liberty as the freedom to act and express oneself, as
long as it does not infringe on the rights of others. It also establishes
that the limits of individual freedoms are established by law and must
be enforced in a just and equitable manner.

Article 5: Laws can only forbid actions harmful to society. Everything
else is permitted, and nothing can be prevented without legal
authorization.

This article emphasizes the importance of individual freedom and the
need to limit government intervention in people's lives. It also outlines
that laws must be just and proportional, meaning they must not infringe
on people's fundamental rights or liberties unnecessarily.

Article 6: No one can be punished for their opinions, as long as they
do not harm others.

This article protects freedom of speech and thought, ensuring that
individuals are free to express their opinions without fear of punishment,
as long as they do not infringe on the rights of others.

Article 7: The guarantee of human rights requires a public force. This
force is for the benefit of all, and not for the benefit of any individual
or group.

This article ensures that the government has the necessary resources to
protect individual rights and ensure justice for all citizens. It also
emphasizes that this force must be used solely for the benefit of society
as a whole, and not for the personal gain of any individual or group.

Article 8: All citizens have the right to vote on taxes and contributions
to the public force.

This article gives citizens the power to hold their government accountable
by requiring that they vote on taxes and other contributions to the public
force. It also ensures that this process is transparent and fair, with
the opportunity for citizens to participate in decision-making processes
that affect their lives.

Article 9: Any society that does not have a constitution guaranteeing
individual rights and limiting government power has no government at all.

This article emphasizes the importance of constitutional guarantees
of individual rights and limitations on government power, as they
are essential for maintaining a just and equitable society. It also
underscores that a government without such guarantees is not truly a
government at all.
```

### Explanation of this program code by starling-lm

```bash
gat main.go -h http://localhost:11434/api/generate -e -m starling-lm -f 
```

Generated output :
------------------

 Here's a detailed explanation of the code provided in markdown format:

## Overview
This Go program is designed to summarize the content of a text file and
generate a markdown-formatted response using OpenAI's GPT-4 model. The
program allows users to specify various options, such as the model to
use, temperature, streaming mode, and more. It also provides colorized
output if requested.

## Imported Packages
The code imports several packages that are used throughout the program:

- `bytes`: Provides functions related to working with byte slices.
- `unicode`: Contains various constants and functions related to Unicode
characters and encoding.
- `utf8`: Contains functions for manipulating UTF-8 strings.
- `encoding/json`: Contains packages and types for JSON marshaling and
unmarshaling.
- `flag`: Provides a package to parse command line flags.
- `fmt`: Standard Go package for formatted I/O functions.
- `net/http`: Package for making HTTP requests.
- `os`: Contains platform-specific functions related to operating systems.
- `os/user`: Provides the current user information.
- `path/filepath`: Contains functions related to working with file paths.
- `strings`: Standard Go package for string manipulation functions.
- `github.com/fatih/color`: A popular package in Go for colorizing
terminal output.

## Configuration Structure
The code defines two structs, `Config` and `Request`, to represent the
configuration and request structures:

```go
type Config struct {
    Model       string `json:"model"`
    OllamaURL   string `json:"ollama_url"`
    Temperature float64 `json:"temperature"`
    Prompt      string `json:"prompt"`
    Stream      bool   `json:"stream"`
    Columns     int    `json:"columns"`
}
```

```go
type Request struct {
    Model       string `json:"model"`
    Prompt      string `json:"prompt"`
    Stream      bool   `json:"stream"`
    Temperature float64 `json:"temperature,omitempty"`
}
```

These structs are used to define the various fields required by the program.

## Function Breakdown

### `wordWrap(text string, lineWidth int) string`
This function takes a text string and a target line width as input and
wraps the text to fit within the specified line width. It uses a custom
algorithm to ensure proper word wrapping and returns the wrapped text
as a string.

### `main()` Function
The `main()` function is the entry point of the program:

1. It first checks if any arguments are provided by the user, and if not,
it displays usage instructions and exits.
2. The function reads the command line flags using the `flag` package
and stores their values in appropriate variables.
3. It loads a configuration file (`config.json`) from the `.config/gat/`
directory using the `loadConfig()` function. If an error occurs, it logs
the error message and exits.
4. The function sets default values for various fields if they are not
provided by the user through flags.
5. It reads the contents of the input file specified by the user and
stores it in a variable.
6. Based on the `explain` flag, it constructs a base prompt string to
ask for either a detailed explanation or a summary of the content using
GPT-4 model.
7. It builds the request payload as a JSON object using the values from
flags and variables.
8. The function sends a POST request to the OpenAI API with the prepared
JSON payload. If any error occurs, it logs the error message and exits.
9. It receives the API response in JSON format and decodes it into a
`Message` struct. If any error occurs, it logs the error message and
exits.
10. The function processes the output message based on user flags
and calls the `colorizeMarkdown()` function to colorize the output if
requested. Finally, it prints the generated markdown-formatted response.

### `loadConfig() (Config, error)` Function
This function loads the configuration from a JSON file named `config.json`
in the `.config/gat/` directory:

1. It creates a `user` object to get the current user's home directory.
2. It constructs the path of the config directory and file.
3. If the directory does not exist, it creates the directory and writes
a sample configuration file (`config.json`) with default values.
4. If the file exists, it reads the contents of the file and unmarshals
it into a `Config` struct. If any error occurs during this process,
it logs the error message and returns an empty `Config` object.
5. The function returns the loaded configuration along with any error
that occurred during the process.

### `colorizeMarkdown(content string) string` Function
This function colorizes a given markdown-formatted content using terminal
colors provided by the `github.com/fatih/color` package:

1. It splits the input content into lines and stores them in an array.
2. For each line, it uses the `colorizeLine()` function to colorize
the text based on its type (heading, emphasis, strong, or code). The
resulting colored lines are stored in another array.
3. Finally, the function joins the colored lines with newline characters
and returns the colored markdown-formatted content.

### `colorizeLine(line string) string` Function
This function colorizes a single line of markdown-formatted content
using terminal colors provided by the `github.com/fatih/color` package:

1. It replaces specific markdown elements (headings, emphases, strong
texts, and codes) in the input line with their respective colored
representations using functions from the `github.com/fatih/color` package.
2. The function returns the colorized line.

## Usage Instructions
To use this program, follow these steps:

1. Make sure you have Go installed on your system. You can download and
install it from [the official website](https://golang.org/doc/install).
2. Save the provided code in a file named `main.go`.
3. Open a terminal window and navigate to the directory where you saved
the `main.go` file.
4. Run the following command to build and execute the program:
   ```bash
   go run main.go [filename] [-c] [-s] [-e] [-h ollama_url] \
   [-t temperature] [-m model] [-p prompt]
   ```
5. Replace `[filename]` with the path to the input text file you want
to summarize. You can also specify various flags, such as `-c` for
colorized output, `-s` for streaming response, `-e` for explaining the
content instead of summarizing it, `-h` for specifying a custom Ollama API
endpoint, `-t` for setting the temperature of the generated response, `-m`
for specifying the model to use, and `-p` for providing a custom prompt.

After executing the command, the program will generate a
markdown-formatted summary or explanation of the specified text file
based on the provided flags and settings.

