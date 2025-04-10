package main

import (
	"flag"
	"log"
	"os"
	"runtime/debug"
	"tg-cli/console"
	"tg-cli/handlers"

	"github.com/chzyer/readline"
	"github.com/joho/godotenv"
	tdlib "github.com/zelenin/go-tdlib/client"
)

func Init() error {
	apiId, apiHash := loadParams()
	updatesChannel := make(chan *tdlib.Message)

	client, err := Auth(apiId, apiHash, updatesChannel)
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

	reader, err := readline.NewEx(&readline.Config{
		Prompt:          ">> ",
		InterruptPrompt: "back",
		EOFPrompt:       "exit",
	})

	if err != nil {
		log.Printf("failed to initialize readline: %v", err)
		return err
	}
	defer reader.Close()

	console.Start(client, reader, updatesChannel)
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
