//go:build cgo_enabled

package main

import (
	"context"
	"math"
	"time"

	"github.com/asticode/go-astiav"
	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media"

	"github.com/harshabose/mediapipe"
	"github.com/harshabose/mediapipe/pkg/consumers"
	"github.com/harshabose/mediapipe/pkg/duplexers"
	"github.com/harshabose/mediapipe/pkg/generators"
	"github.com/harshabose/tools/pkg/buffer"

	"github.com/harshabose/simple_webrtc_comm/client"
	"github.com/harshabose/simple_webrtc_comm/client/pkg/mediasource"
	"github.com/harshabose/simple_webrtc_comm/client/pkg/transcode"
	"github.com/harshabose/simple_webrtc_comm/cmd/fpv"
)

func main() {
	for {
		func() {
			time.Sleep(2 * time.Second)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			l, err := duplexers.NewLoopBack(context.Background(), "127.0.0.1:14559")
			if err != nil {
				panic(err)
			}

			pPool := buffer.CreatePacketPool()
			fPool := buffer.CreateFramePool()

			transcoder, err := transcode.CreateTranscoder(
				transcode.WithGeneralDemuxer(ctx,
					"0",
					transcode.WithAvFoundationInputFormatOption,
					transcode.WithDemuxerBuffer(int(fpv.DefaultVideoFPS), pPool),
				),
				transcode.WithGeneralDecoder(ctx,
					transcode.WithDecoderBuffer(int(fpv.DefaultVideoFPS), fPool),
				),
				transcode.WithGeneralFilter(ctx,
					transcode.VideoFilters,
					transcode.WithFilterBuffer(int(fpv.DefaultVideoFPS), fPool),
					transcode.WithVideoScaleFilterContent(fpv.DefaultVideoWidth, fpv.DefaultVideoHeight),
					transcode.WithVideoPixelFormatFilterContent(fpv.DefaultPixelFormat),
					transcode.WithVideoFPSFilterContent(fpv.DefaultVideoFPS),
				),
				transcode.WithGeneralEncoder(
					ctx,
					astiav.CodecIDH264,
					transcode.WithCodecSettings(transcode.LowLatencyBitrateControlled),
					transcode.WithEncoderBufferSize(int(fpv.DefaultVideoFPS), pPool),
				),
			)
			if err != nil {
				panic(err)
			}

			mediaEngine := &webrtc.MediaEngine{}
			registry := &interceptor.Registry{}
			settings := &webrtc.SettingEngine{}

			drone, err := client.NewClient(
				ctx, cancel, mediaEngine, registry, settings,
				client.WithH264MediaEngine(fpv.DefaultVideoClockRate, mediasource.PacketisationMode1, mediasource.ProfileLevelBaseline31, fpv.DefaultSPSBase64, fpv.DefaultPPSBase64),
				client.WithNACKInterceptor(client.NACKGeneratorLowLatency, client.NACKResponderLowLatency),
				client.WithRTCPReportsInterceptor(client.RTCPReportIntervalLowLatency),
				client.WithSimulcastExtensionHeaders(),
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
			)
			if err != nil {
				panic(err)
			}

			datachannel, err := pc.CreateDataChannel("MAVLINK")
			if err != nil {
				panic(err)
			}

			track, err := pc.CreateMediaSource("A8-MINI",
				mediasource.WithH264Track(fpv.DefaultVideoClockRate, mediasource.PacketisationMode1, mediasource.ProfileLevelBaseline31),
				mediasource.WithPriority(mediasource.Level5),
			)
			if err != nil {
				panic(err)
			}

			if err := pc.Connect("FPV"); err != nil {
				panic(err)
			}

			time.Sleep(10 * time.Second)

			rl := mediapipe.NewIdentityAnyReader[[]byte](l)
			wl := mediapipe.NewIdentityAnyWriter[[]byte](l)

			ird, err := generators.NewIODataChannel(datachannel.DataChannel(), math.MaxUint16)
			if err != nil {
				panic(err)
			}
			iwd, err := consumers.NewIODataChannel(datachannel.DataChannel(), math.MaxUint16)
			if err != nil {
				panic(err)
			}

			rd := mediapipe.NewIdentityAnyReader[[]byte](ird)
			wd := mediapipe.NewIdentityAnyWriter[[]byte](iwd)

			w := mediapipe.NewAnyWriter[media.Sample, *astiav.Packet](consumers.NewPionSampleConsumer(track), nil)
			r := mediapipe.NewAnyReader[media.Sample, *astiav.Packet](transcoder, func(packet *astiav.Packet) (media.Sample, error) {
				s := media.Sample{
					Data:               make([]byte, packet.Size()),
					Timestamp:          time.Now(),
					Duration:           time.Second / time.Duration(fpv.DefaultVideoFPS),
					PacketTimestamp:    uint32(packet.Pts()),
					PrevDroppedPackets: 0,
					Metadata:           nil,
					RTPHeaders:         nil,
				}
				copy(s.Data, packet.Data())
				pPool.Put(packet)

				return s, nil
			})

			transcoder.Start()
			mediapipe.NewAnyPipe(ctx, rl, wd).Start()
			mediapipe.NewAnyPipe(ctx, rd, wl).Start()
			mediapipe.NewAnyPipe(ctx, r, w).Start()

			drone.WaitUntilClosed()
		}()
	}
}
