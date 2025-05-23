package main

import (
	"context"
	"fmt"
	"time"

	"github.com/asticode/go-astiav"
	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v4"

	client "github.com/harshabose/simple_webrtc_comm/client/pkg"
	mediasource "github.com/harshabose/simple_webrtc_comm/mediasource/pkg"
	transcode "github.com/harshabose/simple_webrtc_comm/transcode/pkg"
)

const (
	_cameraBitrate = 300_000
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("\n🔧 Setting up WebRTC...")
	interceptorRegistry := &interceptor.Registry{}
	mediaEngine := &webrtc.MediaEngine{}

	drone, err := client.CreateClient(ctx, cancel, mediaEngine, interceptorRegistry,
		client.WithVP8MediaEngine(90000),
		client.WithBandwidthControlInterceptor(_cameraBitrate, time.Second),
		client.WithTWCCHeaderExtensionSender(),
		client.WithNACKInterceptor(client.NACKGeneratorLowLatency, client.NACKResponderLowLatency),
		client.WithRTCPReportsInterceptor(client.RTCPReportIntervalLowLatency),
		client.WithTWCCSenderInterceptor(client.TWCCIntervalLowLatency),
	)
	if err != nil {
		panic(err)
	}

	pc, err := drone.CreatePeerConnection("MAIN",
		client.WithRTCConfiguration(client.GetRTCConfiguration()),
		client.WithFileOfferSignal("/Users/harshabose/GolandProjects/simple_skyline_sonata/example/bwe-test/exchange/", "/Users/harshabose/GolandProjects/simple_skyline_sonata/example/bwe-test/exchange"),
		client.WithMediaSources(),
		client.WithBandwidthControl(),
	)
	if err != nil {
		panic(err)
	}

	if _, err := pc.CreateMediaSource("A8-MINI", true,
		mediasource.WithVP8Track(90000),
		mediasource.WithPriority(mediasource.Level5),
		mediasource.WithStream(
			mediasource.WithBufferSize(90),
			mediasource.WithDemuxer(
				"0",
				transcode.WithAvFoundationInputFormatOption,
				transcode.WithDemuxerBufferSize(90),
			),
			mediasource.WithDecoder(transcode.WithDecoderBufferSize(90)),
			mediasource.WithFilter(
				transcode.VideoFilters,
				transcode.WithFilterBufferSize(90),
				transcode.WithVideoScaleFilterContent(1280, 720),
				transcode.WithVideoPixelFormatFilterContent(astiav.PixelFormatYuv420P),
				transcode.WithVideoFPSFilterContent(30),
			),
			mediasource.WithEncoder(
				astiav.CodecIDH264,
				transcode.WithEncoderBufferSize(90),
				transcode.WithWebRTCOptimisedOptions,
			),
		),
	); err != nil {
		panic(err)
	}

	estimator, err := pc.GetBWEstimator()
	if err != nil {
		panic(err)
	}

	if err := pc.Connect("TEST"); err != nil {
		panic(err)
	}

	frameCount := 0
	lastBitrateReport := time.Now()

	// Start streaming loop
	for {
		// Report bitrate every 5 seconds
		if time.Since(lastBitrateReport) > 5*time.Second {
			targetBitrate := estimator.GetTargetBitrate()
			fmt.Printf("📊 Target bitrate: %f Mbps | Frames sent: %d\n",
				float64(targetBitrate)/1000000,
				frameCount,
			)
			lastBitrateReport = time.Now()
		}

		// Small sleep to prevent CPU overload
		time.Sleep(10 * time.Millisecond)
	}
}
