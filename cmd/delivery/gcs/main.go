package main

import (
	"context"
	"github.com/harshabose/simple_webrtc_comm/client/pkg"
	"github.com/harshabose/simple_webrtc_comm/cmd/delivery"
	"github.com/harshabose/simple_webrtc_comm/datachannel/pkg"
)

func main() {
	ctx := context.Background()

	gcs, err := client.CreateClient(
		ctx, nil, nil,
		client.WithH264MediaEngine(delivery.DefaultVideoClockRate, client.PacketisationMode1, client.ProfileLevelBaseline42),
		client.WithNACKInterceptor(client.NACKGeneratorLowLatency, client.NACKResponderLowLatency),
		client.WithFLEXFECInterceptor(),
		client.WithRTCPReportsInterceptor(client.RTCPReportIntervalLowLatency),
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

	if err := pc.CreateDataChannel("MAVLINK", data.WithRandomBindPort, data.WithLoopBackPort(14550)); err != nil {
		panic(err)
	}

	if err := pc.Connect("DELIVERY"); err != nil {
		panic(err)
	}

	gcs.WaitUntilClosed()
}
