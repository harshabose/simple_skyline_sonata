# Simple Skyline Sonata

A lightweight, embedded-friendly WebRTC implementation designed specifically for drone-to-any networking and webrtc-mesh networking.

**For Internal Use Only: Indian Institute of Technology Bombay**

## Overview

Simple Skyline Sonata is a specialized WebRTC solution built on [Pion WebRTC](https://github.com/pion/webrtc), optimized for creating robust real-time communication meshes between drones and ground control stations (GCS). The project focuses on minimal resource consumption, high reliability, and low-latency data exchange in aerial networking scenarios.

## Key Features

- **Embedded-Friendly**: Designed to operate efficiently on constrained hardware found in drone systems
- **Mesh Networking**: Creates peer-to-peer connections between multiple drones and ground stations
- **Cross-Platform Support**: Works across Linux, Windows, and macOS environments
- **Multiple Distribution Options**:
  - Authorized SDK access for custom integration
  - Single-file standalone binaries for easy deployment
- **Minimal Dependencies**: Self-contained with minimal external requirements
- **Optimized for Aerial Communications**: Handles intermittent connectivity and bandwidth constraints common in aerial applications
- **Use-case Agnostic**: Developer friendly SDKs with various capabilities for data, video and audio transmission

## Architecture

Simple Skyline Sonata uses a decentralized mesh architecture where each node (drone or GCS) can establish direct WebRTC connections with any other node in the network. This approach:

- Eliminates single points of failure
- Reduces overall latency
- Improves network resilience
- Maximizes bandwidth utilization

The core is built on Pion WebRTC, a pure Go implementation of the WebRTC API, providing excellent performance characteristics for embedded systems.

## Installation

### Standalone Binaries

Pre-compiled binaries are available for authorized users:

| Platform | Architectures | Download Link |
|----------|---------------|--------------|
| Linux    | AMD64, ARM64  | [Internal Portal](#) |
| Windows  | All           | [Internal Portal](#) |
| macOS    | All           | [Internal Portal](#) |

### SDK Integration

For custom integration into your drone platform, SDKs are available for authorized developers. Please contact the project the developer for approval.

## Usage

### Basic Connection Example

```go
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
		client.WithBandwidthControlInterceptor(5_000_000, time.Second),
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
		data.WithMAVP2P(os.Getenv("MAVP2P_EXE_PATH"), os.Getenv("MAVLINK_SERIAL")),
	); err != nil {
		panic(err)
	}

	if err := pc.CreateMediaSource("A8-MINI", false,
		mediasource.WithH264Track(delivery.DefaultVideoClockRate, mediasource.PacketisationMode1, mediasource.ProfileLevelBaseline41),
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

	drone.WaitUntilClosed()
}

```

Documentation available with approval.

## Performance Considerations

For optimal performance in drone applications:

- Set appropriate bandwidth limits based on your hardware capabilities
- Consider reducing video quality for longer range operations
- Enable the adaptive bitrate option for variable network conditions
- Use mesh relaying only when necessary to conserve power

## Binary Distribution

Standalone binaries are provided as single-file executables with no external dependencies:

- **Linux**: Static binaries for AMD64 and ARM64 architectures
- **Windows**: Portable EXE files for all supported architectures
- **macOS**: Universal binaries supporting Intel and Apple Silicon

## Future Work

The following features are planned for upcoming releases:

- **Easy Frontend**: A web-based interface for monitoring and controlling the mesh network
- **Bandwidth Optimization**: Advanced congestion control algorithms specific to aerial networks
- **Mesh Visualization**: Real-time network topology visualization
- **Security Enhancements**: Additional authentication and encryption options

## Security Notice

This project is for internal use only at IIT Bombay. Unauthorized access or distribution is prohibited. All connections require proper authentication credentials.

## Support

For technical support, feature requests, or bug reports, please contact:

- Project Maintainers: [Contact Information]
- Issue Tracker: [Internal Link]

## License

© Indian Institute of Technology Bombay. All rights reserved.
