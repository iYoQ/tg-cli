package main

import (
	"flag"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"tg-cli/app"
	"tg-cli/connection"
	"tg-cli/requests"

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

	if checkFlags(conn) {
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

func checkFlags(conn *connection.Connection) bool {
	chatIdFlag := flag.String("chat", "", "chat id")
	photoFlag := flag.String("ph", "", "send photo --ph path/to/file")
	captionFlag := flag.String("cap", "", "caption to photo --cap text")
	flag.Parse()

	chatIdRaw := *chatIdFlag
	photoPath := *photoFlag
	caption := *captionFlag

	if chatIdRaw == "" && photoPath == "" && caption == "" {
		return false
	}

	if chatIdRaw == "" {
		log.Printf("chat must present")
		return true
	}

	if photoPath == "" {
		log.Printf("path to file must present")
		return true
	}

	chatId64, err := strconv.ParseInt(chatIdRaw, 10, 64)
	if err != nil {
		log.Printf("strconv.Atoi error: %s", err)
		return true
	}

	err = requests.SendPhoto(conn.Client, chatId64, photoPath, caption)
	if err != nil {
		log.Printf("error in sending: %s", err)
		return true
	}

	return true
}

func main() {
	if err := start(); err != nil {
		log.Fatalf("Error: %s", err)
	}
}
