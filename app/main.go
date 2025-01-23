package main

import (
	"fmt"
	"net"

	"github.com/Ashutowwsh/dns-server-go/app/message"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()
	fmt.Println("UDP Connection : ", udpConn)
	fmt.Println("[Info] DNS server is running on 127.0.0.1:2053")

	buf := make([]byte, 512)
	fmt.Println("Buffer created.")

	for {
		n, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			continue
		}
		fmt.Printf("[Info] Received DNS query from %v\n", source)

		response, err := handleDNSPacket(buf[:n])
		if err != nil {
			fmt.Println("Error handling DNS packet:", err)
			continue
		}

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
		fmt.Printf("[Info] Sent DNS response to %v\n", source)
	}
}

func handleDNSPacket(packet []byte) ([]byte, error) {

	// Parse the header
	header := &message.Header{}
	err := header.Unpack(packet[:12])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack header: %w", err)
	}

	// Parse the question
	offset := 12
	question := &message.Question{}
	_, err = question.Unpack(packet, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack question: %w", err)
	}

	// response header
	responseHeader := *header
	responseHeader.QR = 1 // Response flag
	responseHeader.AA = 1 // Authoritative answer
	responseHeader.ANCount = 1

	// answer section
	answer := message.ResourceRecord{
		Name:  question.Name,
		Type:  1, // A record
		Class: 1, // IN class
		TTL:   300,
		RData: net.ParseIP("93.184.216.34").To4(),
	}

	// Pack all components into the response
	headerBytes, err := responseHeader.Pack()
	if err != nil {
		return nil, fmt.Errorf("failed to pack header: %w", err)
	}

	questionBytes, err := question.Pack()
	if err != nil {
		return nil, fmt.Errorf("failed to pack question: %w", err)
	}

	answerBytes, err := answer.Pack()
	if err != nil {
		return nil, fmt.Errorf("failed to pack answer: %w", err)
	}

	return append(append(headerBytes, questionBytes...), answerBytes...), nil
}
