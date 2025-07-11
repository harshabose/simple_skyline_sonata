package main

import (
	"context"
	"math"
	"time"

	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v4"

	"github.com/harshabose/mediapipe"
	"github.com/harshabose/mediapipe/pkg/consumers"
	"github.com/harshabose/mediapipe/pkg/duplexers"
	"github.com/harshabose/mediapipe/pkg/generators"
	"github.com/harshabose/simple_webrtc_comm/client"
	"github.com/harshabose/simple_webrtc_comm/client/pkg/mediasink"
	"github.com/harshabose/simple_webrtc_comm/cmd/fpv"
)

func main() {
	for {
		func() {
			time.Sleep(2 * time.Second)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			l, err := duplexers.NewLoopBack(context.Background(), "127.0.0.1:0", duplexers.WithLoopBackPort(14550))
			if err != nil {
				panic(err)
			}

			mediaEngine := &webrtc.MediaEngine{}
			registry := &interceptor.Registry{}
			settings := &webrtc.SettingEngine{}

			gcs, err := client.CreateClient(
				ctx, cancel, mediaEngine, registry, settings,
				client.WithH264MediaEngine(fpv.DefaultVideoClockRate, client.PacketisationMode1, client.ProfileLevelBaseline31, fpv.DefaultSPSBase64, fpv.DefaultPPSBase64),
				client.WithTWCCHeaderExtensionSender(),
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
				client.GetRTCConfiguration(),
				client.WithFirebaseAnswerSignal,
				client.WithMediaSinks(),
				client.WithDataChannels(),
			)
			if err != nil {
				panic(err)
			}

			datachannel, err := pc.CreateDataChannel("MAVLINK")
			if err != nil {
				panic(err)
			}

			if _, err := pc.CreateMediaSink("A8-MINI", mediasink.RTSPSink(&duplexers.RTSPClientConfig{
				ServerAddr:        "localhost",
				ServerPort:        8554,
				StreamPath:        "DELIVERY/A8-MINI",
				ReadTimeout:       10 * time.Second,
				WriteTimeout:      10 * time.Second,
				DialTimeout:       10 * time.Second,
				ReconnectAttempts: 10,
				ReconnectDelay:    2 * time.Second,
				UserAgent:         "GoRTSP-Host/1.0",
			})); err != nil {
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

			mediapipe.NewAnyPipe(ctx, rl, wd).Start()
			mediapipe.NewAnyPipe(ctx, rd, wl).Start()

			gcs.WaitUntilClosed()
		}()
	}
}
