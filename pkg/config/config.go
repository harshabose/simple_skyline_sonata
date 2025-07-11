package config

import (
	"time"
)

// WARN: DO NOT COMPILE THIS FILE. THIS IS FOR FUTURE

type DataChannelConfig struct {
	Label             string `json:"label" yaml:"label"`
	ID                uint16 `json:"id" yaml:"id"`
	Ordered           bool   `json:"ordered,omitempty" yaml:"ordered,omitempty"`
	Protocol          string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	MaxPacketLifeTime uint16 `json:"max_packet_life_time" yaml:"max_packet_life_time"`
	MaxRetransmits    uint16 `json:"max_retransmits,omitempty" yaml:"max_retransmits,omitempty"`
}

type MediaConfig struct {
	Audio []AudioConfig `json:"audio" yaml:"audio"`
	Video []VideoConfig `json:"video" yaml:"video"`
}

type AudioConfig struct {
	Label      string `json:"label" yaml:"label"`
	Enabled    bool   `json:"enabled" yaml:"enabled"`
	Codec      string `json:"codecs,omitempty" yaml:"codecs,omitempty"`
	MinBitrate int    `json:"min_bitrate,omitempty" yaml:"min_bitrate,omitempty"` // kbps
	MaxBitrate int    `json:"max_bitrate,omitempty" yaml:"max_bitrate,omitempty"` // kbps
	SampleRate int    `json:"sample_rate,omitempty" yaml:"sample_rate,omitempty"` // Hz
}

type VideoConfig struct {
	Label      string `json:"label" yaml:"label"`
	Enabled    bool   `json:"enabled" yaml:"enabled"`
	Codec      string `json:"codecs,omitempty" yaml:"codecs,omitempty"`
	MinBitrate int    `json:"min_bitrate,omitempty" yaml:"min_bitrate,omitempty"` // kbps
	MaxBitrate int    `json:"max_bitrate,omitempty" yaml:"max_bitrate,omitempty"` // kbps
	Width      int    `json:"width,omitempty" yaml:"width,omitempty"`
	Height     int    `json:"height,omitempty" yaml:"height,omitempty"`
	Framerate  int    `json:"framerate,omitempty" yaml:"framerate,omitempty"`
}

// type SecurityConfig struct {
// 	RequireAuth bool   `json:"require_auth" yaml:"require_auth"`
// 	AuthToken   string `json:"auth_token,omitempty" yaml:"auth_token,omitempty"`
// }

type LoggingConfig struct {
	Level          string        `json:"level" yaml:"level"`                       // "debug", "info", "warn", "error"
	EnableStats    bool          `json:"enable_stats" yaml:"enable_stats"`         // Enable WebRTC stats collection
	StatsInterval  time.Duration `json:"stats_interval" yaml:"stats_interval"`     // How often to collect stats
	LogSignaling   bool          `json:"log_signaling" yaml:"log_signaling"`       // Log signaling messages
	LogDataChannel bool          `json:"log_data_channel" yaml:"log_data_channel"` // Log data channel messages
}

type Config struct {
	DataChannels []DataChannelConfig `json:"data_channels" yaml:"data_channels"`
	Media        MediaConfig         `json:"media" yaml:"media"`

	// Security configuration
	// Security SecurityConfig `json:"security" yaml:"security"`

	Logging LoggingConfig `json:"logging" yaml:"logging"`

	UserID string `json:"user_id,omitempty" yaml:"user_id,omitempty"`
	RoomID string `json:"room_id,omitempty" yaml:"room_id,omitempty"`
}
