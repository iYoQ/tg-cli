package main

import (
	"errors"
	"flag"
	"strconv"
	"tg-cli/connection"
	"tg-cli/requests"
)

func loadFlags() flags {
	chatIdFlag := flag.String("chat", "", "chat id")
	fileFlag := flag.String("f", "", "send any file(includes photos)")
	photoFlag := flag.String("p", "", "send photo")
	captionFlag := flag.String("cap", "", "caption to photo or file")
	flag.Parse()

	return flags{
		chatIdFlag:  chatIdFlag,
		fileFlag:    fileFlag,
		photoFlag:   photoFlag,
		captionFlag: captionFlag,
	}
}

func checkFlags(conn *connection.Connection, flags flags) (bool, error) {
	if *flags.chatIdFlag == "" && (*flags.fileFlag == "" || *flags.photoFlag == "") && *flags.captionFlag == "" {
		return false, nil
	}

	if *flags.chatIdFlag == "" {
		return true, errors.New("error: --chat is required when using --f or --p")
	}

	onlyOneProvided := (*flags.fileFlag == "") != (*flags.photoFlag == "")
	if !onlyOneProvided {
		return true, errors.New("error: exactly one of --f or --p must be provided, not both")
	}

	chatId64, err := strconv.ParseInt(*flags.chatIdFlag, 10, 64)
	if err != nil {
		return true, err
	}

	if *flags.fileFlag != "" {
		err = requests.SendFile(conn.Client, chatId64, *flags.fileFlag, *flags.captionFlag)
		if err != nil {
			return true, err
		}
	} else if *flags.photoFlag != "" {
		err = requests.SendPhoto(conn.Client, chatId64, *flags.photoFlag, *flags.captionFlag)
		if err != nil {
			return true, err
		}
	}

	return true, nil
}
