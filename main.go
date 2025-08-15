package main

import (
	"context"
	"os"
	"runtime/debug"
	"strconv"
	"tg-cli/app"
	"tg-cli/connection"
	"tg-cli/logger"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	apiId          int32
	apiHash        string
	verbosityLevel int32
}

type flags struct {
	chatIdFlag  *string
	fileFlag    *string
	photoFlag   *string
	captionFlag *string
}

func start(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)

	cfg := loadParams(ctx)
	flags := loadFlags()

	conn := connection.NewConnection()

	defer func() {
		if r := recover(); r != nil {
			var errMsg string
			switch v := r.(type) {
			case string:
				errMsg = v
			case error:
				errMsg = v.Error()
			default:
				errMsg = "Unknown panic type"
			}

			logger.Error().Str("panic", errMsg).Str("stack", string(debug.Stack())).Msg("panic recovered")
		}

		if conn.Client != nil {
			conn.Close(ctx, cfg.verbosityLevel)
		}
	}()

	if err := auth(cfg, conn, ctx); err != nil {
		return err
	}

	if ok, err := checkFlags(conn, flags); err != nil {
		return err
	} else if ok {
		return nil
	}

	app := tea.NewProgram(app.NewRootModel(conn, ctx), tea.WithAltScreen())
	if _, err := app.Run(); err != nil {
		return err
	}

	return nil
}

func loadParams(ctx context.Context) Config {
	logger := zerolog.Ctx(ctx)

	godotenv.Load()
	apiIdRaw := os.Getenv("API_ID")
	apiHash := os.Getenv("API_HASH")
	devMode := os.Getenv("DEV_ENV")

	if apiIdRaw == "" || apiHash == "" {
		logger.Fatal().Msg("API_ID and API_HASH are required, use .env file, or ENV")
	}

	apiId64, err := strconv.ParseInt(apiIdRaw, 10, 64)
	if err != nil {
		logger.Fatal().Err(err).Msg("from loadParams")
	}

	apiId := int32(apiId64)

	verbosityLevel := int32(0)
	if devMode == "true" {
		verbosityLevel = 1
	}

	return Config{
		apiId:          apiId,
		apiHash:        apiHash,
		verbosityLevel: verbosityLevel,
	}
}

func main() {
	logFile, err := os.OpenFile("tg-cli.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open log file")
	}
	defer logFile.Close()

	ctx := logger.NewLogger(logFile)
	logger := zerolog.Ctx(ctx)

	if err := start(ctx); err != nil {
		logger.Fatal().Err(err).Msg("from main")
	}
}
