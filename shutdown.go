package main

import (
	"context"
	"github.com/zelenin/go-tdlib/client"
	"os"
	"os/signal"
	"syscall"
)

func HandleShutDown(my_client client.Client) {
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	my_client.Close(context.Background())
	os.Exit(1)
}
