package main

import (
	"context"
	"os"
	"time"

	"github.com/asticode/go-astiav"
	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v4"

	"github.com/harshabose/simple_webrtc_comm/client/pkg"
	"github.com/harshabose/simple_webrtc_comm/cmd/audio"
	"github.com/harshabose/simple_webrtc_comm/datachannel/pkg"
	"github.com/harshabose/simple_webrtc_comm/mediasource/pkg"
	"github.com/harshabose/simple_webrtc_comm/transcode/pkg"
)

func main() {
	for {
		func() {
			time.Sleep(2 * time.Second)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			mediaEngine := &webrtc.MediaEngine{}
			registry := &interceptor.Registry{}

			drone, err := client.CreateClient(
				ctx, cancel, mediaEngine, registry,
				client.WithOpusMediaEngine(audio.DefaultAudioSampleRate, audio.DefaultAudioChannelLayout, audio.DefaultAudioStereo),
				client.WithNACKInterceptor(client.NACKGeneratorLowLatency, client.NACKResponderLowLatency),
				client.WithRTCPReportsInterceptor(client.RTCPReportIntervalLowLatency),
				client.WithSimulcastExtensionHeaders(),
				client.WithTWCCSenderInterceptor(client.TWCCIntervalLowLatency),
			)
			if err != nil {
				panic(err)
			}

			pc, err := drone.CreatePeerConnection(
				"MAIN",
				client.WithRTCConfiguration(client.GetRTCConfiguration()),
				client.WithOfferSignal,
				client.WithMediaSources(),
				client.WithMediaSinks(),
				client.WithDataChannels(),
			)
			if err != nil {
				panic(err)
			}

			if _, err := pc.CreateDataChannel("MAVLINK",
				data.WithRandomBindPort,
				data.WithMAVP2P(os.Getenv("MAVP2P_EXE_PATH"), os.Getenv("MAVLINK_SERIAL")),
			); err != nil {
				panic(err)
			}

			notchUpdater, err := transcode.CreatePropNoiseFilterUpdator(ctx, "ttyTHS0", 57600, 200*time.Millisecond)
			notchUpdater.AddNotchFilter(transcode.PropellerOne, 500, 5, 2)
			notchUpdater.AddNotchFilter(transcode.PropellerTwo, 500, 5, 2)
			notchUpdater.AddNotchFilter(transcode.PropellerThree, 500, 5, 2)
			notchUpdater.AddNotchFilter(transcode.PropellerFour, 500, 5, 2)
			notchUpdater.AddNotchFilter(transcode.PropellerFive, 500, 5, 2)
			notchUpdater.AddNotchFilter(transcode.PropellerSix, 500, 5, 2)
			if err != nil {
				panic(err)
			}

			if err := pc.CreateMediaSource("RODE", false,
				mediasource.WithOpusTrack(audio.DefaultAudioSampleRate, audio.DefaultAudioChannelLayout, mediasource.StereoDual),
				mediasource.WithPriority(mediasource.Level5),
				mediasource.WithStream(
					mediasource.WithBufferSize(int(audio.DefaultAudioSamplesPerFrame)),
					mediasource.WithDemuxer(
						"default",
						transcode.WithAlsaInputFormatOption,
						transcode.WithDemuxerBufferSize(int(audio.DefaultAudioSamplesPerFrame)),
					),
					mediasource.WithDecoder(transcode.WithDecoderBufferSize(int(audio.DefaultAudioSamplesPerFrame))),
					mediasource.WithFilter(
						transcode.AudioFilters,
						transcode.WithFilterBufferSize(int(audio.DefaultAudioSamplesPerFrame)),
						transcode.WithAudioSampleFormatChannelLayoutFilter(audio.DefaultAudioSampleFormat, astiav.ChannelLayoutStereo),
						transcode.WithAudioSampleRateFilter(audio.DefaultAudioSampleRate),
						transcode.WithAudioSamplesPerFrameContent(audio.DefaultAudioSamplesPerFrame),
						transcode.WithAudioLowPassFilterContent("4000_cut", 4000, 2),
						transcode.WithAudioHighPassFilterContent("120_cut", 120, 1),
						transcode.WithAudioNotchHarmonicsFilterContent(transcode.PropellerOne.String(), 500, 5, 10.0),
						transcode.WithAudioNotchHarmonicsFilterContent(transcode.PropellerTwo.String(), 500, 5, 10.0),
						transcode.WithAudioNotchHarmonicsFilterContent(transcode.PropellerThree.String(), 500, 5, 10.0),
						transcode.WithAudioNotchHarmonicsFilterContent(transcode.PropellerFour.String(), 500, 5, 10.0),
						transcode.WithAudioNotchHarmonicsFilterContent(transcode.PropellerFive.String(), 500, 5, 10.0),
						transcode.WithAudioNotchHarmonicsFilterContent(transcode.PropellerSix.String(), 500, 5, 10.0),
						transcode.WithMeanBroadBandNoiseFilter("mean_cut", 1.0, 0.010, 0.050),
						transcode.WithAudioEqualiserFilter("mid_eq", 2500, 1000, +5),
						transcode.WithAudioEqualiserFilter("low_eq", 250, 100, +3),
						transcode.WithUpdateFilter(notchUpdater),
					),
					mediasource.WithEncoder(
						astiav.CodecIDOpus,
						transcode.WithEncoderBufferSize(int(audio.DefaultAudioSamplesPerFrame)),
					),
				),
			); err != nil {
				panic(err)
			}

			if err := pc.Connect("AUDIO"); err != nil {
				panic(err)
			}

			drone.WaitUntilClosed()
		}()
	}
}
