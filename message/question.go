package message

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Question struct {
	Name  string
	Type  uint16
	Class uint16
}

func (q *Question) Unpack(data []byte, offset int) (int, error) {
	var err error
	q.Name, offset, err = ParseName(data, offset)
	if err != nil {
		return offset, err
	}

	if offset+4 > len(data) {
		return offset, fmt.Errorf("insufficient data for question")
	}

	q.Type = binary.BigEndian.Uint16(data[offset : offset+2])
	q.Class = binary.BigEndian.Uint16(data[offset+2 : offset+4])
	offset += 4

	return offset, nil
}

func (q *Question) Pack() ([]byte, error) {
	var buf bytes.Buffer

	nameBytes, err := EncodeDomainName(q.Name)
	if err != nil {
		return nil, err
	}
	buf.Write(nameBytes)

	binary.Write(&buf, binary.BigEndian, q.Type)
	binary.Write(&buf, binary.BigEndian, q.Class)

	return buf.Bytes(), nil
}
