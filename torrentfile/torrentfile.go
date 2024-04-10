package torrentfile

import (
	"github.com/jackpal/bencode-go"
	"io"
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
