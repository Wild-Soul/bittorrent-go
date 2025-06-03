package torrent

import (
	"bufio"
	"fmt"
	"os"

	"github.com/codecrafters-io/bittorrent-starter-go/app/internal/bencode"
)

type Torrent struct {
	Announce string `bencode:"announce"`
	Info     Info   `bencode:"info"`
}

type Info struct {
	Length   int64  `bencode:"length"`
	Name     string `bencode:"name"`
	PieceLen int64  `bencode:"piece length"`
	Pieces   string `bencode:"pieces"`
}

// ParseTorrentFile reads and parses a .torrent file.
func ParseTorrentFile(path string) (*Torrent, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoded, err := bencode.Decode(reader)
	if err != nil {
		return nil, err
	}

	// Assert the top-level dictionary
	torrentMap, ok := decoded.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid torrent file structure")
	}

	// Parse announce
	announce, _ := torrentMap["announce"].(string)

	// Parse info dictionary
	infoMap, ok := torrentMap["info"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing or invalid info field")
	}

	// Extract info fields
	length, _ := infoMap["length"].(int64)
	name, _ := infoMap["name"].(string)
	pieceLen, _ := infoMap["piece length"].(int64)
	pieces, _ := infoMap["pieces"].(string)

	return &Torrent{
		Announce: announce,
		Info: Info{
			Length:   length,
			Name:     name,
			PieceLen: pieceLen,
			Pieces:   pieces,
		},
	}, nil
}
