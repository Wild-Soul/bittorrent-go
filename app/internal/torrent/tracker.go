package torrent

import (
	"encoding/binary"
	"fmt"
	"net"
)

type TrackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func (tr *TrackerResponse) unmarshalPeers() ([]Peer, error) {
	const peerSize = 6 // 4 for IP, 2 for port
	peers := []byte(tr.Peers)

	numPeers := len(peers) / peerSize
	if numPeers*peerSize != len(peers) {
		err := fmt.Errorf("received malformed peers")
		return nil, err
	}

	parsedPeers := make([]Peer, numPeers)
	for i := range numPeers {
		offset := i * peerSize
		parsedPeers[i].Ip = net.IP(peers[offset : offset+4])
		parsedPeers[i].Port = binary.BigEndian.Uint16([]byte(peers[offset+4 : offset+6]))
	}

	return parsedPeers, nil
}
