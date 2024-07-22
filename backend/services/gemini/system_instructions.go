package gemini

import (
	"fmt"
	"os"
	"read-robin/utils"

	"github.com/BurntSushi/toml"
)

type SystemInstructions struct {
	QuizModelSystemInstructions      string `toml:"quizModelSystemInstructions"`
	WebscrapeModelSystemInstructions string `toml:"webscrapeModelSystemInstructions"`
	PDFModelSystemInstructions       string `toml:"pdfModelSystemInstructions"`
	ReviewModelSystemInstructions    string `toml:"reviewModelSystemInstructions"`
}

var instructions SystemInstructions

func LoadSystemInstructions() {
	configFile, err := utils.FindConfigFile("gemini_system_instructions.toml")
	if err != nil {
		fmt.Printf("Error locating config file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Reading config file: %s\n", configFile)
	configData, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Unmarshaling config data")
	err = toml.Unmarshal(configData, &instructions)
	if err != nil {
		fmt.Printf("Error unmarshaling config file: %v\n", err)
		os.Exit(1)
	}
	fmt.Print("Loaded system instructions")
}

func init() {
	LoadSystemInstructions()
}
