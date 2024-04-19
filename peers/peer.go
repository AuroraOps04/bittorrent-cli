package peers

import (
	"encoding/binary"
	"errors"
	"net"
	"strconv"
)

type TrackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

type Peer struct {
	IP   net.IP
	Port uint16
}

// Unmarshal parse the byte list to peer list
// 6 byte for one peer
// 4 byte for IP
// 2 byte for Port
func Unmarshal(peers []byte) ([]Peer, error) {
	const peerSize = 6
	if len(peers)%peerSize != 0 {
		return nil, errors.New("invalid peers")
	}
	p := make([]Peer, len(peers)/peerSize)
	for i := 0; i < len(p); i++ {
		offset := i * peerSize
		p[i].IP = net.IP(peers[offset : offset+4])
		p[i].Port = binary.BigEndian.Uint16(peers[offset+4 : offset+6])
	}
	return p, nil

}
func (p Peer) String() string {
	if p.IP == nil {
		return ""
	}
	return net.JoinHostPort(p.IP.String(), strconv.FormatUint(uint64(p.Port), 10))
}

func GetLocalPeerID() [20]byte {
	return [20]byte{'-', 'A', 'u', 'r', 'o', 'r', 'a', 'O', 'p', 's', '2', 'O', 'p', 's', '2', 'O', 'p', 's', '2', '-'}
}
