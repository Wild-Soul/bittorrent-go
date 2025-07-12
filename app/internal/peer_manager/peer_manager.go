package peermanager

import (
	"crypto/rand"

	"github.com/codecrafters-io/bittorrent-starter-go/app/internal/torrent"
)

type PeerManager struct {
	t      *torrent.TorrentFile
	peerId []byte
}

func NewPeerManager(t *torrent.TorrentFile) *PeerManager {
	peerId := make([]byte, 20)
	rand.Read(peerId)

	return &PeerManager{
		t:      t,
		peerId: peerId,
	}
}

func (pm *PeerManager) GetPeerID() []byte {
	return pm.peerId
}
