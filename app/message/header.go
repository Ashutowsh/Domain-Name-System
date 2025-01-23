package message

import (
	"encoding/binary"
	"fmt"
)

type Header struct {
	ID      uint16
	QR      uint8
	Opcode  uint8
	AA      uint8
	TC      uint8
	RD      uint8
	RA      uint8
	Z       uint8
	Rcode   uint8
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

func (h *Header) Pack() ([]byte, error) {
	buf := make([]byte, 12)

	binary.BigEndian.PutUint16(buf[0:2], h.ID)

	flags := uint16(h.QR)<<15 |
		uint16(h.Opcode)<<11 |
		uint16(h.AA)<<10 |
		uint16(h.TC)<<9 |
		uint16(h.RD)<<8 |
		uint16(h.RA)<<7 |
		uint16(h.Z)<<4 |
		uint16(h.Rcode)

	binary.BigEndian.PutUint16(buf[2:4], flags)

	binary.BigEndian.PutUint16(buf[4:6], h.QDCount)
	binary.BigEndian.PutUint16(buf[6:8], h.ANCount)
	binary.BigEndian.PutUint16(buf[8:10], h.NSCount)
	binary.BigEndian.PutUint16(buf[10:12], h.ARCount)

	return buf, nil
}

func (h *Header) Unpack(data []byte) error {
	if len(data) < 12 {
		return fmt.Errorf("invalid header length")
	}

	h.ID = binary.BigEndian.Uint16(data[0:2])

	flags := binary.BigEndian.Uint16(data[2:4])
	h.QR = uint8(flags >> 15)
	h.Opcode = uint8((flags >> 11) & 0xF)
	h.AA = uint8((flags >> 10) & 0x1)
	h.TC = uint8((flags >> 9) & 0x1)
	h.RD = uint8((flags >> 8) & 0x1)
	h.RA = uint8((flags >> 7) & 0x1)
	h.Z = uint8((flags >> 4) & 0x7)
	h.Rcode = uint8(flags & 0xF)

	h.QDCount = binary.BigEndian.Uint16(data[4:6])
	h.ANCount = binary.BigEndian.Uint16(data[6:8])
	h.NSCount = binary.BigEndian.Uint16(data[8:10])
	h.ARCount = binary.BigEndian.Uint16(data[10:12])

	return nil
}
