package client

import (
	"io"
	"net"
	"time"

	"github.com/AuroraOps04/bittorrent-cli/bitfield"
	"github.com/AuroraOps04/bittorrent-cli/handshake"
	"github.com/AuroraOps04/bittorrent-cli/message"
	"github.com/AuroraOps04/bittorrent-cli/peers"
	"github.com/pkg/errors"
)

type Client struct {
	Conn     net.Conn
	Bitfield bitfield.Bitfield
	PeerID   [20]byte
	InfoHash [20]byte
	Peer     peers.Peer
	// Choked is true if the peer is choked
	Choked bool
}

func New(p peers.Peer, id, infoHash [20]byte) (*Client, error) {
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
		conn.Close()
		return nil, errors.WithStack(err)
	}
	bt, err := recvBitfield(conn)
	if err != nil {
		conn.Close()
		return nil, errors.WithStack(err)

	}
	c.Bitfield = bt
	return c, nil
}
func (c *Client) Read() (*message.Message, error) {
	return message.Read(c.Conn)
}
func (c *Client) SendRequest(index, begin, length int) error {
	return nil
}

// SendInterested sends an Interested message to the peer
func (c *Client) SendInterested() error {
	msg := message.Message{ID: message.MsgInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendNotInterested sends a NotInterested message to the peer
func (c *Client) SendNotInterested() error {
	msg := message.Message{ID: message.MsgNotInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendUnchoke sends an Unchoke message to the peer
func (c *Client) SendUnchoke() error {
	msg := message.Message{ID: message.MsgUnchoke}
	_, err := c.Conn.Write(msg.Serialize())
	return err
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
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{})
	// read a message
	msg, err := message.Read(conn)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// check message is bitfield Message
	// keep-alive message
	if msg == nil {
		return nil, errors.New("expected bitfield, got nil")
	}
	if msg.ID != message.MsgBitfield {
		return nil, errors.New("expected bitfield, got different message")
	}
	// msg.Payload is bitfield
	return msg.Payload, nil
}
