package connection

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	tdlib "github.com/zelenin/go-tdlib/client"
)

type Me struct {
	Id        int64
	FirstName string
	LastName  string
}

type Connection struct {
	Client         *tdlib.Client
	UpdatesChannel chan *tdlib.Message
	me             Me
}

func NewConnection() *Connection {
	updatesChannel := make(chan *tdlib.Message)
	return &Connection{
		Client:         nil,
		UpdatesChannel: updatesChannel,
	}
}

func (conn *Connection) SetClient(client *tdlib.Client) {
	conn.Client = client
}

func (conn *Connection) SetMe(tdlibMe *tdlib.User) Me {
	me := Me{
		Id:        tdlibMe.Id,
		FirstName: tdlibMe.FirstName,
		LastName:  tdlibMe.LastName,
	}
	conn.me = me
	return me
}

func (conn Connection) GetMe() Me {
	return conn.me
}

func (conn *Connection) CreateCallbackHandler(result tdlib.Type) {
	go func() {
		switch update := result.(type) {
		case *tdlib.UpdateNewMessage:
			if conn.UpdatesChannel != nil {
				conn.UpdatesChannel <- update.Message
			} else {
				log.Println("channel don't setup, check connection.UpdateChannel")
			}
		}

	}()
}

func (conn *Connection) Close() {
	if conn.Client == nil {
		return
	}

	log.Println("\nShutting down TDLib client...")

	ok, err := conn.Client.Close(context.Background())
	if err != nil {
		log.Printf("Error closing TDLib client: %v\n", err)
		os.Exit(1)
	}
	if ok != nil {
		log.Println("TDLib client closed")
		return
	}

	panic(errors.New("smh very bad happened"))
}

func (conn *Connection) ShutDownListener() {
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP)
	<-ch

	conn.Close()
	os.Exit(0)
}
