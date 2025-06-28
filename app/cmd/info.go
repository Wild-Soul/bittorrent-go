package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/codecrafters-io/bittorrent-starter-go/app/internal/torrent"
)

type InfoCmd struct{}

func (i *InfoCmd) Name() string        { return "info" }
func (i *InfoCmd) Description() string { return "Print info about the input" }

func (i *InfoCmd) Execute(ctx context.Context, args []string) error {
	if len(args) < 1 {
		log.Fatal("missing filename argument")
	}

	fileName := args[0]
	torrentFile, err := torrent.ParseTorrentFile(fileName)
	if err != nil {
		fmt.Printf("Erorr reading file %v, err: %v\n", fileName, err)
		return fmt.Errorf("failed to parse torrent file: %w", err)
	}

	fmt.Println("Tracker URL:", torrentFile.Announce)
	fmt.Println("Length:", torrentFile.Info.Length)

	infohash, err := torrentFile.GetInfohash()
	if err != nil {
		return err
	}
	fmt.Println("Info Hash:", infohash)

	// Piece length and their hashses.
	fmt.Printf("Piece Length: %v\nPiece Hashes:\n", torrentFile.Info.PieceLen)
	for i := 0; i < len(torrentFile.Info.Pieces); i += 20 {
		fmt.Printf("%v\n", hex.EncodeToString(torrentFile.Info.Pieces[i:i+20]))
	}

	return nil
}
