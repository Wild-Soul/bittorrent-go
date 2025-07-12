package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"

	peermanager "github.com/codecrafters-io/bittorrent-starter-go/app/internal/peer_manager"
	"github.com/codecrafters-io/bittorrent-starter-go/app/internal/torrent"
)

type Handshake struct{}

func (command *Handshake) Name() string { return "handshake" }

func (command *Handshake) Description() string { return "peer handshake" }

func (command *Handshake) Execute(ctx context.Context, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("invalid argument: missing torrent file or peer address")
	}

	torrentFile, err := torrent.ParseTorrentFile(args[0])
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Error while parsing torrent file: %s", err.Error()))
		return err
	}

	peerManager := peermanager.NewPeerManager(torrentFile)

	handshakeResponse, err := peerManager.Handshake(ctx, args[1])
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Handshake failed: %s", err.Error()))
		return err
	}

	peerId := handshakeResponse.PeerID()
	fmt.Println("Peer ID:", hex.EncodeToString(peerId[:]))
	return nil
}
