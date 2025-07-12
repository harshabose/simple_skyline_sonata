package message

import (
	"github.com/harshabose/simple_webrtc_comm/pkg/config"
)

const IdentProtocol Protocol = "ident"

type Ident struct {
	BaseMessage
	Config *config.Config `json:"config"`
}

func (m *Ident) GetProtocol() Protocol {
	return IdentProtocol
}

func (m *Ident) Validate() error {
	return nil
}

func (m *Ident) Process(session *Session) error {
	session.mux.Lock()
	defer session.mux.Unlock()

	session.config = m.Config
	// CHECK CONFIG

	return nil
}
