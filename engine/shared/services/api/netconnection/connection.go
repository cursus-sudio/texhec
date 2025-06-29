package netconnection

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"shared/services/codec"
	"shared/services/uuid"
	"shared/utils/connection"
	"shared/utils/httperrors"
	"sync"

	"github.com/ogiusek/relay/v2"
)

var (
	ErrInvalidMessageType error = errors.Join(httperrors.Err400, errors.New("invalid message content"))
)

type NetConnection interface {
	Connect(conn net.Conn, onClose func()) connection.Connection
}

type netConnection struct {
	Codec           codec.Codec
	Connection      connection.Connection
	UuidFactory     uuid.Factory
	PendingMutex    sync.Mutex
	PendingRequests map[string]chan Msg
}

func newNetConnection(
	codec codec.Codec,
	connection connection.Connection,
	uuidFactory uuid.Factory,
) NetConnection {
	return &netConnection{
		Codec:           codec,
		Connection:      connection,
		UuidFactory:     uuidFactory,
		PendingMutex:    sync.Mutex{},
		PendingRequests: map[string]chan Msg{},
	}
}

func (c *netConnection) Request(conn net.Conn, msg Msg) Msg {
	c.PendingMutex.Lock()
	msgChan := make(chan Msg)
	c.PendingRequests[string(msg.ID)] = msgChan
	c.PendingMutex.Unlock()

	c.Send(conn, msg)

	return <-msgChan
}
func (c *netConnection) Respond(conn net.Conn, msg Msg) { c.Send(conn, msg) }
func (c *netConnection) Message(conn net.Conn, msg Msg) { c.Send(conn, msg) }

//

func (c *netConnection) HandleMsg(conn net.Conn, msg Msg) error {
	switch msg.Type {
	case MsgRequest:
		if msg.Payload == nil {
			return httperrors.Err400
		}
		decoded, err := c.Codec.Decode([]byte(*msg.Payload))
		if err != nil {
			return httperrors.Err400
		}
		r := c.Connection.Relay()
		res, err := relay.HandleAny(r, decoded)
		resMsg := NewResponse(msg.ID, string(c.Codec.Encode(res)), err)
		c.Send(conn, resMsg)
		break
	case MsgResponse:
		c.PendingMutex.Lock()
		defer c.PendingMutex.Unlock()
		id := msg.ID
		msgChan, ok := c.PendingRequests[id]
		if !ok {
			return httperrors.Err404
		}
		msgChan <- msg
		break
	case MsgMessage:
		if msg.Payload == nil {
			return httperrors.Err400
		}
		decoded, err := c.Codec.Decode([]byte(*msg.Payload))
		if err != nil {
			return httperrors.Err400
		}
		r := c.Connection.Relay()
		relay.HandleAnyMessage(r, decoded)
		break
	default:
		return ErrInvalidMessageType
	}
	return nil
}

func (c *netConnection) Send(conn net.Conn, msg Msg) error {
	bytes := c.Codec.Encode(msg)
	length := uint64(len(bytes))
	lengthInByes := make([]byte, 8)
	binary.BigEndian.PutUint64(lengthInByes, uint64(length))
	n, err := conn.Write(append(lengthInByes, bytes...)) // n, err := conn.Write(append(lengthInByes, bytes...))
	if err != nil {
		return err
	}
	if uint64(n) != length+8 {
		return fmt.Errorf("sent message length is not of the same length as message payload")
	}
	return nil
}

func (c *netConnection) Listen(conn net.Conn, onClose func()) {
	for {
		messageLengthInBytes := make([]byte, 8)
		if read, err := io.ReadFull(conn, messageLengthInBytes); err != nil || read != 8 {
			break
		}
		messageLength := binary.BigEndian.Uint64(messageLengthInBytes)
		messageBytes := make([]byte, messageLength)
		if read, err := io.ReadFull(conn, messageBytes); err != nil || read != int(messageLength) {
			break
		}

		msgAny, err := c.Codec.Decode(messageBytes)
		if err != nil {
			c.Send(
				conn,
				NewErrorMsg(
					c.UuidFactory.NewUUID().String(),
					err,
				),
			)
			break
		}
		msg, ok := msgAny.(Msg)
		if !ok {
			c.Send(
				conn,
				NewErrorMsg(
					c.UuidFactory.NewUUID().String(),
					ErrInvalidMessageType,
				),
			)
			break
		}

		go c.HandleMsg(conn, msg)
	}
	conn.Close()
	onClose()
}

func (c *netConnection) Connect(conn net.Conn, onClose func()) connection.Connection {
	cb := connection.NewCloseBuilder()
	cb.OnClose(func() {
		conn.Close()
	})

	mlb := connection.NewMessageListenerBuilder()

	mlb.Relay().DefaultHandler(func(ctx relay.AnyContext) {
		req := ctx.Req()
		reqMsg := NewRequest(
			c.UuidFactory.NewUUID().String(),
			string(c.Codec.Encode(req)),
		)
		resMsg := c.Request(conn, reqMsg)
		if resMsg.Err != nil {
			ctx.SetErr(errors.New(*resMsg.Err))
			return
		}
		if resMsg.Payload == nil {
			ctx.SetErr(httperrors.Err500)
			return
		}
		res, err := c.Codec.Decode([]byte(*resMsg.Payload))
		if err != nil {
			ctx.SetErr(httperrors.Err500)
			return
		}
		ctx.SetRes(res)
	})

	mlb.Relay().DefaultMessageHandler(func(ctx relay.AnyMessageCtx) {
		msg := NewMessage(
			c.UuidFactory.NewUUID().String(),
			string(c.Codec.Encode(ctx.Message())),
		)
		c.Send(conn, msg)
	})

	go c.Listen(conn, onClose)

	return connection.NewConnection(
		cb.Build(),
		mlb.Build(),
	)
}
