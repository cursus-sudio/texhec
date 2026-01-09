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

func (f *factory) NewConnection(rawConn net.Conn) connection.Conn {
	messages := make(chan any)
	go func() {
		defer close(messages)

		for {
			messageLengthInBytes := make([]byte, 4)
			if _, err := io.ReadFull(rawConn, messageLengthInBytes); err != nil {
				break
			}
			messageLength := binary.BigEndian.Uint32(messageLengthInBytes)
			messageBytes := make([]byte, messageLength)
			if _, err := io.ReadFull(rawConn, messageBytes); err != nil {
				break
			}

			message, err := f.codec.Decode(messageBytes)
			if err != nil {
				f.logger.Warn(err)
				continue
			}
			// f.logger.Info(fmt.Sprintf("received '***' type '%v'", reflect.TypeOf(message).String()))

			messages <- message
		}

		_ = rawConn.Close()
	}()
	return &conn{
		factory:  f,
		conn:     rawConn,
		messages: messages,
	}
}
