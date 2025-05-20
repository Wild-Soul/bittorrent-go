package main

import (
	"context"
	"fmt"
	"os"

	"github.com/codecrafters-io/bittorrent-starter-go/app/cmd"
)

func main() {
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")
	ctx := context.Background()
	registry := cmd.NewRegistry()

	// Register available commands
	registry.Register(&cmd.DecodeCmd{})
	registry.Register(&cmd.InfoCmd{})

	cmdName := os.Args[1]
	args := os.Args[2:]

	cmd, found := registry.Get(cmdName)
	if !found {
		fmt.Println("Unknown command: " + cmdName)
		return
	}

	if err := cmd.Execute(ctx, args); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
