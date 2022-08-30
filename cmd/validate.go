/*
Copyright © 2022 Yurii Rochniak yrochnyak@gmail.com

*/
package cmd

import (
	"github.com/grem11n/s3bc/action/validate"
	"github.com/grem11n/s3bc/config"
	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Check if files in a bucket have desired storage class.",
	Long: `Example usage:
s3bc validate -b example-bucket -s REDUCED_REDUNDANCY`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config := config.GetConfig(cmd.Flags())

		if err := validate.Run(config); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
