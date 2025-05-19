package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v4"

	"github.com/harshabose/simple_webrtc_comm/client/pkg"
	"github.com/harshabose/simple_webrtc_comm/cmd/delivery"
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
				client.WithRTCConfiguration(client.GetRTCConfiguration()),
				client.WithAnswerSignal,
				client.WithMediaSinks(),
				client.WithDataChannels(),
			)
			if err != nil {
				panic(err)
			}
			fmt.Println("created peer connection")

			if _, err := pc.CreateDataChannel("MAVLINK", data.WithRandomBindPort, data.WithLoopBackPort(14550)); err != nil {
				panic(err)
			}

			if err := pc.Connect("DELIVERY"); err != nil {
				panic(err)
			}

			gcs.WaitUntilClosed()
		}()
	}

}
