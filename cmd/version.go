package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const Version = "v0.1"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "v0.0.1 for jar finder",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		color.HiGreen("version %s", Version)
	},
}
