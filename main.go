package main

import (
	"flag"
	"log"
	"os"
	manual "tg-cli/manualmanager"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	apiIdFlag := flag.String("id", "", "api id")
	apiHashFlag := flag.String("hash", "", "api hash")

	apiId := *apiIdFlag
	if apiId == "" {
		apiId = os.Getenv("API_ID")
	}

	apiHash := *apiHashFlag
	if apiHash == "" {
		apiHash = os.Getenv("API_HASH")
	}

	if apiId == "" || apiHash == "" {
		log.Fatal("API_ID and API_HASH are required, use --id=, --hash= flags or .env file, or ENV")
	}

	my_client := Auth(apiId, apiHash)

	manual.Start(my_client)
}
