package main

import (
	"context"
	"log"
	"path/filepath"
	"strconv"
	"tg-cli/handlers"

	tdlib "github.com/zelenin/go-tdlib/client"
)

func Auth(apiIdRaw string, apiHash string, updatesChannel chan *tdlib.Message) (*tdlib.Client, error) {
	apiId64, err := strconv.ParseInt(apiIdRaw, 10, 32)
	if err != nil {
		log.Printf("strconv.Atoi error: %s", err)
		return nil, err
	}

	apiId := int32(apiId64)

	tdlibParameters := &tdlib.SetTdlibParametersRequest{
		UseTestDc:           false,
		DatabaseDirectory:   filepath.Join("./", "database"),
		FilesDirectory:      filepath.Join("./", "files"),
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		UseMessageDatabase:  true,
		UseSecretChats:      false,
		ApiId:               apiId,
		ApiHash:             apiHash,
		SystemLanguageCode:  "en",
		DeviceModel:         "Server",
		SystemVersion:       "1.0.0",
		ApplicationVersion:  "1.0.0",
	}

	authorizer := tdlib.ClientAuthorizer(tdlibParameters)
	go tdlib.CliInteractor(authorizer)

	_, err = tdlib.SetLogVerbosityLevel(&tdlib.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: 1,
	})
	if err != nil {
		log.Printf("SetLogVerbosityLevel error: %s", err)
		return nil, err
	}

	resHandCallback := func(result tdlib.Type) {
		go func() {
			switch update := result.(type) {
			case *tdlib.UpdateNewMessage:
				updatesChannel <- update.Message
			}
		}()
	}

	client, err := tdlib.NewClient(authorizer, tdlib.WithResultHandler(tdlib.NewCallbackResultHandler(resHandCallback)))
	if err != nil {
		log.Printf("NewClient error: %s", err)
		return nil, err
	}

	go handlers.HandleShutDown(client)

	versionOption, err := client.GetOption(&tdlib.GetOptionRequest{
		Name: "version",
	})
	if err != nil {
		log.Printf("GetOption error: %s", err)
		return client, err
	}

	commitOption, err := client.GetOption(&tdlib.GetOptionRequest{
		Name: "commit_hash",
	})
	if err != nil {
		log.Printf("GetOption error: %s", err)
		return client, err
	}

	log.Printf("TDLib version: %s (commit: %s)", versionOption.(*tdlib.OptionValueString).Value, commitOption.(*tdlib.OptionValueString).Value)

	if commitOption.(*tdlib.OptionValueString).Value != tdlib.TDLIB_VERSION {
		log.Printf("TDLib verson supported by the library (%s) is not the same as TDLIB version (%s)", tdlib.TDLIB_VERSION, commitOption.(*tdlib.OptionValueString).Value)
	}

	me, err := client.GetMe(context.Background())
	if err != nil {
		log.Printf("GetMe error: %s", err)
		return client, err
	}

	log.Printf("Me: %s %s", me.FirstName, me.LastName)

	return client, nil
}
