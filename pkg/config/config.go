package config

import (
	"time"
)

// WARN: DO NOT COMPILE THIS FILE. THIS IS FOR FUTURE

// DataChannelConfig represents configuration for a data channel
type DataChannelConfig struct {
	Label    string `json:"label" yaml:"label"`
	Ordered  bool   `json:"ordered" yaml:"ordered"`
	Protocol string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
}

// MediaConfig represents media-related configuration
type MediaConfig struct {
	Audio []AudioConfig `json:"audio" yaml:"audio"`
	Video []VideoConfig `json:"video" yaml:"video"`
}

// AudioConfig represents audio stream configuration
type AudioConfig struct {
	Label      string `json:"label" yaml:"label"`
	Enabled    bool   `json:"enabled" yaml:"enabled"`
	Codec      string `json:"codecs,omitempty" yaml:"codecs,omitempty"`
	MinBitrate int    `json:"min_bitrate,omitempty" yaml:"min_bitrate,omitempty"` // kbps
	MaxBitrate int    `json:"max_bitrate,omitempty" yaml:"max_bitrate,omitempty"` // kbps
	SampleRate int    `json:"sample_rate,omitempty" yaml:"sample_rate,omitempty"` // Hz
}

// VideoConfig represents video stream configuration
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

// SecurityConfig represents security-related configuration
type SecurityConfig struct {
	RequireAuth bool   `json:"require_auth" yaml:"require_auth"`
	AuthToken   string `json:"auth_token,omitempty" yaml:"auth_token,omitempty"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level          string        `json:"level" yaml:"level"`                       // "debug", "info", "warn", "error"
	EnableStats    bool          `json:"enable_stats" yaml:"enable_stats"`         // Enable WebRTC stats collection
	StatsInterval  time.Duration `json:"stats_interval" yaml:"stats_interval"`     // How often to collect stats
	LogSignaling   bool          `json:"log_signaling" yaml:"log_signaling"`       // Log signaling messages
	LogDataChannel bool          `json:"log_data_channel" yaml:"log_data_channel"` // Log data channel messages
}

// Config represents the complete peer connection configuration
type Config struct {
	// Core WebRTC configuration
	DataChannels []DataChannelConfig `json:"data_channels" yaml:"data_channels"`

	// Media configuration
	Media MediaConfig `json:"media" yaml:"media"`

	// Security configuration
	Security SecurityConfig `json:"security" yaml:"security"`

	// Logging configuration
	Logging LoggingConfig `json:"logging" yaml:"logging"`

	// Application-specific settings
	UserID string `json:"user_id,omitempty" yaml:"user_id,omitempty"`
	RoomID string `json:"room_id,omitempty" yaml:"room_id,omitempty"`
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig() *Config {
	return &Config{
		DataChannels: []DataChannelConfig{
			{
				Label:   "data",
				Ordered: true,
			},
		},
		Media: MediaConfig{
			Audio: []AudioConfig{{
				Enabled:    false,
				Codec:      "opus",
				MinBitrate: 64,
				MaxBitrate: 128,
				SampleRate: 48000,
			}},
			Video: []VideoConfig{{
				Enabled:    false,
				Codec:      "h264",
				Width:      1280,
				Height:     720,
				Framerate:  30,
				MinBitrate: 1000,
				MaxBitrate: 2500,
			}},
		},
		Security: SecurityConfig{
			RequireAuth: false,
		},
		Logging: LoggingConfig{
			Level:          "info",
			EnableStats:    false,
			StatsInterval:  5 * time.Second,
			LogSignaling:   false,
			LogDataChannel: false,
		},
	}
}
