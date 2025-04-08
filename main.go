package main

import (
	"flag"
	"log"
	"os"
	"runtime/debug"
	"tg-cli/console"
	"tg-cli/handlers"

	"github.com/joho/godotenv"
)

func Init() error {
	apiId, apiHash := loadParams()

	client, err := Auth(apiId, apiHash)
	if err != nil {
		if client != nil {
			handlers.ShutDown(client)
		}
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic recovered: %v\n%s", r, debug.Stack())
		}

		if client != nil {
			handlers.ShutDown(client)
		}
	}()

	console.Start(client)
	return nil
}

func loadParams() (string, string) {
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
		log.Printf("API_ID and API_HASH are required, use --id=, --hash= flags or .env file, or ENV")
	}

	return apiId, apiHash
}

func main() {
	if err := Init(); err != nil {
		log.Fatalf("Initialization failed: %s", err)
	}
}
