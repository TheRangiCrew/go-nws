package cmd

import (
	"github.com/TheRangiCrew/go-nws/internal/awips"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(awipsCmd)
	awipsCmd.AddCommand(parseCmd)
}

var awipsCmd = &cobra.Command{
	Use:   "awips",
	Short: "Tools for AWIPS products",
}

var parseCmd = &cobra.Command{
	Use:   "parse <filename>",
	Short: "Parses the provided AWIPS product",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		awips.ParseAwips(args[0])
	},
}
