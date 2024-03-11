package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Packer is a generic interface to pack and unpack message packet.
type Packer interface {
	// Pack packs Message into the packet to be written.
	Pack(data []byte) ([]byte, error)

	// Unpack unpacks the message packet from reader,
	// returns the message, and error if error occurred.
	Unpack(reader io.Reader) ([]byte, error)
}

type defaultPacker struct {
}

func newDefaultPacker() *defaultPacker {
	return &defaultPacker{}
}

func (d *defaultPacker) Pack(data []byte) ([]byte, error) {
	dataSize := len(data)
	buffer := make([]byte, 4+dataSize)
	binary.LittleEndian.PutUint32(buffer[:4], uint32(dataSize)) // write dataSize
	copy(buffer[4:], data)                                      // write data
	return buffer, nil
}

func (d *defaultPacker) Unpack(reader io.Reader) ([]byte, error) {
	headBuffer := make([]byte, 8)
	_, err := io.ReadFull(reader, headBuffer)
	if err != nil {
		return nil, err
	}
	dataSize := binary.LittleEndian.Uint32(headBuffer[:4])
	data := make([]byte, dataSize)
	if _, err := io.ReadFull(reader, data); err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, fmt.Errorf("read data err: %s", err)
	}
	return data, nil
}
