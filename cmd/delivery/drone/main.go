package main

import (
	"context"
	
	"github.com/asticode/go-astiav"
	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v4"
	
	"github.com/harshabose/simple_webrtc_comm/client/pkg"
	"github.com/harshabose/simple_webrtc_comm/cmd/delivery"
	"github.com/harshabose/simple_webrtc_comm/datachannel/pkg"
	"github.com/harshabose/simple_webrtc_comm/mediasource/pkg"
	"github.com/harshabose/simple_webrtc_comm/transcode/pkg"
)

func main() {
	ctx := context.Background()
	
	mediaEngine := &webrtc.MediaEngine{}
	registry := &interceptor.Registry{}
	
	drone, err := client.CreateClient(
		ctx, mediaEngine, registry,
		// client.WithBandwidthControlInterceptor(5_000_000, time.Second),
		client.WithH264MediaEngine(delivery.DefaultVideoClockRate, client.PacketisationMode1, client.ProfileLevelBaseline41, delivery.DefaultSPSBase64, delivery.DefaultPPSBase64),
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
		client.WithOfferSignal,
		client.WithMediaSources(),
		client.WithDataChannels(),
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
	
	if err := pc.CreateMediaSource("A8-MINI", false,
		mediasource.WithH264Track(delivery.DefaultVideoClockRate, mediasource.PacketisationMode1, mediasource.ProfileLevelBaseline41),
		mediasource.WithPriority(mediasource.Level5),
		mediasource.WithStream(
			mediasource.WithBufferSize(int(delivery.DefaultVideoFPS*3)),
			mediasource.WithDemuxer(
				"/dev/video0",
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
				transcode.WithX264LowLatencyOptions,
			),
		),
	); err != nil {
		panic(err)
	}
	
	if err := pc.Connect("DELIVERY"); err != nil {
		panic(err)
	}
	
	drone.WaitUntilClosed()
}
