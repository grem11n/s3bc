/*

Copyright © 2022 Yurii Rochniak yrochnyak@gmail.com

*/
package cmd

import (
	"github.com/grem11n/s3bc/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print s3bc version.",
	Run: func(cmd *cobra.Command, args []string) {
		version.PrintInfo()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
