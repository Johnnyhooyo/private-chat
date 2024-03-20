package common

import (
	"encoding/binary"
	"errors"
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

type DefaultPacker struct {
	encrypt Encrypt
}

func NewDefaultPacker(key []byte) Packer {
	return &DefaultPacker{
		encrypt: NewAesEncryptor(key),
	}
}

func (d *DefaultPacker) Pack(data []byte) ([]byte, error) {
	cipher := d.encrypt.Encrypt(data)
	dataSize := len(cipher)
	buffer := make([]byte, 4+dataSize)
	binary.LittleEndian.PutUint32(buffer[:4], uint32(dataSize)) // write dataSize
	copy(buffer[4:], cipher)                                    // write data
	return buffer, nil
}

func (d *DefaultPacker) Unpack(reader io.Reader) ([]byte, error) {
	headBuffer := make([]byte, 4)
	_, err := io.ReadFull(reader, headBuffer)
	if err != nil {
		return nil, err
	}
	dataSize := binary.LittleEndian.Uint32(headBuffer[:4])
	data := make([]byte, dataSize)
	if err = d.readData(reader, data, dataSize); err != nil {
		return nil, err
	}

	return d.encrypt.Decrypt(data), nil
}

// readData 处理tcp读取数据时的异常
func (d *DefaultPacker) readData(reader io.Reader, data []byte, dataSize uint32) error {
	if n, err := io.ReadFull(reader, data); err != nil {
		if uint32(n) == dataSize {
			fmt.Printf("read expected size data %d, and extra error is %s \n", n, err.Error())
			return nil
		}
		if errors.Is(err, io.EOF) {
			return err
		} else if errors.Is(err, io.ErrUnexpectedEOF) {
			// 如果数据还没读完，可以继续去读取 tcp保证消息顺序到达
			return d.readData(reader, data[n:], dataSize-uint32(n))
		}
		return fmt.Errorf("read data err: %s", err)
	}
	return nil
}
