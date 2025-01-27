package cmd

import (
	"fmt"
	"huffman/encoder"

	"github.com/spf13/cobra"
)

var decompressedFilename string

var decompressCmd = &cobra.Command{
	Use:     "decompress",
	Aliases: []string{"decode", "d"},
	Short:   "Decompress a file",
	Args:    cobra.ExactArgs(2),
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Decompressing %s\n", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		encoder.Decode(args[0], args[1], &decompressedFilename)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s decompressed\n", args[0])
	},
}

func init() {
	decompressCmd.Flags().StringVarP(&decompressedFilename, "out", "o", "", "Output filename")
	rootCmd.AddCommand(decompressCmd)
}
