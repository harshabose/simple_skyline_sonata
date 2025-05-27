package main

import (
	"context"
	"time"

	"github.com/asticode/go-astiav"
	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v4"

	"github.com/harshabose/simple_webrtc_comm/client/pkg"
	"github.com/harshabose/simple_webrtc_comm/cmd/delivery"
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
				client.WithH264MediaEngine(delivery.DefaultVideoClockRate, client.PacketisationMode1, client.ProfileLevelBaseline31, delivery.DefaultSPSBase64, delivery.DefaultPPSBase64),
				client.WithBandwidthControlInterceptor(delivery.InitialBitrate, delivery.MinimumBitrate, delivery.MaximumBitrate, time.Second),
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

			// if _, err := pc.CreateDataChannel("MAVLINK",
			// 	data.WithRandomBindPort,
			// 	// data.WithMAVP2P(os.Getenv("MAVP2P_EXE_PATH"), os.Getenv("MAVLINK_SERIAL")),
			// ); err != nil {
			// 	panic(err)
			// }

			if _, err := pc.CreateMediaSource("A8-MINI", true,
				mediasource.WithH264Track(delivery.DefaultVideoClockRate, mediasource.PacketisationMode1, mediasource.ProfileLevelBaseline31),
				mediasource.WithPriority(mediasource.Level5),
				mediasource.WithStream(
					mediasource.WithBufferSize(int(delivery.DefaultVideoFPS*3)),
					mediasource.WithTranscoder(
						transcode.WithGeneralDemuxer(ctx,
							"/dev/video0",
							// transcode.WithAvFoundationInputFormatOption,
							transcode.WithDemuxerBufferSize(int(delivery.DefaultVideoFPS)),
						),
						transcode.WithGeneralDecoder(ctx,
							transcode.WithDecoderBufferSize(int(delivery.DefaultVideoFPS)),
						),
						transcode.WithGeneralFilter(ctx,
							transcode.VideoFilters,
							transcode.WithFilterBufferSize(int(delivery.DefaultVideoFPS)),
							transcode.WithVideoScaleFilterContent(delivery.DefaultVideoWidth, delivery.DefaultVideoHeight),
							transcode.WithVideoPixelFormatFilterContent(delivery.DefaultPixelFormat),
							transcode.WithVideoFPSFilterContent(delivery.DefaultVideoFPS),
						),
						transcode.WithBitrateControlEncoder(ctx,
							astiav.CodecIDH264,
							transcode.UpdateConfig{
								MinBitrate:              delivery.MinimumBitrate,              // 500kbps
								MaxBitrate:              delivery.MaximumBitrate,              // 2Mbps
								CutVideoBelowMinBitrate: delivery.CutVideoBelowMinimumBitrate, // Enable pausing
							},
							transcode.LowLatencyBitrateControlled,
							int(delivery.DefaultVideoFPS),
						),
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
