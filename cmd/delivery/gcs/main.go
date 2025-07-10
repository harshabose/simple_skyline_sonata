package main

import (
	"context"
	"math"
	"time"

	"github.com/harshabose/mediapipe"
	"github.com/harshabose/mediapipe/pkg/dioreader"
	"github.com/harshabose/mediapipe/pkg/diowriter"
	"github.com/harshabose/mediapipe/pkg/loopback"
	"github.com/harshabose/mediapipe/pkg/rtsp"
	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v4"

	"github.com/harshabose/simple_webrtc_comm/client/pkg"
	"github.com/harshabose/simple_webrtc_comm/client/pkg/mediasink"
	"github.com/harshabose/simple_webrtc_comm/cmd/delivery"
)

func main() {
	for {
		func() {
			time.Sleep(2 * time.Second)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			l, err := loopback.NewLoopBack(context.Background(), "127.0.0.1:0", loopback.WithLoopBackPort(14550))
			if err != nil {
				panic(err)
			}

			mediaEngine := &webrtc.MediaEngine{}
			registry := &interceptor.Registry{}
			settings := &webrtc.SettingEngine{}

			gcs, err := client.CreateClient(
				ctx, cancel, mediaEngine, registry, settings,
				client.WithH264MediaEngine(delivery.DefaultVideoClockRate, client.PacketisationMode1, client.ProfileLevelBaseline31, delivery.DefaultSPSBase64, delivery.DefaultPPSBase64),
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

			if _, err := pc.CreateMediaSink("A8-MINI", mediasink.RTSPSink(&rtsp.ClientConfig{
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

			mediapipe.NewAnyPipe(ctx, rl, wd).Start()
			mediapipe.NewAnyPipe(ctx, rd, wl).Start()

			gcs.WaitUntilClosed()
		}()
	}
}
