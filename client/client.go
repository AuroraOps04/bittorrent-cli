package client

import (
	"io"
	"net"

	"github.com/AuroraOps04/bittorrent-cli/bitfield"
	"github.com/AuroraOps04/bittorrent-cli/handshake"
	"github.com/AuroraOps04/bittorrent-cli/peer"
	"github.com/pkg/errors"
)

type Client struct {
	Conn     net.Conn
	Bitfield bitfield.Bitfield
	PeerID   [20]byte
	InfoHash [20]byte
	Peer     peer.Peer
}

func New(p peer.Peer, id, infoHash [20]byte) (*Client, error) {
	c := &Client{
		PeerID:   id,
		Peer:     p,
		InfoHash: infoHash,
	}
	conn, err := net.Dial("tcp", p.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c.Conn = conn
	// handShake
	err = completeHandshake(conn, id, infoHash)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	bt, err := recvBitfield(conn)
	if err != nil {
		return nil, errors.WithStack(err)

	}
	c.Bitfield = bt
	return c, nil
}

func completeHandshake(conn net.Conn, peerID, infoHash [20]byte) error {
	sendHandshake := handshake.Handshake{
		Pstr:     "BitTorrent protocol",
		InfoHash: infoHash,
		PeerID:   peerID,
	}
	sendBytes := sendHandshake.Serialize()
	_, err := conn.Write(sendBytes)
	if err != nil {
		return err
	}
	reader := io.Reader(conn)
	recvHandshake, err := handshake.Read(reader)
	if err != nil {
		return errors.WithStack(err)
	}
	if recvHandshake.InfoHash != infoHash {
		return errors.New("infohash does not match")
	}
	return nil
}

func recvBitfield(conn net.Conn) (bitfield.Bitfield, error) {
	return bitfield.Read(conn)
}
