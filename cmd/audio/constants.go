package audio

import (
	"github.com/asticode/go-astiav"

	client "github.com/harshabose/simple_webrtc_comm/client/pkg"
)

const (
	DefaultAudioSampleRate      uint32 = 48000
	DefaultAudioChannelLayout   uint16 = 2
	DefaultAudioStereo                 = client.Dual
	DefaultAudioSampleFormat           = astiav.SampleFormatFlt
	DefaultAudioSamplesPerFrame uint16 = 960
)
