package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v4"

	"github.com/harshabose/simple_webrtc_comm/client/pkg"
	"github.com/harshabose/simple_webrtc_comm/cmd/delivery"
	"github.com/harshabose/simple_webrtc_comm/mediasink/pkg"
	"github.com/harshabose/simple_webrtc_comm/mediasink/pkg/rtsp"
	"github.com/harshabose/simple_webrtc_comm/mediasink/pkg/socket"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	socketServer := socket.NewServer(socket.Config{
		Addr:             "0.0.0.0",
		Port:             14552,
		ReadTimout:       100 * time.Millisecond,
		WriteTimout:      100 * time.Millisecond,
		TotalConnections: 10,
		KeepHosting:      true,
	})

	rtspServer := rtsp.NewNewServer(nil)

	go socketServer.Start(ctx)
	go rtspServer.Start(ctx)

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
		client.WithFirebaseAnswerSignal,
		client.WithMediaSinks(),
		client.WithDataChannels(),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("created peer connection")

	socketSink, err := pc.CreateMediaSink("TELEMETRY-RESTREAM", mediasink.WithLocalSocket("localhost", 14552, socket.HostConfig{
		WriteTimeout: 5 * time.Second,
		ReadTimout:   5 * time.Second,
		ConnectRetry: true,
	}, nil))
	if err != nil {
		panic(err)
	}

	// if _, err := pc.CreateMediaSink("A8-MINI", mediasink.WithWebRTCRestream(webrtc_js.Config{
	// 	Addr:             "0.0.0.0",
	// 	Port:             8080,
	// 	ReadTimout:       5 * time.Second,
	// 	WriteTimout:      5 * time.Second,
	// 	KeepHosting:      true,
	// 	TotalConnections: 10,
	// 	AllowLocalOnly:   false,
	// 	ClientTTL:        10 * time.Minute,
	// })); err != nil {
	// 	panic(err)
	// }
	//
	if _, err := pc.CreateDataChannel("MAVLINK", socketSink); err != nil {
		panic(err)
	}

	fmt.Println("connecting")
	if err := pc.Connect("DELIVERY"); err != nil {
		panic(err)
	}

	select {}
}
