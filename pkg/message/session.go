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
	config         config.Config
	Conn           *websocket.Conn
	PeerConnection *client.Client
	mux            sync.RWMutex
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

func NewSessionWithIdent(ident *Ident, conn *websocket.Conn, client *client.Client) *Session {
	return &Session{
		config:         ident.Config,
		Conn:           conn,
		PeerConnection: client,
	}
}
