package torrentfile

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/AuroraOps04/bittorrent-cli/client"
	"github.com/AuroraOps04/bittorrent-cli/peer"
	"github.com/stretchr/testify/assert"
)

func TestRequestPeers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := []byte(
			"d" +
				"8:interval" + "i900e" +
				"5:peers" + "12:" +
				string([]byte{
					192, 0, 2, 123, 0x1A, 0xE1, // 0x1AE1 = 6881
					127, 0, 0, 1, 0x1A, 0xE9, // 0x1AE9 = 6889
				}) + "e")
		w.Write(response)
	}))
	defer ts.Close()
	tf := TorrentFile{
		Announce: ts.URL,
		InfoHash: [20]byte{216, 247, 57, 206, 195, 40, 149, 108, 204, 91, 191, 31, 134, 217, 253, 207, 219, 168, 206, 182},
		PieceHashes: [][20]byte{
			{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106},
			{97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48},
		},
		PieceLength: 262144,
		Length:      351272960,
		Name:        "debian-10.2.0-amd64-netinst.iso",
	}
	peerID := [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	const port uint16 = 6882
	expected := []peer.Peer{
		{IP: net.IP{192, 0, 2, 123}, Port: 6881},
		{IP: net.IP{127, 0, 0, 1}, Port: 6889},
	}
	p, err := tf.GetPeers(peerID, port)
	assert.Nil(t, err)
	assert.Equal(t, expected, p)
}
func TestBuildTrackerURL(t *testing.T) {
	to := TorrentFile{
		Announce: "http://bttracker.debian.org:6969/announce",
		InfoHash: [20]byte{216, 247, 57, 206, 195, 40, 149, 108, 204, 91, 191, 31, 134, 217, 253, 207, 219, 168, 206, 182},
		PieceHashes: [][20]byte{
			{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106},
			{97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48},
		},
		PieceLength: 262144,
		Length:      351272960,
		Name:        "debian-10.2.0-amd64-netinst.iso",
	}
	peerID := [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	const port uint16 = 6882
	url, err := to.buildTrackerURL(peerID, port)
	expected := "http://bttracker.debian.org:6969/announce?compact=1&downloaded=0&info_hash=%D8%F79%CE%C3%28%95l%CC%5B%BF%1F%86%D9%FD%CF%DB%A8%CE%B6&left=351272960&peer_id=%01%02%03%04%05%06%07%08%09%0A%0B%0C%0D%0E%0F%10%11%12%13%14&port=6882&uploaded=0"
	assert.Nil(t, err)
	assert.Equal(t, url, expected)
}

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
	peers, err := tf.GetPeers(peer.GetLocalPeerID(), 6882)
	if err != nil {
		t.Fatalf("Get peers error: %v", err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(len(peers))
	for _, p := range peers {
		go func() {
			defer wg.Done()
			fmt.Println(p.String())
			c, err := client.New(p, peer.GetLocalPeerID(), tf.InfoHash)
			if err != nil {
				t.Error(err)
				return
			}
			fmt.Println(c.Bitfield)
		}()
	}
	wg.Wait()
}
