package torrentfile

import (
	"fmt"
	"os"
	"testing"
)

func TestOpen(t *testing.T) {
	f, err := os.Open("../testdata/debian-12.5.0-amd64-netinst.iso.torrent")
	if err != nil {
		t.Errorf("Error opening file: %s", err)
	}
	bto, err := Open(f)
	if err != nil {
		t.Errorf("error to parse torrent file: %s", err)
	}
	fmt.Printf("%v\n",*bto)
}
