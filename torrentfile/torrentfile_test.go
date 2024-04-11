package torrentfile

import (
	"fmt"
	"os"
	"testing"

	"github.com/AuroraOps04/bittorrent-cli/peer"
)

func TestOpen(t *testing.T) {
	f, err := os.Open("../testdata/debian-12.5.0-amd64-netinst.iso.torrent")
	if err != nil {
		t.Errorf("Error opening file: %s", err)
	}
	bto, err := open(f)
	if err != nil {
		t.Errorf("error to parse torrent file: %s", err)
	}
	fmt.Printf("%v\n", *bto)
}

func TestGetPeers(t *testing.T) {

	f, err := os.Open("../testdata/debian-12.5.0-amd64-netinst.iso.torrent")
	if err != nil {
		t.Errorf("Error opening file: %s", err)
	}
	tf, err := New(f)
	if err != nil {
		t.Fatalf("Error create torrentfile: %v", err)
	}
	peers, err := tf.GetPeers(peer.GetLocalPeerID())
	if err != nil {
		t.Fatalf("Get peers error: %v", err)
	}
	t.Log(peers)
}
