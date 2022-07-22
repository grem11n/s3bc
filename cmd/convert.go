/*

Copyright © 2022 Yurii Rochniak yrochnyak@gmail.com

*/
package cmd

import (
	"github.com/grem11n/s3bc/action/convert"
	"github.com/grem11n/s3bc/config"
	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Bulk convert objects in an S3 bucket to the given Storage Class",
	Long: `Example usage:
s3bc ...
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket, err := cmd.Flags().GetString("bucket")
		if err != nil {
			return err
		}

		storageClass, err := cmd.Flags().GetString("storage-class")
		if err != nil {
			return err
		}

		exclude, err := cmd.Flags().GetStringSlice("exclude")
		if err != nil {
			return err
		}

		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			return err
		}

		config := &config.Config{
			Bucket:       bucket,
			StorageClass: storageClass,
			Excluded:     exclude,
			DryRun:       dryRun,
		}

		if err := convert.Run(config); err != nil {
			return err
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
