package message

import (
	"bytes"
	"fmt"
)

// EncodeDomainName encodes a domain name into the label format used in DNS packets.
func EncodeDomainName(domain string) ([]byte, error) {
	labels := bytes.Split([]byte(domain), []byte("."))
	var buf bytes.Buffer

	for _, label := range labels {
		if len(label) > 63 {
			return nil, fmt.Errorf("label length exceeds 63 characters")
		}

		buf.WriteByte(byte(len(label)))
		buf.Write(label)
	}

	buf.WriteByte(0) // Null byte to terminate the domain
	return buf.Bytes(), nil
}

// ParseName decodes a label-encoded domain name from a DNS packet.
func ParseName(data []byte, offset int) (string, int, error) {
	var domain []byte

	for {
		if offset >= len(data) {
			return "", offset, fmt.Errorf("offset out of range")
		}

		length := int(data[offset])
		offset++

		if length == 0 { // Null byte signals end of the domain name
			break
		}

		if offset+length > len(data) {
			return "", offset, fmt.Errorf("invalid label length")
		}

		domain = append(domain, data[offset:offset+length]...)
		domain = append(domain, '.') // Add a dot between labels
		offset += length
	}

	if len(domain) > 0 && domain[len(domain)-1] == '.' {
		domain = domain[:len(domain)-1] // Remove the trailing dot
	}

	return string(domain), offset, nil
}
