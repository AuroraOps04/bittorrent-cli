package torrentfile

import (
	"io"
	"net/url"
	"strconv"

	"github.com/jackpal/bencode-go"
	"google.golang.org/genproto/googleapis/cloud/aiplatform/v1/schema/predict/params"
)

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
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
	PicesHashes  [][20]byte
	PieceLength  int
	Length       int
	Name         string
}

func (bto bencodeTorrent) toTorrentFile() (TorrentFile, error) {

	var t TorrentFile
	t.Announce = bto.Announce
	t.Comment = bto.Comment
	t.CreationDate = bto.CreationDate
	t.PieceLength = bto.Info.PieceLength
	t.Length = bto.Info.Length
	t.Name = bto.Info.Name
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
		"perr_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}
	base.RawQuery = params.Encode()
	return base.User.String(), nil
}
