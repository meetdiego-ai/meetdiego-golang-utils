package openai

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

var (
	openaiClient *openai.Client
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: Error loading .env file: %v\n", err)
	}

	if os.Getenv("OPENAI_API_KEY") == "" {
		panic("Error: OPENAI_API_KEY is not set")
	}

	openaiClient = openai.NewClient(os.Getenv("OPENAI_API_KEY"))
}
