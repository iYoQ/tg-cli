package main

import (
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"tg-cli/app"
	"tg-cli/connection"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

type Config struct {
	apiId   int32
	apiHash string
}

type flags struct {
	chatIdFlag  *string
	fileFlag    *string
	photoFlag   *string
	captionFlag *string
}

func start() error {
	cfg := loadParams()
	flags := loadFlags()

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

	if ok, err := checkFlags(conn, flags); err != nil {
		return err
	} else if ok {
		return nil
	}

	app := tea.NewProgram(app.NewRootModel(conn), tea.WithAltScreen())
	if _, err := app.Run(); err != nil {
		return err
	}

	return nil
}

func loadParams() Config {
	godotenv.Load()
	apiIdRaw := os.Getenv("API_ID")
	apiHash := os.Getenv("API_HASH")

	if apiIdRaw == "" || apiHash == "" {
		log.Fatalf("API_ID and API_HASH are required, use .env file, or ENV")
	}

	apiId64, err := strconv.ParseInt(apiIdRaw, 10, 64)
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
		log.Fatalf("Error: %s", err)
	}
}
