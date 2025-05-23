package main

import (
	"context"
	"fmt"

	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v4"

	client "github.com/harshabose/simple_webrtc_comm/client/pkg"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("\n🔧 Setting up WebRTC...")
	interceptorRegistry := &interceptor.Registry{}
	mediaEngine := &webrtc.MediaEngine{}

	gcs, err := client.CreateClient(ctx, cancel, mediaEngine, interceptorRegistry,
		client.WithVP8MediaEngine(90000),
		client.WithTWCCHeaderExtensionSender(),
		client.WithNACKInterceptor(client.NACKGeneratorLowLatency, client.NACKResponderLowLatency),
		client.WithRTCPReportsInterceptor(client.RTCPReportIntervalLowLatency),
		client.WithTWCCSenderInterceptor(client.TWCCIntervalLowLatency),
	)
	if err != nil {
		panic(err)
	}

	pc, err := gcs.CreatePeerConnection("MAIN",
		client.WithRTCConfiguration(client.GetRTCConfiguration()),
		client.WithFileAnswerSignal("/Users/harshabose/GolandProjects/simple_skyline_sonata/example/bwe-test/exchange/", "/Users/harshabose/GolandProjects/simple_skyline_sonata/example/bwe-test/exchange"),
		client.WithMediaSinks(),
	)
	if err != nil {
		panic(err)
	}

	// peerConnection, err := pc.GetPeerConnection()
	// if err != nil {
	// 	panic(err)
	// }
	//
	// // Try multiple possible locations for offer.txt
	// possiblePaths := []string{
	// 	"../exchange/TEST/MAIN/offer.txt", // Same directory
	// 	"../offer.txt",                    // Parent directory
	// 	"../sender/offer.txt",
	// 	"../../offer.txt", // Grandparent directory
	// 	"./offer.txt",     // Explicit current directory
	// }
	//
	// fmt.Println("\n📥 Looking for offer.txt...")
	// var offerPath string
	// for _, path := range possiblePaths {
	// 	absPath, _ := filepath.Abs(path)
	// 	fmt.Printf("🔍 Checking: %s\n", absPath)
	// 	if _, err := os.Stat(path); err == nil {
	// 		offerPath = path
	// 		fmt.Printf("✓ Found offer.txt at: %s\n", absPath)
	// 		break
	// 	}
	// }
	//
	// if offerPath == "" {
	// 	fmt.Println("\n❌ offer.txt not found in any expected location")
	// 	fmt.Println("📋 Please place offer.txt in one of these locations:")
	// 	for _, path := range possiblePaths {
	// 		absPath, _ := filepath.Abs(path)
	// 		fmt.Printf("   - %s\n", absPath)
	// 	}
	// 	fmt.Println("\nPress Enter to check again...")
	// 	waitForFile(possiblePaths)
	// }
	//
	// // Load offer from file
	// offer := loadSDPFromFile(offerPath)
	// if err := peerConnection.SetRemoteDescription(offer); err != nil {
	// 	panic(fmt.Sprintf("❌ Error setting remote description: %v", err))
	// }
	//
	// // Create answer
	// fmt.Println("⚡ Creating answer...")
	// answer, err := peerConnection.CreateAnswer(nil)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// // Wait for ICE gathering
	// gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
	// if err = peerConnection.SetLocalDescription(answer); err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Println("⏳ Gathering ICE candidates...")
	// <-gatherComplete
	// fmt.Println("✓ ICE gathering complete")
	//
	// // Save answer in the same directory as offer
	// answerPath := filepath.Join(filepath.Dir(offerPath), "answer.txt")
	// saveSDPToFile(peerConnection.LocalDescription(), answerPath)
	//
	// fmt.Println("\n✓ Receiver setup complete. Waiting for connection...")
	// fmt.Println("Press Ctrl+C to exit")

	// Keep the program running

	if err := pc.Connect("TEST"); err != nil {
		panic(err)
	}

	select {}
}
