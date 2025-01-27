package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const microVersion = "1"
const minorVersion = "0"
const majorVersion = "0"

var rootCmd = &cobra.Command{
	Use:     "huffman",
	Version: fmt.Sprintf("%s.%s.%s", majorVersion, minorVersion, microVersion),
	Short:   "huffman - a simple CLI to encode and decode files using the huffman algorithm",
	Long:    `Welcome to huffman! This is a simple CLI to encode and decode files using the huffman algorithm.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
