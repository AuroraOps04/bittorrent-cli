package handshake

import (
	"io"

	"github.com/pkg/errors"
)

type Handshake struct {
	Pstr     string
	InfoHash [20]byte
	PeerID   [20]byte
}

func (h *Handshake) Serialize() []byte {
	// 49 = 1 + 8 + 20 + 20
	buf := make([]byte, len(h.Pstr)+49)
	buf[0] = byte(len(h.Pstr))
	cur := 1
	cur += copy(buf[cur:], h.Pstr)
	// 8 byte of zero
	cur += copy(buf[cur:], make([]byte, 8))
	cur += copy(buf[cur:], h.InfoHash[:])
	cur += copy(buf[cur:], h.PeerID[:])
	return buf
}

func Read(r io.Reader) (*Handshake, error) {
	lengthBuf := make([]byte, 1)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read pstr length")
	}
	length := int(lengthBuf[0])
	anotherBuf := make([]byte, 48+length)
	_, err = io.ReadFull(r, anotherBuf)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read the rest of the handshake")
	}
	pstr := string(anotherBuf[0:length])
	infohash := anotherBuf[length+8 : length+28]
	peerid := anotherBuf[length+28 : length+48]
	h := Handshake{
		Pstr:     pstr,
		InfoHash: [20]byte(infohash),
		PeerID:   [20]byte(peerid),
	}
	return &h, nil
}
