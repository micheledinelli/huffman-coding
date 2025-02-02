package cmd

import (
	"fmt"
	"huffman/encoder"

	"github.com/spf13/cobra"
)

var compressedFilename string

var compressCmd = &cobra.Command{
	Use:     "compress",
	Aliases: []string{"encode", "c"},
	Short:   "Compress a file",
	Args:    cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Compressing %s with Huffman algorithm\n", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		encoder.Encode(args[0], &compressedFilename)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s compressed\n", args[0])
	},
}

func init() {
	compressCmd.Flags().StringVarP(&compressedFilename, "out", "o", "", "Output filename")
	rootCmd.AddCommand(compressCmd)
}
