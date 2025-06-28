package cmd

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"

	"github.com/codecrafters-io/bittorrent-starter-go/app/internal/torrent"
)

const Port uint16 = 6681

type PeersCmd struct{}

func (peerCmd *PeersCmd) Name() string        { return "peers" }
func (peerCmd *PeersCmd) Description() string { return "Discover peers" }

func (peerCmd *PeersCmd) Execute(tx context.Context, args []string) error {
	if len(args) < 1 {
		log.Fatal("missing filepath argument")
	}

	filepath := args[0]
	torrentFile, err := torrent.ParseTorrentFile(filepath)
	if err != nil {
		fmt.Printf("Error while reading file: %s, err: %s\n", filepath, err)
		return err
	}

	peerId := make([]byte, 20)
	rand.Read(peerId)

	peers, err := torrentFile.RequestPeers([20]byte(peerId), Port)
	if err != nil {
		log.Fatal(err)
	}

	for _, peer := range peers {
		fmt.Printf("%s\n", peer)
	}

	return nil
}
