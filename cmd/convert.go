/*

Copyright © 2022 Yurii Rochniak yrochnyak@gmail.com

*/
package cmd

import (
	"log"

	"github.com/grem11n/s3bc/action/convert"
	"github.com/grem11n/s3bc/config"
	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Bulk convert objects in an S3 bucket to the given storage class",
	Long: `Example usage:
s3bc convert -b example-bucket -s REDUCED_REDUNDANCY`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config := config.GetConfig(cmd.Flags())

		if err := convert.Run(config); err != nil {
			log.Fatalf("Error running convert command: %s", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringSliceP(
		"exclude",
		"e",
		nil,
		"Patterns to exclude from the conversion. You can use Go Regexp. Also, you can provide multiple patterns separated by comma.",
	)

	convertCmd.Flags().Bool(
		"dry-run",
		false,
		"Dry run only outputs the list of keys to be converted.",
	)
}
