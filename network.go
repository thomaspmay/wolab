package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// MagicPacket is a slice of 102 bytes containing the magic packet data.
type MagicPacket [102]byte

// createMagicPacket allocates a new MagicPacket with the specified MAC.
func createMagicPacket(macAddr string) (packet MagicPacket, err error) {
	mac, err := net.ParseMAC(macAddr)
	if err != nil {
		return packet, err
	}

	if len(mac) != 6 {
		return packet, errors.New("invalid EUI-48 MAC address")
	}

	// write magic bytes to packet
	copy(packet[0:], []byte{255, 255, 255, 255, 255, 255})
	offset := 6

	for i := 0; i < 16; i++ {
		copy(packet[offset:], mac)
		offset += 6
	}

	return packet, nil
}

func sendUDPPacket(mp MagicPacket, addr string) (err error) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			if err == nil {
				err = fmt.Errorf("failed to close connection: %w", closeErr)
			} else {
				err = fmt.Errorf("%v; also failed to close connection: %w", err, closeErr)
			}
		}
	}()

	_, err = conn.Write(mp[:])
	if err != nil {
		return fmt.Errorf("failed to write to UDP connection: %w", err)
	}
	return err
}

// Send writes the MagicPacket to the specified address on port 9.
func (mp MagicPacket) Send(addr string) error {
	return sendUDPPacket(mp, addr+":9")
}

// SendPort writes the MagicPacket to the specified address and port.
func (mp MagicPacket) SendPort(addr string, port string) error {
	return sendUDPPacket(mp, addr+":"+port)
}

type PingResult struct {
	IP    string
	Alive bool
}

func Ping(ip string, timeout time.Duration) bool {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listening for ICMP packets: %v\n", err)
		return false
	}
	defer conn.Close()

	message := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("PING"),
		},
	}

	messageBytes, err := message.Marshal(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling ICMP message: %v\n", err)
		return false
	}

	dst, err := net.ResolveIPAddr("ip4", ip)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving IP address: %v\n", err)
		return false
	}

	_, err = conn.WriteTo(messageBytes, dst)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending ICMP message: %v\n", err)
		return false
	}

	reply := make([]byte, 1500)
	err = conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting read deadline: %v\n", err)
		return false
	}

	_, _, err = conn.ReadFrom(reply)
	if err != nil {
		return false
	}

	return true
}

func PingNetwork(subnet string, timeout time.Duration) ([]PingResult, error) {
	ip, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return nil, fmt.Errorf("error parsing subnet: %v", err)
	}

	var wg sync.WaitGroup
	results := make([]PingResult, 0)

	mutex := &sync.Mutex{}

	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			alive := Ping(ip, timeout)
			mutex.Lock()
			results = append(results, PingResult{IP: ip, Alive: alive})
			mutex.Unlock()
		}(ip.String())
	}

	wg.Wait()
	return results, nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// todo: Create function which hosts a fake machine to test discovery and WoL
