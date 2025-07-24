package cmd

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

type DownloadPiece struct{}

func (command *DownloadPiece) Name() string { return "download_piece" }

func (command *DownloadPiece) Description() string { return "download a piece" }

func (command *DownloadPiece) Execute(ctx context.Context, args []string) error {

	downloadPieceCmd := flag.NewFlagSet("download_piece", flag.ExitOnError)
	output := downloadPieceCmd.String("o", "", "Output path")

	// Parse the flags for 'download_piece'
	downloadPieceCmd.Parse(args)

	// After flags, remaining args are positional
	params := downloadPieceCmd.Args()
	if len(params) < 2 {
		fmt.Println("usage: go_executable download_piece -o <output_path> <torrent_file> <piece_index>")
		os.Exit(1)
	}

	torrentFile := params[0]
	pieceIndexStr := params[1]
	pieceIndex, err := strconv.Atoi(pieceIndexStr)
	if err != nil {
		fmt.Printf("Invalid piece index: %s\n", pieceIndexStr)
		os.Exit(1)
	}

	slog.InfoContext(ctx, "Starting download_piece command",
		"torrent_file", torrentFile,
		"piece_index", pieceIndex,
		"output", *output,
	)

	return nil
}
