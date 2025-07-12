package peermanager

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strconv"
	"time"
)

const (
	bitTorrentProtocolLength = 19
	bitTorrentProtocolName   = "BitTorrent protocol"
	connectionTimeout        = 5 * time.Second
	reservedByteLen          = 8
)

type HandshakeResponse struct {
	protocolName string
	infoHash     [20]byte
	peerID       [20]byte
}

func (h *HandshakeResponse) PeerID() [20]byte {
	return h.peerID
}

func (pm *PeerManager) Handshake(ctx context.Context, peer string) (*HandshakeResponse, error) {

	tcpConn, err := net.DialTimeout("tcp", peer, connectionTimeout)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Error while connecting to peer: %s", err.Error()))
		return &HandshakeResponse{}, err
	}
	defer tcpConn.Close()

	// Write to tcp connection.
	if err = pm.writeHandshake(tcpConn); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Failed to write to tcp connection: %s", err.Error()))
		return &HandshakeResponse{}, err
	}

	resp, err := pm.readHandshakeResponse(tcpConn)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Failed to read handshake response: %s", err.Error()))
		return &HandshakeResponse{}, err
	}

	return resp, nil
}

func (pm *PeerManager) writeHandshake(tcpConn net.Conn) error {
	reservedBytes := make([]byte, 8)
	sha1Hash, err := pm.t.GetInfohash()
	if err != nil {
		return err
	}
	peerID := pm.GetPeerID()

	handshake := make([]byte, 49+len(bitTorrentProtocolName))
	handshake[0] = byte(len(bitTorrentProtocolName))
	copy(handshake[1:], []byte(bitTorrentProtocolName))
	copy(handshake[1+len(bitTorrentProtocolName):], reservedBytes)
	copy(handshake[1+len(bitTorrentProtocolName)+reservedByteLen:], sha1Hash[:])
	copy(handshake[1+len(bitTorrentProtocolName)+reservedByteLen+20:], peerID[:])

	// Send handshake
	_, err = tcpConn.Write(handshake)
	if err != nil {
		return err
	}
	return nil
}

func (m *PeerManager) readHandshakeResponse(conn net.Conn) (*HandshakeResponse, error) {
	protocolLen := make([]byte, 1)
	if _, err := io.ReadFull(conn, protocolLen); err != nil {
		return nil, err
	}

	lenOfProtocolName, err := strconv.Atoi(strconv.Itoa(int(protocolLen[0])))
	if err != nil {
		return nil, err
	}

	protocolName := make([]byte, lenOfProtocolName)

	if _, err = io.ReadFull(conn, protocolName); err != nil {
		return nil, err
	}

	if _, err = io.ReadFull(conn, make([]byte, reservedByteLen)); err != nil {
		return nil, err
	}

	var infoHash [20]byte
	if _, err = io.ReadFull(conn, infoHash[:]); err != nil {
		return nil, err
	}

	var peerID [20]byte
	if _, err = io.ReadFull(conn, peerID[:]); err != nil {
		return nil, err
	}

	return &HandshakeResponse{
		protocolName: string(protocolName),
		infoHash:     infoHash,
		peerID:       peerID,
	}, nil
}
