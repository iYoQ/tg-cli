package main

import (
	"context"
	"log"
	"path/filepath"
	"tg-cli/connection"

	tdlib "github.com/zelenin/go-tdlib/client"
)

func auth(cfg Config, conn *connection.Connection) error {
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
		NewVerbosityLevel: 1,
	})
	if err != nil {
		log.Printf("SetLogVerbosityLevel error: %s", err)
		return err
	}

	client, err := tdlib.NewClient(authorizer, tdlib.WithResultHandler(tdlib.NewCallbackResultHandler(conn.CreateCallbackHandler)))
	if err != nil {
		log.Printf("NewClient error: %s", err)
		return err
	}

	conn.SetClient(client)

	go conn.ShutDownListener()

	versionOption, err := client.GetOption(&tdlib.GetOptionRequest{
		Name: "version",
	})
	if err != nil {
		log.Printf("GetOption error: %s", err)
		return err
	}

	commitOption, err := client.GetOption(&tdlib.GetOptionRequest{
		Name: "commit_hash",
	})
	if err != nil {
		log.Printf("GetOption error: %s", err)
		return err
	}

	log.Printf("TDLib version: %s (commit: %s)", versionOption.(*tdlib.OptionValueString).Value, commitOption.(*tdlib.OptionValueString).Value)

	if commitOption.(*tdlib.OptionValueString).Value != tdlib.TDLIB_VERSION {
		log.Printf("TDLib verson supported by the library (%s) is not the same as TDLIB version (%s)", tdlib.TDLIB_VERSION, commitOption.(*tdlib.OptionValueString).Value)
	}

	tdlibMe, err := client.GetMe(context.Background())
	if err != nil {
		log.Printf("GetMe error: %s", err)
		return err
	}

	me := conn.SetMe(tdlibMe)

	log.Printf("Me: %s %s", me.FirstName, me.LastName)

	return nil
}
