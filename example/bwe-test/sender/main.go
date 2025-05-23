// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

//go:build !js
// +build !js

// bandwidth-estimation-from-disk demonstrates how to use Pion's Bandwidth Estimation APIs.
package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v4"
)

const (
	lowFile    = "low.ivf"
	lowBitrate = 300_000

	medFile    = "med.ivf"
	medBitrate = 1_000_000

	highFile    = "high.ivf"
	highBitrate = 2_500_000

	ivfHeaderSize = 32
)

// func main() { //nolint:gocognit,cyclop,maintidx
// 	qualityLevels := []struct {
// 		fileName string
// 		bitrate  int
// 	}{
// 		{lowFile, lowBitrate},
// 		{medFile, medBitrate},
// 		{highFile, highBitrate},
// 	}
// 	currentQuality := 0
//
// 	// Check if IVF files exist
// 	fmt.Println("🎬 Checking for video files...")
// 	for _, level := range qualityLevels {
// 		_, err := os.Stat(level.fileName)
// 		if os.IsNotExist(err) {
// 			panic(fmt.Sprintf("❌ File %s was not found", level.fileName))
// 		}
// 		fmt.Printf("✓ Found %s\n", level.fileName)
// 	}
//
// 	// Setup WebRTC
// 	fmt.Println("\n🔧 Setting up WebRTC...")
// 	interceptorRegistry := &interceptor.Registry{}
// 	mediaEngine := &webrtc.MediaEngine{}
// 	if err := mediaEngine.RegisterDefaultCodecs(); err != nil {
// 		panic(err)
// 	}
//
// 	congestionController, err := cc.NewInterceptor(func() (cc.BandwidthEstimator, error) {
// 		return gcc.NewSendSideBWE(gcc.SendSideBWEInitialBitrate(lowBitrate))
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	estimatorChan := make(chan cc.BandwidthEstimator, 1)
// 	congestionController.OnNewPeerConnection(func(id string, estimator cc.BandwidthEstimator) { // nolint: revive
// 		estimatorChan <- estimator
// 	})
//
// 	interceptorRegistry.Add(congestionController)
// 	if err = webrtc.ConfigureTWCCHeaderExtensionSender(mediaEngine, interceptorRegistry); err != nil {
// 		panic(err)
// 	}
//
// 	if err = webrtc.RegisterDefaultInterceptors(mediaEngine, interceptorRegistry); err != nil {
// 		panic(err)
// 	}
//
// 	// Create a new RTCPeerConnection
// 	peerConnection, err := webrtc.NewAPI(
// 		webrtc.WithInterceptorRegistry(interceptorRegistry), webrtc.WithMediaEngine(mediaEngine),
// 	).NewPeerConnection(webrtc.Configuration{
// 		ICEServers: []webrtc.ICEServer{
// 			{
// 				URLs: []string{"stun:stun.l.google.com:19302"},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer func() {
// 		if cErr := peerConnection.Close(); cErr != nil {
// 			fmt.Printf("cannot close peerConnection: %v\n", cErr)
// 		}
// 	}()
//
// 	// Wait until our Bandwidth Estimator has been created
// 	estimator := <-estimatorChan
//
// 	// Create a video track
// 	videoTrack, err := webrtc.NewTrackLocalStaticSample(
// 		webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "pion",
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	rtpSender, err := peerConnection.AddTrack(videoTrack)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	// Read incoming RTCP packets
// 	go func() {
// 		readRTCPWithAnalysis(rtpSender)
// 		// rtcpBuf := make([]byte, 1500)
// 		// for {
// 		// 	if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
// 		// 		return
// 		// 	}
// 		// }
// 	}()
//
// 	// Set the handler for ICE connection state
// 	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
// 		fmt.Printf("🔗 ICE Connection State: %s\n", connectionState.String())
// 	})
//
// 	// Set the handler for Peer connection state
// 	peerConnection.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
// 		fmt.Printf("📡 Peer Connection State: %s\n", state.String())
// 	})
//
// 	// ================================
// 	// FILE-BASED SDP EXCHANGE
// 	// ================================
//
// 	fmt.Println("\n=======================================================")
// 	fmt.Println("🎥 WebRTC Bandwidth Estimation Demo (File-based SDP)")
// 	fmt.Println("=======================================================")
// 	fmt.Println("Choose your role:")
// 	fmt.Println("1. Create an offer (generate offer.txt)")
// 	fmt.Println("2. Wait for an offer (read offer.txt and generate answer.txt)")
// 	fmt.Print("Enter choice (1 or 2): ")
//
// 	choice := readInput()
//
// 	var isOfferer bool
// 	switch strings.TrimSpace(choice) {
// 	case "1":
// 		isOfferer = true
// 		fmt.Println("✓ You chose to create an offer")
// 	case "2":
// 		isOfferer = false
// 		fmt.Println("✓ You chose to wait for an offer")
// 	default:
// 		panic("Invalid choice. Please enter 1 or 2.")
// 	}
//
// 	// File-based SDP exchange
// 	if isOfferer {
// 		handleOffererFlow(peerConnection)
// 	} else {
// 		handleAnswererFlow(peerConnection)
// 	}
//
// 	fmt.Println("\n🎬 Starting video streaming with bandwidth adaptation...")
// 	fmt.Println("📊 Monitor bandwidth changes and quality switching...")
// 	fmt.Println("Press Ctrl+C to stop")
//
// 	// ================================
// 	// VIDEO STREAMING WITH BWE
// 	// ================================
//
// 	// Open initial IVF file
// 	file, err := os.Open(qualityLevels[currentQuality].fileName)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	ivf, header, err := ivfreader.NewWith(file)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	// Calculate frame timing based on IVF header
// 	ticker := time.NewTicker(
// 		time.Millisecond * time.Duration((float32(header.TimebaseNumerator)/float32(header.TimebaseDenominator))*1000),
// 	)
// 	defer ticker.Stop()
//
// 	frame := []byte{}
// 	frameHeader := &ivfreader.IVFFrameHeader{}
// 	currentTimestamp := uint64(0)
// 	frameCount := 0
//
// 	switchQualityLevel := func(newQualityLevel int) {
// 		fmt.Printf("🔄 Switching quality: %s (%.1f Mbps) → %s (%.1f Mbps)\n",
// 			qualityLevels[currentQuality].fileName, float64(qualityLevels[currentQuality].bitrate)/1_000_000,
// 			qualityLevels[newQualityLevel].fileName, float64(qualityLevels[newQualityLevel].bitrate)/1_000_000,
// 		)
//
// 		currentQuality = newQualityLevel
// 		ivf.ResetReader(setReaderFile(qualityLevels[currentQuality].fileName))
//
// 		// Find a suitable frame to start from (keyframe preferred)
// 		for {
// 			if frame, frameHeader, err = ivf.ParseNextFrame(); err != nil {
// 				break
// 			} else if frameHeader.Timestamp >= currentTimestamp && frame[0]&0x1 == 0 {
// 				break
// 			}
// 		}
// 	}
//
// 	// Start streaming loop
// 	// lastBitrateReport := time.Now()
// 	for ; true; <-ticker.C {
// 		targetBitrate := estimator.GetTargetBitrate()
//
// 		// Report bitrate every 5 seconds
// 		// if time.Since(lastBitrateReport) > 5*time.Second {
// 		fmt.Printf("📊 Target bitrate: %f Mbps | Current quality: %s | Frames sent: %d\n",
// 			float64(targetBitrate)/1000000,
// 			qualityLevels[currentQuality].fileName,
// 			frameCount,
// 		)
// 		targetBitrate = lowBitrate
// 		// lastBitrateReport = time.Now()
// 		// }
//
// 		// Bandwidth-based quality switching
// 		switch {
// 		// Downgrade quality if target bitrate is below current level
// 		case currentQuality != 0 && targetBitrate < qualityLevels[currentQuality].bitrate:
// 			switchQualityLevel(currentQuality - 1)
//
// 		// Upgrade quality if target bitrate is above next level
// 		case len(qualityLevels) > (currentQuality+1) && targetBitrate > qualityLevels[currentQuality+1].bitrate:
// 			switchQualityLevel(currentQuality + 1)
//
// 		// Normal frame processing
// 		default:
// 			frame, frameHeader, err = ivf.ParseNextFrame()
// 		}
//
// 		switch {
// 		// Loop the video when we reach the end
// 		case errors.Is(err, io.EOF):
// 			ivf.ResetReader(setReaderFile(qualityLevels[currentQuality].fileName))
//
// 		// Send the frame
// 		case err == nil:
// 			currentTimestamp = frameHeader.Timestamp
// 			if err = videoTrack.WriteSample(media.Sample{Data: frame, Duration: time.Second}); err != nil {
// 				// Don't panic on frame errors, just log them
// 				fmt.Printf("⚠ Error sending frame: %v\n", err)
// 			} else {
// 				frameCount++
// 			}
//
// 		// Handle other errors
// 		default:
// 			fmt.Printf("⚠ Frame parsing error: %v\n", err)
// 		}
// 	}
// }

// Handle the offerer (initiator) flow
func handleOffererFlow(pc *webrtc.PeerConnection) {
	fmt.Println("\n⚡ Creating offer...")
	offer, err := pc.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	if err = pc.SetLocalDescription(offer); err != nil {
		panic(err)
	}

	// Wait for ICE gathering to complete
	gatherComplete := webrtc.GatheringCompletePromise(pc)
	fmt.Println("⏳ Gathering ICE candidates...")
	<-gatherComplete
	fmt.Println("✓ ICE gathering complete")

	// Save offer to file
	saveSDPToFile(pc.LocalDescription(), "offer.txt")
	fmt.Println("\n📤 Offer saved to offer.txt")
	fmt.Println("📋 Instructions:")
	fmt.Println("   1. Send offer.txt to the remote peer")
	fmt.Println("   2. Wait for them to send you answer.txt")
	fmt.Println("   3. Place answer.txt in this directory")
	fmt.Println("   4. Press Enter to continue...")

	// Wait for user to place answer.txt
	waitForFile("answer.txt")

	// Load answer from file
	answer := loadSDPFromFile("answer.txt")
	if err = pc.SetRemoteDescription(answer); err != nil {
		panic(fmt.Sprintf("❌ Error setting remote description: %v", err))
	}
	fmt.Println("✓ Answer processed successfully")
}

// Handle the answerer (responder) flow
func handleAnswererFlow(pc *webrtc.PeerConnection) {
	fmt.Println("\n📥 Waiting for offer...")
	fmt.Println("📋 Instructions:")
	fmt.Println("   1. Obtain offer.txt from the remote peer")
	fmt.Println("   2. Place offer.txt in this directory")
	fmt.Println("   3. Press Enter to continue...")

	// Wait for user to place offer.txt
	waitForFile("offer.txt")

	// Load offer from file
	offer := loadSDPFromFile("offer.txt")
	if err := pc.SetRemoteDescription(offer); err != nil {
		panic(fmt.Sprintf("❌ Error setting remote description: %v", err))
	}

	// Create answer
	fmt.Println("⚡ Creating answer...")
	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Wait for ICE gathering
	gatherComplete := webrtc.GatheringCompletePromise(pc)
	if err = pc.SetLocalDescription(answer); err != nil {
		panic(err)
	}

	fmt.Println("⏳ Gathering ICE candidates...")
	<-gatherComplete
	fmt.Println("✓ ICE gathering complete")

	// Save answer to file
	saveSDPToFile(pc.LocalDescription(), "answer.txt")
	fmt.Println("\n📤 Answer saved to answer.txt")
	fmt.Println("📋 Instructions:")
	fmt.Println("   1. Send answer.txt back to the remote peer")
	fmt.Println("   2. Connection should establish automatically")
}

// Save SDP to file
func saveSDPToFile(sdp *webrtc.SessionDescription, filename string) {
	encoded := encode(sdp)
	err := os.WriteFile(filename, []byte(encoded), 0644)
	if err != nil {
		panic(fmt.Sprintf("❌ Error saving SDP to %s: %v", filename, err))
	}
	fmt.Printf("✓ SDP saved to %s (%d bytes)\n", filename, len(encoded))
}

// Load SDP from file
func loadSDPFromFile(filename string) webrtc.SessionDescription {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("❌ Error reading %s: %v", filename, err))
	}

	var sdp webrtc.SessionDescription
	encoded := strings.TrimSpace(string(data))
	decode(encoded, &sdp)
	fmt.Printf("✓ SDP loaded from %s (%d bytes)\n", filename, len(data))
	return sdp
}

// Wait for file to exist and user confirmation
func waitForFile(filename string) {
	for {
		readInput() // Wait for user to press Enter

		if _, err := os.Stat(filename); err == nil {
			fmt.Printf("✓ Found %s\n", filename)
			break
		} else {
			fmt.Printf("❌ %s not found. Please place the file and press Enter again...\n", filename)
		}
	}
}

// Simple input reader
func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func setReaderFile(filename string) func(_ int64) io.Reader {
	return func(_ int64) io.Reader {
		file, err := os.Open(filename) // nolint
		if err != nil {
			panic(err)
		}
		if _, err = file.Seek(ivfHeaderSize, io.SeekStart); err != nil {
			panic(err)
		}
		return file
	}
}

// JSON encode + base64 a SessionDescription
func encode(obj *webrtc.SessionDescription) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

// Decode a base64 and unmarshal JSON into a SessionDescription
func decode(in string, obj *webrtc.SessionDescription) {
	fmt.Printf("🔍 Decoding SDP (%d chars)\n", len(in))

	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(fmt.Sprintf("❌ Base64 decode error: %v", err))
	}

	if err = json.Unmarshal(b, obj); err != nil {
		panic(fmt.Sprintf("❌ JSON unmarshal error: %v", err))
	}

	fmt.Printf("✓ SDP parsed successfully (type: %s)\n", obj.Type)
}

// Enhanced RTCP reader that parses and analyzes TWCC feedback
func readRTCPWithAnalysis(rtpSender *webrtc.RTPSender) {
	rtcpBuf := make([]byte, 1500)
	lastTWCCReport := time.Now()

	for {
		n, _, rtcpErr := rtpSender.Read(rtcpBuf)
		if rtcpErr != nil {
			fmt.Printf("RTCP read error: %v\n", rtcpErr)
			return
		}

		// Parse the RTCP packet
		packets, err := rtcp.Unmarshal(rtcpBuf[:n])
		if err != nil {
			fmt.Printf("Failed to unmarshal RTCP: %v\n", err)
			continue
		}

		// Analyze each RTCP packet
		for _, packet := range packets {
			switch p := packet.(type) {
			case *rtcp.TransportLayerCC:
				// TWCC feedback packet - this is what we're interested in!
				analyzeTWCCFeedback(p)
				lastTWCCReport = time.Now()

			case *rtcp.ReceiverReport:
				analyzeReceiverReport(p)

			case *rtcp.SenderReport:
				analyzeSenderReport(p)

			case *rtcp.ReceiverEstimatedMaximumBitrate:
				// REMB packet
				analyzeREMBFeedback(p)

			default:
				fmt.Printf("📄 Other RTCP packet: %T\n", packet)
			}
		}

		// Check if we're missing TWCC feedback
		if time.Since(lastTWCCReport) > 5*time.Second {
			fmt.Printf("⚠️  No TWCC feedback received for %v - BWE may be stuck!\n",
				time.Since(lastTWCCReport))
		}
	}
}

func analyzeTWCCFeedback(twcc *rtcp.TransportLayerCC) {
	fmt.Printf("📈 TWCC Feedback Received:\n")
	fmt.Printf("  Media SSRC: %d\n", twcc.MediaSSRC)
	fmt.Printf("  Feedback Packet Count: %d\n", twcc.FbPktCount)
	fmt.Printf("  Packet Status Count: %d\n", twcc.PacketStatusCount)
	fmt.Printf("  Reference Time: %d\n", twcc.ReferenceTime)

	// Count received packets
	receivedCount := 0
	lostCount := 0

	for i, _ := range twcc.PacketChunks {
		if i >= int(twcc.PacketStatusCount) {
			break
		}
	}

	totalPackets := receivedCount + lostCount
	packetLossRate := float64(lostCount) / float64(totalPackets) * 100

	fmt.Printf("  Packets: %d received, %d lost (%.1f%% loss)\n",
		receivedCount, lostCount, packetLossRate)

	// Analyze receive deltas for jitter
	if len(twcc.RecvDeltas) > 0 {
		analyzeReceiveDeltas(twcc.RecvDeltas)
	}

	fmt.Println("  ----------------------------------------")
}

func analyzeReceiveDeltas(deltas []*rtcp.RecvDelta) {
	if len(deltas) < 2 {
		return
	}

	// Calculate jitter and timing information
	var totalDelta time.Duration
	var maxDelta, minDelta time.Duration = 0, time.Hour

	for _, delta := range deltas {
		deltaDuration := delta.Delta
		totalDelta += time.Duration(deltaDuration)

		if time.Duration(deltaDuration) > (maxDelta) {
			maxDelta = time.Duration(deltaDuration)
		}
		if time.Duration(deltaDuration) < minDelta {
			minDelta = time.Duration(deltaDuration)
		}
	}

	avgDelta := totalDelta / time.Duration(len(deltas))
	jitter := maxDelta - minDelta

	fmt.Printf("  Timing Analysis:\n")
	fmt.Printf("    Avg Delta: %v\n", avgDelta)
	fmt.Printf("    Jitter: %v (max: %v, min: %v)\n", jitter, maxDelta, minDelta)

	// Warning for high jitter
	if jitter > 50*time.Millisecond {
		fmt.Printf("    ⚠️  High jitter detected: %v\n", jitter)
	}
}

func analyzeReceiverReport(rr *rtcp.ReceiverReport) {
	fmt.Printf("📊 Receiver Report:\n")
	fmt.Printf("  SSRC: %d\n", rr.SSRC)

	for _, report := range rr.Reports {
		fmt.Printf("  Report for SSRC %d:\n", report.SSRC)
		fmt.Printf("    Fraction Lost: %d/256 (%.1f%%)\n",
			report.FractionLost, float64(report.FractionLost)/256*100)
		fmt.Printf("    Total Lost: %d packets\n", report.TotalLost)
		fmt.Printf("    Highest Seq: %d\n", report.LastSequenceNumber)
		fmt.Printf("    Jitter: %d timestamp units\n", report.Jitter)
	}
	fmt.Println("  ----------------------------------------")
}

func analyzeSenderReport(sr *rtcp.SenderReport) {
	fmt.Printf("📤 Sender Report:\n")
	fmt.Printf("  SSRC: %d\n", sr.SSRC)
	fmt.Printf("  NTP Time: %d\n", sr.NTPTime)
	fmt.Printf("  RTP Time: %d\n", sr.RTPTime)
	fmt.Printf("  Packet Count: %d\n", sr.PacketCount)
	fmt.Printf("  Octet Count: %d\n", sr.OctetCount)
	fmt.Println("  ----------------------------------------")
}

func analyzeREMBFeedback(remb *rtcp.ReceiverEstimatedMaximumBitrate) {
	fmt.Printf("📶 REMB Feedback:\n")
	fmt.Printf("  Sender SSRC: %d\n", remb.SenderSSRC)
	fmt.Printf("  Bitrate: %.2f Mbps\n", float64(remb.Bitrate)/1_000_000)
	fmt.Printf("  Media SSRCs: %v\n", remb.SSRCs)
	fmt.Println("  ----------------------------------------")
}
