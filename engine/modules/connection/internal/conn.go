package internal

import (
	"encoding/binary"
	"net"
)

type conn struct {
	*factory
	conn     net.Conn
	messages chan any
}

func (conn conn) Close() error       { return conn.conn.Close() }
func (conn conn) Messages() chan any { return conn.messages }

func (conn conn) Send(message any) error {
	bytes, err := conn.codec.Encode(message)
	if err != nil {
		return err
	}

	// conn.logger.Info(fmt.Sprintf("sending '***' of type '%v'", reflect.TypeOf(message).String()))
	length := uint32(len(bytes))
	lengthInByes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthInByes, length)
	if _, err := conn.conn.Write(append(lengthInByes, bytes...)); err != nil {
		return err
	}

	return nil
}
