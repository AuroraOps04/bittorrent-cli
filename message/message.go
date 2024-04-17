package message

import (
	"encoding/binary"
	"github.com/pkg/errors"
	"io"
)

type messageID uint8

const (
	MsgChoke         messageID = 0
	MsgUnchoke       messageID = 1
	MsgInterested    messageID = 2
	MsgNotInterested messageID = 3
	MsgHave          messageID = 4
	MsgBitfield      messageID = 5
	MsgRequest       messageID = 6
	MsgPiece         messageID = 7
	MsgCancel        messageID = 8
)

type Message struct {
	ID      messageID
	Payload []byte
}

func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}
	length := uint32(len(m.Payload) + 1) // +1 for id
	buf := make([]byte, 4+length)        // 4 for length
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload)
	return buf
}

func Read(r io.Reader) (*Message, error) {
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	length := binary.BigEndian.Uint32(lengthBuf)
	// keep-alive message
	if length == 0 {
		return nil, nil
	}
	messageBuf := make([]byte, length)
	_, err = io.ReadFull(r, messageBuf)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	m := Message{
		ID:      messageID(messageBuf[0]),
		Payload: messageBuf[1:],
	}
	return &m, nil
}

// FormatRequest create a REUQEST message
func FormatRequest(index, begin, length int) *Message {
	m := Message{
		ID: MsgRequest,
	}
	buf := make([]byte, 12)
	binary.BigEndian.PutUint32(buf[0:4], uint32(index))
	binary.BigEndian.PutUint32(buf[4:8], uint32(begin))
	binary.BigEndian.PutUint32(buf[8:], uint32(length))
	m.Payload = buf
	return &m
}

// FormatHave create a HAVE message
func FormatHave(index int) *Message {
	m := Message{
		ID: MsgHave,
	}
	buf := make([]byte, 4)
	binary.BigEndian.AppendUint32(buf[:], uint32(index))
	m.Payload = buf
	return &m
}
