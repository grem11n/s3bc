/*
Copyright © 2022 Yurii Rochniak yrochnyak@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "s3bc",
	Short: "Bulk Storage Class update in an S3 bucket.",
	Long:  `S3BC or S3 Bulk Convert is a CLI tool to update the storage class of the files in an AWS S3 bucket.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cmd, _, err := rootCmd.Find(os.Args[1:])
	// Default to convert command if no other subcommand is provided.
	if err == nil && cmd.Use == rootCmd.Use && cmd.Flags().Parse(os.Args[1:]) != pflag.ErrHelp {
		args := append([]string{convertCmd.Use}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}
	if err := rootCmd.Execute(); err != nil {
		log.Fatal()
	}
}

func init() {
	rootCmd.PersistentFlags().StringP(
		"bucket",
		"b",
		"",
		"Target S3 bucket",
	)

	rootCmd.PersistentFlags().StringP(
		"storage-class",
		"s",
		"STANDARD",
		"Storage class to set",
	)

	rootCmd.PersistentFlags().BoolP(
		"verbose",
		"v",
		false,
		"Verbose output",
	)
}
