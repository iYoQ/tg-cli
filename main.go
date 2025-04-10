package main

import (
	"flag"
	"log"
	"os"
	"runtime/debug"
	"tg-cli/connection"
	"tg-cli/console"

	"github.com/chzyer/readline"
	"github.com/joho/godotenv"
)

func Init() error {
	apiId, apiHash := loadParams()

	conn := connection.NewConnection()

	err := Auth(apiId, apiHash, conn)
	if err != nil {
		conn.Close()
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic recovered: %v\n%s", r, debug.Stack())
		}

		if conn.Client != nil {
			conn.Close()
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

	console.Start(conn, reader)
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
