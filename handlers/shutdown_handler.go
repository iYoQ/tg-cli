package handlers

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	tdlib "github.com/zelenin/go-tdlib/client"
)

func HandleShutDown(client *tdlib.Client) {
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP)
	<-ch

	ShutDown(client)
	os.Exit(0)
}

func ShutDown(client *tdlib.Client) *tdlib.Ok {
	log.Println("\nShutting down TDLib client...")

	ok, err := client.Close(context.Background())
	if err != nil {
		log.Printf("Error closing TDLib client: %v\n", err)
		os.Exit(1)
	}
	if ok != nil {
		log.Println("TDLib client closed")
		return ok
	}

	panic(errors.New("smh very bad happened"))
}
