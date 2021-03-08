package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/panjf2000/gnet"
)

/**
*自定义编码解码器
@author zhangchengji
*/
const (
	DefaultHeadLength = 2
)

type LogLengthFieldProtocol struct {
}

func (cc *LogLengthFieldProtocol) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
	panic("implement me")
}

// Decode ...
func (cc *LogLengthFieldProtocol) Decode(c gnet.Conn) ([]byte, error) {
	// parse header
	headerLen := DefaultHeadLength // uint16+uint16+uint32
	if size, header := c.ReadN(headerLen); size == headerLen {
		byteBuffer := bytes.NewBuffer(header)
		var dataLength uint16
		_ = binary.Read(byteBuffer, binary.BigEndian, &dataLength)

		// parse payload
		dataLen := int(dataLength) //max int32 can contain 210MB payload
		protocolLen := headerLen + dataLen
		if dataSize, data := c.ReadN(protocolLen); dataSize == protocolLen {
			c.ShiftN(protocolLen)
			//log.Println("parse success:", data, dataSize)

			// return the payload of the data
			return data[headerLen:], nil
		}
		// log.Println("not enough payload data:", dataLen, protocolLen, dataSize)
		return nil, errors.New("not enough payload data")

	}
	// log.Println("not enough header data:", size)
	return nil, errors.New("not enough header data")
}
