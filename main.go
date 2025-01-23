package main

import (
	"fmt"
	"net"

	"github.com/Ashutowwsh/dns-server-go/cache"
	"github.com/Ashutowwsh/dns-server-go/db"
	"github.com/Ashutowwsh/dns-server-go/handlers"
)

func main() {
	fmt.Println("Starting DNS server on 127.0.0.1:2053")

	redisClient := cache.NewRedisClient()
	postgresDB := db.NewPostgresDB("postgres://user:password@localhost:5432/dnsdb")

	dnsHandler := &handlers.DNSHandler{
		RedisClient: redisClient,
		PostgresDB:  postgresDB,
	}

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

	buf := make([]byte, 512)
	for {
		n, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			continue
		}

		response, err := dnsHandler.HandleDNSPacket(buf[:n], source.IP.String())
		if err != nil {
			fmt.Println("Error handling DNS packet:", err)
			continue
		}

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
