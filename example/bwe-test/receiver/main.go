package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pion/webrtc/v4"
)

// func main() {
// 	// Show current working directory for debugging
// 	cwd, _ := os.Getwd()
// 	fmt.Printf("📁 Current working directory: %s\n", cwd)
//
// 	fmt.Println("\n🔧 Setting up WebRTC...")
// 	interceptorRegistry := &interceptor.Registry{}
// 	mediaEngine := &webrtc.MediaEngine{}
// 	if err := mediaEngine.RegisterDefaultCodecs(); err != nil {
// 		panic(err)
// 	}
//
// 	if err := webrtc.ConfigureTWCCHeaderExtensionSender(mediaEngine, interceptorRegistry); err != nil {
// 		panic(err)
// 	}
//
// 	if err := webrtc.RegisterDefaultInterceptors(mediaEngine, interceptorRegistry); err != nil {
// 		panic(err)
// 	}
//
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
// 		if err := peerConnection.Close(); err != nil {
// 			fmt.Printf("cannot close peerConnection: %v\n", err)
// 		}
// 	}()
//
// 	peerConnection.OnTrack(func(remote *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
// 		fmt.Printf("📹 Received track: %s (SSRC: %d)\n", remote.Kind(), remote.SSRC())
//
// 		// Read RTP packets from track
// 		go func() {
// 			for {
// 				_, _, err := remote.ReadRTP()
// 				if err != nil {
// 					return
// 				}
// 				// Track received packets for debugging
// 			}
// 		}()
// 	})
//
// 	peerConnection.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
// 		fmt.Printf("🔗 ICE Connection State: %s\n", state.String())
// 	})
//
// 	peerConnection.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
// 		fmt.Printf("📡 Peer Connection State: %s\n", state.String())
// 	})
//
// 	// Try multiple possible locations for offer.txt
// 	possiblePaths := []string{
// 		"offer.txt",    // Same directory
// 		"../offer.txt", // Parent directory
// 		"../sender/offer.txt",
// 		"../../offer.txt", // Grandparent directory
// 		"./offer.txt",     // Explicit current directory
// 	}
//
// 	fmt.Println("\n📥 Looking for offer.txt...")
// 	var offerPath string
// 	for _, path := range possiblePaths {
// 		absPath, _ := filepath.Abs(path)
// 		fmt.Printf("🔍 Checking: %s\n", absPath)
// 		if _, err := os.Stat(path); err == nil {
// 			offerPath = path
// 			fmt.Printf("✓ Found offer.txt at: %s\n", absPath)
// 			break
// 		}
// 	}
//
// 	if offerPath == "" {
// 		fmt.Println("\n❌ offer.txt not found in any expected location")
// 		fmt.Println("📋 Please place offer.txt in one of these locations:")
// 		for _, path := range possiblePaths {
// 			absPath, _ := filepath.Abs(path)
// 			fmt.Printf("   - %s\n", absPath)
// 		}
// 		fmt.Println("\nPress Enter to check again...")
// 		waitForFile(possiblePaths)
// 	}
//
// 	// Load offer from file
// 	offer := loadSDPFromFile(offerPath)
// 	if err := peerConnection.SetRemoteDescription(offer); err != nil {
// 		panic(fmt.Sprintf("❌ Error setting remote description: %v", err))
// 	}
//
// 	// Create answer
// 	fmt.Println("⚡ Creating answer...")
// 	answer, err := peerConnection.CreateAnswer(nil)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	// Wait for ICE gathering
// 	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
// 	if err = peerConnection.SetLocalDescription(answer); err != nil {
// 		panic(err)
// 	}
//
// 	fmt.Println("⏳ Gathering ICE candidates...")
// 	<-gatherComplete
// 	fmt.Println("✓ ICE gathering complete")
//
// 	// Save answer in the same directory as offer
// 	answerPath := filepath.Join(filepath.Dir(offerPath), "../sender/answer.txt")
// 	saveSDPToFile(peerConnection.LocalDescription(), answerPath)
//
// 	fmt.Println("\n✓ Receiver setup complete. Waiting for connection...")
// 	fmt.Println("Press Ctrl+C to exit")
//
// 	// Keep the program running
// 	select {}
// }

// Wait for file to exist in any of the possible locations
func waitForFile(possiblePaths []string) string {
	for {
		readInput() // Wait for user to press Enter

		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				absPath, _ := filepath.Abs(path)
				fmt.Printf("✓ Found %s\n", absPath)
				return path
			}
		}

		fmt.Println("❌ Still not found. Please check the file placement and press Enter again...")

		// Show what files are actually in each directory
		fmt.Println("\n🔍 Files in current directory:")
		if files, err := os.ReadDir("."); err == nil {
			for _, file := range files {
				fmt.Printf("   - %s\n", file.Name())
			}
		}

		fmt.Println("\n🔍 Files in parent directory:")
		if files, err := os.ReadDir(".."); err == nil {
			for _, file := range files {
				fmt.Printf("   - %s\n", file.Name())
			}
		}
	}
}

// Simple input reader with timeout
func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// Save SDP to file
func saveSDPToFile(sdp *webrtc.SessionDescription, filename string) {
	encoded := encode(sdp)

	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(fmt.Sprintf("❌ Error creating directory %s: %v", dir, err))
	}

	err := os.WriteFile(filename, []byte(encoded), 0644)
	if err != nil {
		panic(fmt.Sprintf("❌ Error saving SDP to %s: %v", filename, err))
	}

	absPath, _ := filepath.Abs(filename)
	fmt.Printf("✓ SDP saved to %s (%d bytes)\n", absPath, len(encoded))
}

// JSON encode + base64 a SessionDescription
func encode(obj *webrtc.SessionDescription) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

// Load SDP from file
func loadSDPFromFile(filename string) webrtc.SessionDescription {
	absPath, _ := filepath.Abs(filename)
	fmt.Printf("📖 Reading SDP from: %s\n", absPath)

	data, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("❌ Error reading %s: %v", filename, err))
	}

	var sdp webrtc.SessionDescription
	encoded := strings.TrimSpace(string(data))
	decode(encoded, &sdp)
	fmt.Printf("✓ SDP loaded from %s (%d bytes)\n", absPath, len(data))
	return sdp
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
