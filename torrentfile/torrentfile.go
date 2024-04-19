package torrentfile

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/AuroraOps04/bittorrent-cli/peers"
	"github.com/jackpal/bencode-go"
	"github.com/pkg/errors"
)

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

func (i *bencodeInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

type bencodeTorrent struct {
	Announce     string      `bencode:"announce"`
	Info         bencodeInfo `bencode:"info"`
	Comment      string      `bencode:"comment"`
	CreationDate int64       `bencode:"creation date"`
}

func open(r io.Reader) (*bencodeTorrent, error) {
	var bto bencodeTorrent
	err := bencode.Unmarshal(r, &bto)
	if err != nil {
		return nil, err
	}
	return &bto, err
}

type TorrentFile struct {
	Announce     string
	Comment      string
	CreationDate int64
	InfoHash     [20]byte
	PieceHashes  [][20]byte
	PieceLength  int
	Length       int
	Name         string
}

func (bto bencodeTorrent) toTorrentFile() (TorrentFile, error) {

	h, err := bto.Info.hash()
	if err != nil {
		return TorrentFile{}, err
	}
	var t TorrentFile
	t.Announce = bto.Announce
	t.Comment = bto.Comment
	t.CreationDate = bto.CreationDate
	t.PieceLength = bto.Info.PieceLength
	t.Length = bto.Info.Length
	t.Name = bto.Info.Name
	t.InfoHash = h
	return t, nil
}

func New(r io.Reader) (*TorrentFile, error) {
	bto, err := open(r)
	if err != nil {
		return nil, err
	}
	t, err := bto.toTorrentFile()
	return &t, err

}

func (t *TorrentFile) buildTrackerURL(peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}
func (t *TorrentFile) GetPeers(peerId [20]byte, port uint16) ([]peers.Peer, error) {
	url, err := t.buildTrackerURL(peerId, port)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// 不使用代理
	c := http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			Proxy: nil,
		},
	}
	resp, err := c.Get(url)
	if err != nil {
		return nil, errors.WithMessage(err, "request tracker")

	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request tracker failed status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	trackerResp := peers.TrackerResponse{}
	// 解析peers有问题，后面好像多了个peers的个数
	// byteArr, err := io.ReadAll(resp.Body)

	err = bencode.Unmarshal(resp.Body, &trackerResp)
	if err != nil {
		return nil, errors.WithMessage(err, "unmarshal tracker response")
	}
	return peers.Unmarshal([]byte(trackerResp.Peers))
}
