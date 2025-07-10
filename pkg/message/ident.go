package message

// WARN: DO NOT COMPILE THIS FILE. THIS IS FOR FUTURE

import (
	"github.com/harshabose/simple_webrtc_comm/pkg/config"
)

const IdentProtocol Protocol = "ident"

type Ident struct {
	BaseMessage
	Config config.Config `json:"config"`
}

func (m *Ident) GetProtocol() Protocol {
	return IdentProtocol
}

func (m *Ident) Validate() error {
	return nil
}

func (m *Ident) signValidate() error {
	// serverSignatire := os.Getenv("SERVER_SIGNATURE")

	return nil
}

func (m *Ident) Process(session *Session) error {
	session.mux.Lock()
	defer session.mux.Unlock()

	if err := m.signValidate(); err != nil {
		return err
	}

	session.config = m.Config
	// CHECK CONFIG

	return nil
}
