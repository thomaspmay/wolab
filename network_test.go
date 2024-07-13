package main

import (
	"bytes"
	"net"
	"testing"
	"time"
)

const (
	MagicPacketSize = 102
)

func TestNewMagicPacket(t *testing.T) {
	packet, err := createMagicPacket("AA:BB:CC:DD:EE:FF")
	if err != nil {
		t.Fatalf("magic packet creation failed: %v", err)
	}

	// Check for magic padding
	if !bytes.Equal(packet[:6], []byte{255, 255, 255, 255, 255, 255}) {
		t.Fatal("packet doesn't contain 6 bytes of 0xFF padding")
	}

	// AA:BB:CC:DD:EE:FF
	rawMac := []byte{170, 187, 204, 221, 238, 255}
	pos := 6
	for i := 0; i < 16; i++ {
		if !bytes.Equal(packet[pos:pos+6], rawMac) {
			t.Fatalf("magic packet contains incorrect MAC address at iteration %d", i)
		}
		pos += 6
	}
}

func TestEui48(t *testing.T) {
	// only support EUI-48 addresses, EUI-64 should fail
	_, err := createMagicPacket("AA:BB:CC:AA:BB:CC:AA:BB")
	if err == nil {
		t.Error("able to construct magic packet with invalid MAC")
	}

	_, err = createMagicPacket("AA:BB:CC:AA:BB:CC")
	if err != nil {
		t.Error(err)
	}
}

func TestPing(t *testing.T) {
	// Test a known live IP address (127.0.0.1 is localhost and should always be reachable)
	ip := "127.0.0.1"
	timeout := 1 * time.Second

	alive := Ping(ip, timeout)
	if !alive {
		t.Errorf("Expected %s to be alive, but it is not responding", ip)
	}

	// Test a known unreachable IP address (unlikely to be in use)
	ip = "192.0.2.0" // Reserved IP for documentation purposes, typically not in use
	alive = Ping(ip, timeout)
	if alive {
		t.Errorf("Expected %s to be unreachable, but it is responding", ip)
	}
}

func TestPingNetwork(t *testing.T) {
	// Test a small local subnet
	subnet := "127.0.0.0/30"
	timeout := 1 * time.Second

	results, err := PingNetwork(subnet, timeout)
	if err != nil {
		t.Fatalf("Error pinging network: %v", err)
	}

	expectedResults := map[string]bool{
		"127.0.0.0": false,
		"127.0.0.1": true,
		"127.0.0.2": true,
		"127.0.0.3": true,
	}

	if len(results) != len(expectedResults) {
		t.Fatalf("Expected %d results, but got %d", len(expectedResults), len(results))
	}

	for _, result := range results {
		expectedAlive, exists := expectedResults[result.IP]
		if !exists {
			t.Errorf("Unexpected IP address in results: %s", result.IP)
		} else if result.Alive != expectedAlive {
			t.Errorf("Expected IP %s to be alive: %v, but got: %v", result.IP, expectedAlive, result.Alive)
		}
	}
}

// Test IP increment function
func TestInc(t *testing.T) {
	ip := net.ParseIP("192.168.1.1").To4()
	inc(ip)
	if ip.String() != "192.168.1.2" {
		t.Errorf("Expected IP to increment to 192.168.1.2, but got %s", ip.String())
	}

	ip = net.ParseIP("192.168.1.255").To4()
	inc(ip)
	if ip.String() != "192.168.2.0" {
		t.Errorf("Expected IP to increment to 192.168.2.0, but got %s", ip.String())
	}
}
