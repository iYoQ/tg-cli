package handlers

import (
	"context"
	"github.com/zelenin/go-tdlib/client"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func HandleShutDown(my_client *client.Client) {
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP)
	<-ch
	log.Println("\nShutting down TDLib client...")

	ok, err := my_client.Close(context.Background())

	if err != nil {
		log.Printf("Error closing TDLib client: %v\n", err)
		os.Exit(1)
	}

	if ok != nil {
		log.Println("TDLib client closed")
	}

	os.Exit(0)
}
