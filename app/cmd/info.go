package cmd

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/codecrafters-io/bittorrent-starter-go/app/internal/bencode"
)

type InfoCmd struct{}
type torrentStruct struct {
}

func (i *InfoCmd) Name() string        { return "info" }
func (i *InfoCmd) Description() string { return "Print info about the input" }

func (i *InfoCmd) Execute(ctx context.Context, args []string) error {
	if len(args) < 1 {
		log.Fatal("missing filename argument")
	}
	fileName := os.Args[2]
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Erorr reading file %v, err: %v\n", fileName, err)
	}

	byteReader := bytes.NewReader(data)
	dictionary, err := bencode.DecodeBencode(bufio.NewReader(byteReader))
	if err != nil {
		fmt.Printf("Failed to parse file %v\n", fileName)
		return err
	}

	// TODO:: Create TorrentFileStruct and parse into that.
	if assertedData, ok := dictionary.(map[string]interface{}); ok {
		fmt.Println("Tracker URL:", assertedData["announce"])
		if info, ok := assertedData["info"].(map[string]interface{}); ok {
			fmt.Println("Length:", info["length"])
		}
	}

	return nil
}
