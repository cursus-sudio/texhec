package internal

import (
	"encoding/binary"
	"engine/modules/connection"
	"engine/services/codec"
	"engine/services/logger"
	"io"
	"net"
)

type factory struct {
	codec  codec.Codec
	logger logger.Logger
}

func NewFactory(codec codec.Codec, logger logger.Logger) *factory {
	return &factory{codec, logger}
}

func (f *factory) NewConnection(rawConn net.Conn) connection.Connection {
	messages := make(chan any)
	go func() {
		defer rawConn.Close()
		defer close(messages)

		for {
			messageLengthInBytes := make([]byte, 2)
			if _, err := io.ReadFull(rawConn, messageLengthInBytes); err != nil {
				break
			}
			messageLength := binary.BigEndian.Uint16(messageLengthInBytes)
			messageBytes := make([]byte, messageLength)
			if _, err := io.ReadFull(rawConn, messageBytes); err != nil {
				break
			}

			message, err := f.codec.Decode(messageBytes)
			if err != nil {
				f.logger.Warn(err)
				continue
			}

			messages <- message
		}
	}()
	return conn{
		factory:  f,
		conn:     rawConn,
		messages: messages,
	}
}
