package connection

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
				log.Fatal().Msg("channel don't setup")
			}
		}

	}()
}

func (conn *Connection) Close(ctx context.Context, verbosityLevel int32) {
	if conn.Client == nil {
		return
	}

	logger := zerolog.Ctx(ctx)

	if verbosityLevel > 0 {
		logger.Info().Msg("Shutting down TDLib client...")
	}

	ok, err := conn.Client.Close(context.Background())
	if err != nil {
		logger.Fatal().Err(err).Msg("Error closing TDLib client")
		os.Exit(1)
	}
	if ok != nil {
		if verbosityLevel > 0 {
			logger.Info().Msg("TDLib client closed")
		}
		return
	}

	panic(errors.New("smh very bad happened"))
}

func (conn *Connection) ShutDownListener(ctx context.Context, verbosityLevel int32) {
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP)
	<-ch

	conn.Close(ctx, verbosityLevel)
	os.Exit(0)
}
