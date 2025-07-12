package config

import (
	"context"
	"time"

	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v4"

	"github.com/harshabose/simple_webrtc_comm/client"
)

type DataChannelConfig struct {
	Label             string `json:"label" yaml:"label"`
	ID                uint16 `json:"id" yaml:"id"`
	Ordered           bool   `json:"ordered,omitempty" yaml:"ordered,omitempty"`
	Protocol          string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	MaxPacketLifeTime uint16 `json:"max_packet_life_time" yaml:"max_packet_life_time"`
	MaxRetransmits    uint16 `json:"max_retransmits,omitempty" yaml:"max_retransmits,omitempty"`
}

type MediaConfig struct {
	Audio []AudioConfig `json:"audio,omitempty" yaml:"audio,omitempty"`
	Video []VideoConfig `json:"video,omitempty" yaml:"video,omitempty"`
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
	ClientConfig          client.ClientConfig           `json:"client_config" yaml:"client_config"`
	PeerConnectionConfigs []client.PeerConnectionConfig `json:"peer_connection_configs" yaml:"peer_connection_configs"`

	MediaEngine         webrtc.MediaEngine   `json:"-"`
	InterceptorRegistry interceptor.Registry `json:"-"`
	SettingsEngine      webrtc.SettingEngine `json:"-"`

	// Security SecurityConfig `json:"security" yaml:"security"`

	Logging LoggingConfig `json:"logging" yaml:"logging"`

	UserID string `json:"user_id,omitempty" yaml:"user_id,omitempty"`
	RoomID string `json:"room_id,omitempty" yaml:"room_id,omitempty"`
}

func (c *Config) GenerateClient(ctx context.Context, cancel context.CancelFunc) (*client.Client, error) {
	return client.NewClientFromConfig(ctx, cancel, &c.MediaEngine, &c.InterceptorRegistry, &c.SettingsEngine, &c.ClientConfig)
}

func (c *Config) GeneratePeerConnection(client *client.Client, pcc client.PeerConnectionConfig) (*client.PeerConnection, error) {
	pc, err := client.CreatePeerConnection(pcc.Name, pcc.RTCConfig, pcc.ToOptions()...)
	if err != nil {
		return nil, err
	}

	if err := pcc.CreateDataChannels(pc); err != nil {
		// TODO: REMOVE PEER CONNECTION
		return nil, err
	}

	if err := pcc.CreateMediaSources(pc); err != nil {
		// TODO: REMOVE PEER CONNECTION
		return nil, err
	}

	if err := pcc.CreateMediaSinks(pc); err != nil {
		// TODO: REMOVE PEER CONNECTION
		return nil, err
	}

	return pc, nil

}
