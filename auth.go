package main

import (
	"context"
	"path/filepath"
	"tg-cli/connection"

	"github.com/rs/zerolog"
	tdlib "github.com/zelenin/go-tdlib/client"
)

func auth(cfg Config, conn *connection.Connection, ctx context.Context) error {
	logger := zerolog.Ctx(ctx)

	tdlibParameters := &tdlib.SetTdlibParametersRequest{
		UseTestDc:           false,
		DatabaseDirectory:   filepath.Join("./", "database"),
		FilesDirectory:      filepath.Join("./", "files"),
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		UseMessageDatabase:  true,
		UseSecretChats:      false,
		ApiId:               cfg.apiId,
		ApiHash:             cfg.apiHash,
		SystemLanguageCode:  "en",
		DeviceModel:         "Server",
		SystemVersion:       "1.0.0",
		ApplicationVersion:  "1.0.0",
	}

	authorizer := tdlib.ClientAuthorizer(tdlibParameters)
	go tdlib.CliInteractor(authorizer)

	_, err := tdlib.SetLogVerbosityLevel(&tdlib.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: cfg.verbosityLevel,
	})
	if err != nil {
		return err
	}

	client, err := tdlib.NewClient(authorizer, tdlib.WithResultHandler(tdlib.NewCallbackResultHandler(conn.CreateCallbackHandler)))
	if err != nil {
		return err
	}

	conn.SetClient(client)

	go conn.ShutDownListener(ctx, cfg.verbosityLevel)

	versionOption, err := client.GetOption(&tdlib.GetOptionRequest{
		Name: "version",
	})
	if err != nil {
		return err
	}

	commitOption, err := client.GetOption(&tdlib.GetOptionRequest{
		Name: "commit_hash",
	})
	if err != nil {
		return err
	}

	if cfg.verbosityLevel > 0 {
		logger.Info().Str("TDLib version", versionOption.(*tdlib.OptionValueString).Value).Str("commit", commitOption.(*tdlib.OptionValueString).Value).Msg("")

		if commitOption.(*tdlib.OptionValueString).Value != tdlib.TDLIB_VERSION {
			logger.Warn().Str("TDLib supported version", tdlib.TDLIB_VERSION).Str("your version", commitOption.(*tdlib.OptionValueString).Value).Msg("")
		}
	}

	tdlibMe, err := client.GetMe(context.Background())
	if err != nil {
		return err
	}

	me := conn.SetMe(tdlibMe)

	if cfg.verbosityLevel > 0 {
		logger.Info().Str("FirstName", me.FirstName).Str("LastName", me.LastName).Msg("me")
	}

	return nil
}
