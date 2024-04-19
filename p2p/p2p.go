package p2p

import (
	"github.com/AuroraOps04/bittorrent-cli/client"
	"github.com/AuroraOps04/bittorrent-cli/message"
	"github.com/AuroraOps04/bittorrent-cli/peers"
)

type pieceState struct {
	// index of the piece in bitfield
	index int
	// store the data in memory
	buf        []byte
	client     *client.Client
	downloaded int
	requested  int
	backlog    int
}

// readMessage process the message envent
func (state *pieceState) readMessage() error {
	msg, err := state.client.Read()
	if err != nil {
		return err
	}

	// keep-alive
	if msg == nil {
		return nil
	}
	switch msg.ID {
	case message.MsgUnchoke:
		state.client.Choked = false
	case message.MsgChoke:
		state.client.Choked = true
		// the peer have new pieces
	case message.MsgHave:
	// recive the piece data
	case message.MsgPiece:
		// message.
	}
	return nil
}

func attemptDownloadPiece(peer peers.Peer) error{
	return nil
}

