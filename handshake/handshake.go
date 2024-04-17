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
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if len(buf) == 0{
		return nil, errors.New("invalid handshake")
	}
	pstrLen := buf[0]
	if len(buf) != int(pstrLen)+49 {
		return nil, errors.New("invalid handshake")
	}
	h := Handshake{}
	h.Pstr = string(buf[1 : 1+pstrLen])
	// 8 is reserved bytes
	h.InfoHash = [20]byte(buf[pstrLen+9 : pstrLen+29])
	h.PeerID = [20]byte(buf[pstrLen+29:])
	return &h, nil

}
