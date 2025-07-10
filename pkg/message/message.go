package message

// WARN: DO NOT COMPILE THIS FILE. THIS IS FOR FUTURE

import "errors"

type Protocol string

const NoneProtocol Protocol = "none-protocol"

type Message interface {
	Validate() error
	GetProtocol() Protocol
	Process(*Session) error
}

type BaseMessage struct {
	Protocol Protocol `json:"protocol"`
}

func (m *BaseMessage) Validate() error {
	return errors.New("not valid")
}

func (m *BaseMessage) GetProtocol() Protocol {
	return NoneProtocol
}

func (m *BaseMessage) Process(*Session) error {
	return errors.New("process not implemented")
}
