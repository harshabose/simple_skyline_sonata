package main

import (
	"context"
	"time"

	"github.com/asticode/go-astiav"
	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v4"

	"github.com/harshabose/simple_webrtc_comm/client/pkg"
	"github.com/harshabose/simple_webrtc_comm/datachannel/pkg"
	"github.com/harshabose/simple_webrtc_comm/mediasource/pkg"
	"github.com/harshabose/simple_webrtc_comm/transcode/pkg"

	"github.com/harshabose/simple_webrtc_comm/cmd/delivery"
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
				client.WithH264MediaEngine(delivery.DefaultVideoClockRate, client.PacketisationMode1, client.ProfileLevelBaseline31, delivery.DefaultSPSBase64, delivery.DefaultPPSBase64),
				// client.WithDefaultMediaEngine(),
				// client.WithVP8MediaEngine(delivery.DefaultVideoClockRate),
				client.WithBandwidthControlInterceptor(300_000, time.Second),
				client.WithTWCCHeaderExtensionSender(),
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
				client.WithFirebaseOfferSignal,
				client.WithMediaSources(),
				client.WithDataChannels(),
				client.WithBandwidthControl(),
			)
			if err != nil {
				panic(err)
			}

			if _, err := pc.CreateDataChannel("MAVLINK",
				data.WithRandomBindPort,
				// data.WithMAVP2P(os.Getenv("MAVP2P_EXE_PATH"), os.Getenv("MAVLINK_SERIAL")),
			); err != nil {
				panic(err)
			}

			if _, err := pc.CreateMediaSource("A8-MINI", true,
				mediasource.WithH264Track(delivery.DefaultVideoClockRate, mediasource.PacketisationMode1, mediasource.ProfileLevelBaseline31),
				// mediasource.WithVP8Track(delivery.DefaultVideoClockRate),
				mediasource.WithPriority(mediasource.Level5),
				mediasource.WithStream(
					mediasource.WithBufferSize(int(delivery.DefaultVideoFPS*3)),
					mediasource.WithDemuxer(
						"0",
						transcode.WithAvFoundationInputFormatOption,
						// "rtsp://192.168.144.25:8554/main.264",
						// transcode.WithRTSPInputOption,
						transcode.WithDemuxerBufferSize(int(delivery.DefaultVideoFPS)*3),
					),
					mediasource.WithDecoder(transcode.WithDecoderBufferSize(int(delivery.DefaultVideoFPS)*3)),
					mediasource.WithFilter(
						transcode.VideoFilters,
						transcode.WithFilterBufferSize(int(delivery.DefaultVideoFPS)*3),
						transcode.WithVideoScaleFilterContent(delivery.DefaultVideoWidth, delivery.DefaultVideoHeight),
						transcode.WithVideoPixelFormatFilterContent(delivery.DefaultPixelFormat),
						transcode.WithVideoFPSFilterContent(delivery.DefaultVideoFPS),
					),
					mediasource.WithEncoder(
						astiav.CodecIDH264,
						transcode.WithEncoderBufferSize(int(delivery.DefaultVideoFPS)*3),
						transcode.WithWebRTCOptimisedOptions,
					),
				),
			); err != nil {
				panic(err)
			}

			if err := pc.Connect("DELIVERY"); err != nil {
				panic(err)
			}

			drone.WaitUntilClosed()
		}()
	}

}
