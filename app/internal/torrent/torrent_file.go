package torrent

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/jackpal/bencode-go"
)

// lenght of hash of each peice in bytes
const PIECE_LEN = 20

type TorrentFile struct {
	Announce string `bencode:"announce"`
	Info     Info   `bencode:"info"`
}

type Info struct {
	Length   int64  `bencode:"length"`
	Name     string `bencode:"name"`
	PieceLen int64  `bencode:"piece length"`
	Pieces   []byte `bencode:"pieces"`
}

// ParseTorrentFile reads and parses a .torrent file.
func ParseTorrentFile(path string) (*TorrentFile, error) {
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
	piecesStr, _ := infoMap["pieces"].(string)
	pieces := []byte(piecesStr)

	return &TorrentFile{
		Announce: announce,
		Info: Info{
			Length:   length,
			Name:     name,
			PieceLen: pieceLen,
			Pieces:   pieces,
		},
	}, nil
}

func (tf *TorrentFile) GetInfohash() (string, error) {
	buf := new(bytes.Buffer)
	err := bencode.Marshal(buf, tf.Info)
	if err != nil {
		log.Printf("Error while encoding info: %v", err)
		return "", fmt.Errorf("error encoding: (%w)", err)
	}

	// compute info hash.
	// TODO:: Introuduce Torrent to keep only torrent info, that way won't have to do decode/encode.
	sha1Hasher := sha1.New()
	sha1Hasher.Write(buf.Bytes())
	infoHash := hex.EncodeToString(sha1Hasher.Sum(nil))

	return infoHash, nil
}
