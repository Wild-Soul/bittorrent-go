package cmd

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/codecrafters-io/bittorrent-starter-go/app/internal/bencode"
	"github.com/codecrafters-io/bittorrent-starter-go/app/internal/torrent"
)

type InfoCmd struct{}

func (i *InfoCmd) Name() string        { return "info" }
func (i *InfoCmd) Description() string { return "Print info about the input" }

func (i *InfoCmd) Execute(ctx context.Context, args []string) error {
	if len(args) < 1 {
		log.Fatal("missing filename argument")
	}
	fileName := os.Args[2]
	parsedTorrent, err := torrent.ParseTorrentFile(fileName)
	if err != nil {
		fmt.Printf("Erorr reading file %v, err: %v\n", fileName, err)
		return fmt.Errorf("failed to parse torrent file: %w", err)
	}

	buf := new(bytes.Buffer)
	err = bencode.Encode(buf, parsedTorrent.Info)
	if err != nil {
		log.Printf("Error while encoding info: %v", err)
		return fmt.Errorf("error encoding: (%w)", err)
	}

	fmt.Println("Tracker URL:", parsedTorrent.Announce)
	fmt.Println("Length:", parsedTorrent.Info.Length)

	// compute sha-1 of buffer.
	sha1Hasher := sha1.New()
	sha1Hasher.Write(buf.Bytes())
	fmt.Println("Info Hash:", hex.EncodeToString(sha1Hasher.Sum(nil)))

	// Piece length and their hashses.
	fmt.Printf("Piece Length: %v\nPiece Hashes:\n", parsedTorrent.Info.PieceLen)
	for i := 0; i < len(parsedTorrent.Info.Pieces); i += 20 {
		fmt.Printf("%v\n", hex.EncodeToString(parsedTorrent.Info.Pieces[i:i+20]))
	}

	return nil
}
