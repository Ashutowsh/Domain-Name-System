package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/Ashutowwsh/dns-server-go/cache"
	"github.com/Ashutowwsh/dns-server-go/db"
	"github.com/Ashutowwsh/dns-server-go/message"
)

type DNSHandler struct {
	RedisClient *cache.RedisClient
	PostgresDB  *db.PostgresDB
}

func (h *DNSHandler) HandleDNSPacket(packet []byte, clientIP string) ([]byte, error) {
	ctx := context.Background()

	// Parse DNS header
	header := &message.Header{}
	err := header.Unpack(packet[:12])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack header: %w", err)
	}

	// Parse DNS question
	offset := 12
	question := &message.Question{}
	_, err = question.Unpack(packet, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack question: %w", err)
	}

	// Rate limiting
	allowed, err := h.RedisClient.RateLimit(ctx, clientIP, 5, time.Minute)
	if err != nil {
		return nil, fmt.Errorf("rate limit error: %w", err)
	}
	if !allowed {
		return nil, fmt.Errorf("rate limit exceeded for client: %s", clientIP)
	}

	// Check Redis cache
	cacheKey := fmt.Sprintf("%s:%d", question.Name, question.Type)
	cachedResponse, err := h.RedisClient.GetCache(ctx, cacheKey)
	if err != nil {
		return nil, err
	}
	if cachedResponse != "" {
		return []byte(cachedResponse), nil // Return cached response
	}

	// Create DNS response
	responseHeader := *header
	responseHeader.QR = 1 // Response flag
	responseHeader.AA = 1 // Authoritative answer
	responseHeader.ANCount = 1

	answer := message.ResourceRecord{
		Name:  question.Name,
		Type:  1, // A record
		Class: 1, // IN class
		TTL:   300,
		RData: net.ParseIP("93.184.216.34").To4(),
	}

	headerBytes, _ := responseHeader.Pack()
	questionBytes, _ := question.Pack()
	answerBytes, _ := answer.Pack()

	response := append(append(headerBytes, questionBytes...), answerBytes...)

	// Cache the response
	responseJSON, _ := json.Marshal(response)
	_ = h.RedisClient.SetCache(ctx, cacheKey, string(responseJSON), 5*time.Minute)

	// Log the query
	_ = h.PostgresDB.LogQuery(ctx, question.Name, "A", "93.184.216.34")

	return response, nil
}
