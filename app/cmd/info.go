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
	torrent, err := torrent.ParseTorrentFile(fileName)
	if err != nil {
		fmt.Printf("Erorr reading file %v, err: %v\n", fileName, err)
		return fmt.Errorf("failed to parse torrent file: %w", err)
	}

	buf := new(bytes.Buffer)
	err = bencode.Encode(buf, torrent.Info)
	if err != nil {
		log.Printf("Error while encoding info: %v", err)
		return fmt.Errorf("erro encoding: (%w)", err)
	}

	fmt.Println("Tracker URL:", torrent.Announce)
	fmt.Println("Length:", torrent.Info.Length)

	// compute sha-1 of buffer.
	sha1Hasher := sha1.New()
	sha1Hasher.Write(buf.Bytes())
	fmt.Println("Info Hash:", hex.EncodeToString(sha1Hasher.Sum(nil)))

	return nil
}
