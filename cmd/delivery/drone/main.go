//go:build cgo_enabled

package main

import (
	"context"
	"math"
	"time"

	"github.com/asticode/go-astiav"
	"github.com/harshabose/mediapipe"
	"github.com/harshabose/mediapipe/pkg/dioreader"
	"github.com/harshabose/mediapipe/pkg/diowriter"
	"github.com/harshabose/mediapipe/pkg/loopback"
	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media"

	"github.com/harshabose/simple_webrtc_comm/client/pkg"
	"github.com/harshabose/simple_webrtc_comm/client/pkg/mediasource"
	"github.com/harshabose/simple_webrtc_comm/client/pkg/transcode"
	"github.com/harshabose/simple_webrtc_comm/cmd/delivery"
	"github.com/harshabose/tools/buffer/pkg"
)

func main() {
	for {
		func() {
			time.Sleep(2 * time.Second)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			l, err := loopback.NewLoopBack(context.Background(), "127.0.0.1:14559")
			if err != nil {
				panic(err)
			}

			pPool := buffer.CreatePacketPool()
			fPool := buffer.CreateFramePool()

			transcoder, err := transcode.CreateTranscoder(
				transcode.WithGeneralDemuxer(ctx,
					"0",
					transcode.WithAvFoundationInputFormatOption,
					transcode.WithDemuxerBuffer(int(delivery.DefaultVideoFPS), pPool),
				),
				transcode.WithGeneralDecoder(ctx,
					transcode.WithDecoderBuffer(int(delivery.DefaultVideoFPS), fPool),
				),
				transcode.WithGeneralFilter(ctx,
					transcode.VideoFilters,
					transcode.WithFilterBuffer(int(delivery.DefaultVideoFPS), fPool),
					transcode.WithVideoScaleFilterContent(delivery.DefaultVideoWidth, delivery.DefaultVideoHeight),
					transcode.WithVideoPixelFormatFilterContent(delivery.DefaultPixelFormat),
					transcode.WithVideoFPSFilterContent(delivery.DefaultVideoFPS),
				),
				transcode.WithGeneralEncoder(
					ctx,
					astiav.CodecIDH264,
					transcode.WithCodecSettings(transcode.LowLatencyBitrateControlled),
					transcode.WithEncoderBufferSize(int(delivery.DefaultVideoFPS), pPool),
				),
				// transcode.WithMultiEncoderBitrateControl(ctx,
				// 	astiav.CodecIDH264,
				// 	transcode.NewMultiConfig(delivery.MinimumBitrate, delivery.MaximumBitrate, 10),
				// 	transcode.LowLatencyBitrateControlled,
				// 	int(delivery.DefaultVideoFPS), buffer.CreatePacketPool(),
				// ),
			)
			if err != nil {
				panic(err)
			}

			mediaEngine := &webrtc.MediaEngine{}
			registry := &interceptor.Registry{}
			settings := &webrtc.SettingEngine{}

			drone, err := client.CreateClient(
				ctx, cancel, mediaEngine, registry, settings,
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
				client.GetRTCConfiguration(),
				client.WithFirebaseOfferSignal,
				client.WithMediaSources(),
				client.WithDataChannels(),
				// client.WithBandwidthControl(),
			)
			if err != nil {
				panic(err)
			}

			datachannel, err := pc.CreateDataChannel("MAVLINK")
			if err != nil {
				panic(err)
			}

			track, err := pc.CreateMediaSource("A8-MINI",
				mediasource.WithH264Track(delivery.DefaultVideoClockRate, mediasource.PacketisationMode1, mediasource.ProfileLevelBaseline31),
				mediasource.WithPriority(mediasource.Level5),
			)
			if err != nil {
				panic(err)
			}

			// bwe, err := pc.GetBWEstimator()
			// if err != nil {
			// 	panic(err)
			// }
			//
			// if err := bwe.Subscribe("A8-MINI", track.GetPriority(), transcoder.OnUpdateBitrate()); err != nil {
			// 	panic(err)
			// }

			if err := pc.Connect("DELIVERY"); err != nil {
				panic(err)
			}

			time.Sleep(5 * time.Second)

			rl := mediapipe.NewIdentityAnyReader[[]byte](l)
			wl := mediapipe.NewIdentityAnyWriter[[]byte](l)

			ird, err := dioreader.NewDataChannel(datachannel.DataChannel(), math.MaxUint16)
			if err != nil {
				panic(err)
			}
			iwd, err := diowriter.NewDataChannel(datachannel.DataChannel(), math.MaxUint16)
			if err != nil {
				panic(err)
			}

			rd := mediapipe.NewIdentityAnyReader[[]byte](ird)
			wd := mediapipe.NewIdentityAnyWriter[[]byte](iwd)

			w := mediapipe.NewAnyWriter[media.Sample, *astiav.Packet](track, nil)
			r := mediapipe.NewAnyReader[media.Sample, *astiav.Packet](transcoder, func(packet *astiav.Packet) (media.Sample, error) {
				s := media.Sample{
					Data:               make([]byte, packet.Size()),
					Timestamp:          time.Now(),
					Duration:           time.Second / time.Duration(delivery.DefaultVideoFPS),
					PacketTimestamp:    uint32(packet.Pts()),
					PrevDroppedPackets: 0,
					Metadata:           nil,
					RTPHeaders:         nil,
				}
				copy(s.Data, packet.Data())
				transcoder.PutBack(packet)

				return s, nil
			})

			transcoder.Start()
			mediapipe.NewAnyPipe(ctx, rl, wd).Start()
			mediapipe.NewAnyPipe(ctx, rd, wl).Start()
			mediapipe.NewAnyPipe(ctx, r, w).Start()
			// bwe.Start()

			drone.WaitUntilClosed()
		}()
	}
}
