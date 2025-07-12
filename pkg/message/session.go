package message

// WARN: DO NOT COMPILE THIS FILE. THIS IS FOR FUTURE

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/coder/websocket"

	"github.com/harshabose/simple_webrtc_comm/client"
	"github.com/harshabose/simple_webrtc_comm/pkg/config"
)

type Session struct {
	config *config.Config
	Conn   *websocket.Conn
	Client *client.Client
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	mux    sync.RWMutex
}

func (s *Session) Write(ctx context.Context, msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshalling; err: %s", err.Error())
	}

	if err := s.Conn.Write(ctx, websocket.MessageText, data); err != nil {
		return fmt.Errorf("error while sending message; err: %s", err.Error())
	}

	return nil
}

func NewSessionWithIdent(ctx context.Context, ident *Ident, conn *websocket.Conn, client *client.Client) *Session {
	ctx2, cancel2 := context.WithCancel(ctx)
	return &Session{
		config: ident.Config,
		Conn:   conn,
		Client: client,
		ctx:    ctx2,
		cancel: cancel2,
	}
}
