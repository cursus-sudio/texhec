package netconnection

type MsgType int

const (
	MsgRequest MsgType = iota
	MsgResponse
	MsgMessage
	MsgError
)

type Msg struct {
	ID      string  `json:"id"`
	Payload *string `json:"payload"`
	Err     *string `json:"err"`
	Type    MsgType `json:"type"`
}

func NewRequest(id string, payload string) Msg {
	payloadHeap := payload
	return Msg{
		ID:      id,
		Payload: &payloadHeap,
		Err:     nil,
		Type:    MsgRequest,
	}
}

func NewResponse(id string, payload string, err error) Msg {
	payloadHeap := payload
	if err == nil {
		return Msg{
			ID:      id,
			Payload: &payloadHeap,
			Err:     nil,
			Type:    MsgResponse,
		}
	}
	errHeap := err.Error()
	return Msg{
		ID:      id,
		Payload: &payloadHeap,
		Err:     &errHeap,
		Type:    MsgResponse,
	}
}

func NewMessage(id string, payload string) Msg {
	payloadHeap := payload
	return Msg{
		ID:      id,
		Payload: &payloadHeap,
		Err:     nil,
		Type:    MsgMessage,
	}
}

func NewErrorMsg(id string, err error) Msg {
	errHeap := err.Error()
	return Msg{
		ID:      id,
		Payload: nil,
		Err:     &errHeap,
		Type:    MsgError,
	}
}
