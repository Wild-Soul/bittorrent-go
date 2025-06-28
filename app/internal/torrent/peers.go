package torrent

import (
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jackpal/bencode-go"
)

type Peer struct {
	Ip   net.IP
	Port uint16
}

func (p Peer) String() string {
	return net.JoinHostPort(p.Ip.String(), strconv.Itoa(int(p.Port)))
}

func (tf *TorrentFile) generatePeersReqUrl(peerId [20]byte, port uint16) (string, error) {
	parsedUrl, err := url.Parse(tf.Announce)
	if err != nil {
		return "", err
	}

	infohash, err := tf.GetInfohash()
	if err != nil {
		return "", err
	}

	params := url.Values{
		"info_hash":  []string{infohash},
		"peer_id":    []string{string(peerId[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(int(tf.Info.Length))},
	}

	parsedUrl.RawQuery = params.Encode()
	return parsedUrl.String(), nil
}

func (tf *TorrentFile) RequestPeers(peerId [20]byte, port uint16) ([]Peer, error) {
	url, err := tf.generatePeersReqUrl(peerId, port)

	if err != nil {
		return []Peer{}, err
	}

	res, err := http.Get(url)
	if err != nil {
		return []Peer{}, err
	}

	defer res.Body.Close()

	var trackerResponse TrackerResponse
	err = bencode.Unmarshal(res.Body, &trackerResponse)
	if err != nil {
		return []Peer{}, err
	}

	return trackerResponse.unmarshalPeers()
}
