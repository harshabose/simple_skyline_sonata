package main

import (
	"context"
	"github.com/asticode/go-astiav"
	"github.com/harshabose/simple_webrtc_comm/client/pkg"
	"github.com/harshabose/simple_webrtc_comm/cmd/delivery"
	"github.com/harshabose/simple_webrtc_comm/datachannel/pkg"
	"github.com/harshabose/simple_webrtc_comm/mediasource/pkg"
	"github.com/harshabose/simple_webrtc_comm/transcode/pkg"
	"os"
	"time"
)

func main() {
	ctx := context.Background()

	deliveryDrone, err := client.CreateClient(
		ctx, nil, nil,
		client.WithBandwidthControlInterceptor(2500, 50*time.Millisecond),
		client.WithH264MediaEngine(delivery.DefaultVideoClockRate, client.PacketisationMode1, client.ProfileLevelBaseline42),
		client.WithNACKInterceptor(client.NACKGeneratorLowLatency, client.NACKResponderLowLatency),
		client.WithFLEXFECInterceptor(),
		client.WithRTCPReportsInterceptor(client.RTCPReportIntervalLowLatency),
		client.WithTWCCSenderInterceptor(client.TWCCIntervalLowLatency),
	)
	if err != nil {
		panic(err)
	}

	pc, err := deliveryDrone.CreatePeerConnection(
		"MAIN",
		client.WithRTCConfiguration(client.GetRTCConfiguration()),
		client.WithOfferSignal,
		client.WithMediaSources(),
		client.WithDataChannels(),
	)
	if err != nil {
		panic(err)
	}

	if err := pc.CreateDataChannel("MAVLINK",
		data.WithRandomBindPort,
		data.WithMAVP2P(os.Getenv("MAVP2P_EXE_PATH"), os.Getenv("MAVLINK_SERIAL")),
	); err != nil {
		panic(err)
	}

	if err := pc.CreateMediaSource("A8-MINI", true,
		mediasource.WithH264Track(delivery.DefaultVideoClockRate, mediasource.PacketisationMode1, mediasource.ProfileLevelBaseline42),
		mediasource.WithPriority(mediasource.Level5),
		mediasource.WithStream(
			mediasource.WithBufferSize(int(delivery.DefaultVideoFPS*3)),
			mediasource.WithDemuxer(
				"rtsp://192.168.144.25:8554/main.264",
				transcode.WithRTSPInputOption,
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

	deliveryDrone.WaitUntilClosed()
}
