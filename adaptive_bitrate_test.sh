#!/bin/bash

# WebRTC Bitrate Estimation Test Script
# Network Interface: veth0 (virtual network namespace)

INTERFACE="wlx202351285376"
LOG_FILE="webrtc_test_$(date +%Y%m%d_%H%M%S).log"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log() {
    echo -e "$(date '+%H:%M:%S') - $1" | tee -a "$LOG_FILE"
}

cleanup() {
    log "${RED}Cleaning up network rules...${NC}"
    sudo tc qdisc del dev "$INTERFACE" root 2>/dev/null || true
    log "${GREEN}Cleanup complete${NC}"
    exit 0
}

# Trap cleanup on script exit
trap cleanup EXIT INT TERM

wait_for_input() {
    echo -e "${YELLOW}Press Enter to continue to next test (or Ctrl+C to exit)...${NC}"
    read -r
}

apply_tc_rule() {
    local rule="$1"
    local description="$2"

    log "${BLUE}Applying: $description${NC}"
    log "Command: sudo tc qdisc change dev $INTERFACE root $rule"

    # Try change first, if it fails, add new rule
    if ! sudo tc qdisc change dev "$INTERFACE" root $rule 2>/dev/null; then
        sudo tc qdisc del dev "$INTERFACE" root 2>/dev/null || true
        sudo tc qdisc add dev "$INTERFACE" root $rule
    fi

    log "${GREEN}Applied: $description${NC}"

    # Show what was actually applied
    log "Actual rule: $(sudo tc qdisc show dev $INTERFACE | grep root)"
}

echo "================================================"
echo "WebRTC Bitrate Estimation Test Suite"
echo "Interface: $INTERFACE"
echo "Log file: $LOG_FILE"
echo "================================================"

# Verify interface exists
if ! ip link show "$INTERFACE" >/dev/null 2>&1; then
    echo "Error: Interface $INTERFACE not found. Make sure to run 'make setup-test-network' first."
    exit 1
fi

# Test 1: Baseline - No restrictions
log "${GREEN}=== TEST 1: BASELINE (No restrictions) ===${NC}"
sudo tc qdisc del dev "$INTERFACE" root 2>/dev/null || true
log "Network: Clean slate - no artificial constraints"
log "Expected: Should reach MaximumBitrate (2Mbps)"
wait_for_input

# Test 2: Static Bandwidth Limits
log "${GREEN}=== TEST 2: STATIC BANDWIDTH LIMITS ===${NC}"

log "${BLUE}Test 2.1: 1.5Mbps limit${NC}"
apply_tc_rule "netem rate 1500kbit" "1.5Mbps bandwidth limit"
log "Expected: Stabilize ~1.3-1.4Mbps"
wait_for_input

log "${BLUE}Test 2.2: 1Mbps limit${NC}"
apply_tc_rule "netem rate 1mbit" "1Mbps bandwidth limit"
log "Expected: Stabilize ~900kbps"
wait_for_input

log "${BLUE}Test 2.3: 600kbps limit${NC}"
apply_tc_rule "netem rate 600kbit" "600kbps bandwidth limit"
log "Expected: Stabilize ~550kbps"
wait_for_input

log "${BLUE}Test 2.4: 400kbps limit (below minimum)${NC}"
apply_tc_rule "netem rate 400kbit" "400kbps bandwidth limit"
log "Expected: Hit MinimumBitrate and pause video"
wait_for_input

# Test 3: Dynamic Bandwidth Changes
log "${GREEN}=== TEST 3: DYNAMIC BANDWIDTH ADAPTATION ===${NC}"

log "${BLUE}Test 3.1: Degradation Sequence${NC}"
log "Starting with 2Mbps..."
apply_tc_rule "netem rate 2mbit" "2Mbps start"
sleep 30

log "Degrading to 1Mbps..."
apply_tc_rule "netem rate 1mbit" "Degrade to 1Mbps"
sleep 30

log "Degrading to 600kbps..."
apply_tc_rule "netem rate 600kbit" "Degrade to 600kbps"
sleep 30

log "Degrading to 400kbps (should pause)..."
apply_tc_rule "netem rate 400kbit" "Degrade to 400kbps"
log "Degradation sequence complete"
wait_for_input

log "${BLUE}Test 3.2: Recovery Sequence${NC}"
log "Starting with 400kbps..."
apply_tc_rule "netem rate 400kbit" "400kbps start"
sleep 30

log "Improving to 800kbps..."
apply_tc_rule "netem rate 800kbit" "Improve to 800kbps"
sleep 30

log "Improving to 1.5Mbps..."
apply_tc_rule "netem rate 1500kbit" "Improve to 1.5Mbps"
sleep 30

log "Removing all limits..."
sudo tc qdisc del dev "$INTERFACE" root 2>/dev/null || true
log "Recovery sequence complete"
wait_for_input

# Test 4: Packet Loss Tests
log "${GREEN}=== TEST 4: PACKET LOSS IMPACT ===${NC}"

log "${BLUE}Test 4.1: 1% packet loss${NC}"
apply_tc_rule "netem rate 1mbit loss 1%" "1Mbps + 1% packet loss"
log "Expected: Slight bitrate reduction"
wait_for_input

log "${BLUE}Test 4.2: 2% packet loss${NC}"
apply_tc_rule "netem rate 1mbit loss 2%" "1Mbps + 2% packet loss"
log "Expected: More aggressive bitrate reduction"
wait_for_input

log "${BLUE}Test 4.3: 5% burst loss${NC}"
apply_tc_rule "netem rate 1mbit loss 5% 25%" "1Mbps + 5% burst loss"
log "Expected: Significant bitrate reduction due to NACK storms"
wait_for_input

# Test 5: Latency and Jitter
log "${GREEN}=== TEST 5: LATENCY AND JITTER TESTS ===${NC}"

log "${BLUE}Test 5.1: High latency (200ms RTT)${NC}"
apply_tc_rule "netem delay 100ms rate 1mbit" "100ms delay (200ms RTT) + 1Mbps"
log "Expected: Slower adaptation due to feedback delay"
wait_for_input

log "${BLUE}Test 5.2: High jitter${NC}"
apply_tc_rule "netem delay 100ms 50ms rate 1mbit" "100ms ±50ms jitter + 1Mbps"
log "Expected: Less stable bitrate due to variable feedback timing"
wait_for_input

# Test 6: Real-world Simulations
log "${GREEN}=== TEST 6: REAL-WORLD SIMULATIONS ===${NC}"

log "${BLUE}Test 6.1: 4G Mobile Network${NC}"
apply_tc_rule "netem delay 80ms 20ms loss 1% rate 2mbit" "4G simulation: 80±20ms delay, 1% loss, 2Mbps"
log "Expected: Moderate adaptation with occasional drops"
wait_for_input

log "${BLUE}Test 6.2: Poor WiFi${NC}"
apply_tc_rule "netem delay 150ms 50ms loss 3% 50% rate 800kbit" "Poor WiFi: high jitter, burst loss, 800kbps"
log "Expected: Aggressive bitrate reduction and instability"
wait_for_input

log "${BLUE}Test 6.3: Congested Network${NC}"
sudo tc qdisc del dev "$INTERFACE" root 2>/dev/null || true
sudo tc qdisc add dev "$INTERFACE" root handle 1: netem delay 120ms 30ms
# Use netem rate instead of tbf for consistency
sudo tc qdisc change dev "$INTERFACE" root handle 1: netem delay 120ms 30ms rate 600kbit
log "Applied: Congested network with variable delay and limited bandwidth"
log "Expected: Conservative bitrate with high variability"
wait_for_input

# Oscillating test
log "${GREEN}=== TEST 7: OSCILLATING CONDITIONS ===${NC}"
log "${BLUE}Test 7.1: Bandwidth oscillation (15s intervals)${NC}"
for i in {1..6}; do
    if [ $((i % 2)) -eq 1 ]; then
        apply_tc_rule "netem rate 1mbit" "Oscillation: 1Mbps"
    else
        apply_tc_rule "netem rate 500kbit" "Oscillation: 500kbps"
    fi
    sleep 15
done
log "Oscillation test complete"

# Test 8: Extreme Conditions
log "${GREEN}=== TEST 8: EXTREME CONDITIONS ===${NC}"

log "${BLUE}Test 8.1: Very high loss${NC}"
apply_tc_rule "netem rate 1mbit loss 10%" "1Mbps + 10% packet loss"
log "Expected: Significant bitrate reduction or connection drops"
wait_for_input

log "${BLUE}Test 8.2: Very high latency${NC}"
apply_tc_rule "netem delay 500ms rate 1mbit" "500ms delay + 1Mbps"
log "Expected: Very slow adaptation, potential timeouts"
wait_for_input

log "${BLUE}Test 8.3: Very low bandwidth${NC}"
apply_tc_rule "netem rate 200kbit" "200kbps (very low)"
log "Expected: Video should pause or drop to minimum quality"
wait_for_input

log "${GREEN}=== ALL TESTS COMPLETE ===${NC}"
log "Check your WebRTC logs and compare with network conditions"
log "Test results logged to: $LOG_FILE"
echo ""
echo "Summary of test results can be found in: $LOG_FILE"
echo "Make sure to analyze your WebRTC application logs for bitrate adaptation patterns!"

# Final cleanup will happen automatically via trap