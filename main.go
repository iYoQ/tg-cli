package main

import (
	"flag"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"tg-cli/connection"
	"tg-cli/view"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

type Config struct {
	apiId   int32
	apiHash string
}

func start() error {
	cfg := loadParams()

	conn := connection.NewConnection()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic recovered: %v\n%s", r, debug.Stack())
		}

		if conn.Client != nil {
			conn.Close()
		}
	}()

	if err := auth(cfg, conn); err != nil {
		return err
	}

	p := tea.NewProgram(view.NewModel(conn))
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

func loadParams() Config {
	godotenv.Load()

	apiIdFlag := flag.String("id", "", "api id")
	apiHashFlag := flag.String("hash", "", "api hash")
	flag.Parse()

	apiIdRaw := *apiIdFlag
	if apiIdRaw == "" {
		apiIdRaw = os.Getenv("API_ID")
	}

	apiHash := *apiHashFlag
	if apiHash == "" {
		apiHash = os.Getenv("API_HASH")
	}

	if apiIdRaw == "" || apiHash == "" {
		log.Fatalf("API_ID and API_HASH are required, use --id=, --hash= flags or .env file, or ENV")
	}

	apiId64, err := strconv.ParseInt(apiIdRaw, 10, 32)
	if err != nil {
		log.Fatalf("strconv.Atoi error: %s", err)
	}

	apiId := int32(apiId64)

	return Config{
		apiId:   apiId,
		apiHash: apiHash,
	}
}

func main() {
	if err := start(); err != nil {
		log.Fatalf("Initialization failed: %s", err)
	}
}
