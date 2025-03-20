package main

import (
	"context"
	"time"

	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v4"

	"github.com/harshabose/simple_webrtc_comm/client/pkg"
	"github.com/harshabose/simple_webrtc_comm/cmd/audio"
	"github.com/harshabose/simple_webrtc_comm/datachannel/pkg"
)

func main() {
	for {
		func() {
			time.Sleep(2 * time.Second)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			mediaEngine := &webrtc.MediaEngine{}
			registry := &interceptor.Registry{}

			gcs, err := client.CreateClient(
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

			pc, err := gcs.CreatePeerConnection(
				"MAIN",
				client.WithRTCConfiguration(client.GetRTCConfiguration()),
				client.WithAnswerSignal,
				client.WithMediaSources(),
				client.WithMediaSinks(),
				client.WithDataChannels(),
			)

			if _, err := pc.CreateDataChannel("MAVLINK", data.WithRandomBindPort, data.WithLoopBackPort(14550)); err != nil {
				panic(err)
			}

			// if err := pc.CreateMediaSource("MAC", false,
			// 	mediasource.WithOpusTrack(audio.DefaultAudioSampleRate, audio.DefaultAudioChannelLayout, mediasource.StereoDual),
			// 	mediasource.WithPriority(mediasource.Level5),
			// 	mediasource.WithStream(
			// 		mediasource.WithBufferSize(int(audio.DefaultAudioSamplesPerFrame)),
			// 		mediasource.WithDemuxer(
			// 			"default",
			// 			transcode.WithAvFoundationInputFormatOption,
			// 			transcode.WithDemuxerBufferSize(int(audio.DefaultAudioSamplesPerFrame)),
			// 		),
			// 		mediasource.WithDecoder(transcode.WithDecoderBufferSize(int(audio.DefaultAudioSamplesPerFrame))),
			// 		mediasource.WithFilter(
			// 			transcode.VideoFilters,
			// 			transcode.WithFilterBufferSize(int(audio.DefaultAudioSamplesPerFrame)),
			// 			transcode.WithAudioSampleFormatChannelLayoutFilter(audio.DefaultAudioSampleFormat, astiav.ChannelLayoutStereo),
			// 			transcode.WithAudioSampleRateFilter(audio.DefaultAudioSampleRate),
			// 			transcode.WithAudioSamplesPerFrameContent(audio.DefaultAudioSamplesPerFrame),
			// 		),
			// 		mediasource.WithEncoder(
			// 			astiav.CodecIDOpus,
			// 			transcode.WithEncoderBufferSize(int(audio.DefaultAudioSamplesPerFrame)),
			// 		),
			// 	),
			// ); err != nil {
			// 	panic(err)
			// }

			if err := pc.Connect("AUDIO"); err != nil {
				panic(err)
			}

			gcs.WaitUntilClosed()
		}()
	}
}
