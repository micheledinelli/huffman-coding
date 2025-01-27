package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number ",
	Long:  `All software has versions. This is Huffman's encoder/decoder`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Huffman encoder/decoder " + majorVersion + "." + minorVersion + "." + microVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
