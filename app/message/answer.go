package message

import (
	"bytes"
	"encoding/binary"
	"net"
)

type ResourceRecord struct {
	Name  string
	Type  uint16
	Class uint16
	TTL   uint32
	RData net.IP
}

func (rr *ResourceRecord) Pack() ([]byte, error) {
	var buf bytes.Buffer

	nameBytes, err := EncodeDomainName(rr.Name)
	if err != nil {
		return nil, err
	}

	buf.Write(nameBytes)

	binary.Write(&buf, binary.BigEndian, rr.Type)
	binary.Write(&buf, binary.BigEndian, rr.Class)
	binary.Write(&buf, binary.BigEndian, rr.TTL)
	binary.Write(&buf, binary.BigEndian, uint16(len(rr.RData)))
	buf.Write(rr.RData)

	return buf.Bytes(), nil
}
