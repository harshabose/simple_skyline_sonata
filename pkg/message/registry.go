package message

// WARN: DO NOT COMPILE THIS FILE. THIS IS FOR FUTURE

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Factory func() Message

type RegistryBuilder struct {
	registry map[Protocol]Factory
}

func NewRegistryBuilder() *RegistryBuilder {
	return &RegistryBuilder{
		registry: make(map[Protocol]Factory),
	}
}

func (b *RegistryBuilder) Register(protocol Protocol, factory Factory) error {
	if _, exists := b.registry[protocol]; exists {
		return errors.New("protocol already registered")
	}
	b.registry[protocol] = factory
	return nil
}

func (b *RegistryBuilder) Build() *Registry {
	registry := make(map[Protocol]Factory)
	for k, v := range b.registry {
		registry[k] = v
	}

	return &Registry{
		registry: registry,
	}
}

type Registry struct {
	registry map[Protocol]Factory
}

func (r *Registry) UnMarshal(data []byte) (Message, error) {
	if data == nil {
		return nil, errors.New("data cannot be nil")
	}

	base := &BaseMessage{}

	if err := json.Unmarshal(data, base); err != nil {
		return nil, fmt.Errorf("error while unmarshalling; err: %s", err.Error())
	}

	return r.UnMarshalWithProtocol(base.Protocol, data)
}

func (r *Registry) UnMarshalWithProtocol(protocol Protocol, data []byte) (Message, error) {
	if data == nil {
		return nil, errors.New("data cannot be nil")
	}

	if protocol == NoneProtocol || protocol == "" {
		return nil, errors.New("protocol is none")
	}

	factory, exists := r.registry[protocol]
	if !exists {
		return nil, errors.New("protocol not registered")
	}

	msg := factory()

	if err := json.Unmarshal(data, msg); err != nil {
		return nil, fmt.Errorf("error while unmarshalling; err: %s", err.Error())
	}

	return msg, nil
}
