package cmd

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/codecrafters-io/bittorrent-starter-go/app/internal/bencode"
)

type DecodeCmd struct{}

func (d *DecodeCmd) Name() string { return "decode" }

func (d *DecodeCmd) Description() string { return "Decode a bencoded string" }

func (d *DecodeCmd) Execute(ctx context.Context, args []string) error {
	if len(args) < 1 {
		log.Fatal("missing filename argument")
	}
	bencodedValue := args[0]
	inputBytes := []byte(bencodedValue)

	byteReader := bytes.NewReader(inputBytes)
	decoded, err := bencode.DecodeBencode(bufio.NewReader(byteReader))
	if err != nil {
		fmt.Printf("Failed to parse value %v\n", bencodedValue)
		return err
	}

	jsonOutput, _ := json.Marshal(decoded)
	fmt.Println(string(jsonOutput))
	return nil
}
