package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "s3-directory-listing",
	Short: "Generate directory listings for S3-compatible storage",
	Long: `A tool to generate customizable directory listings for S3-compatible storage
such as CloudFlare R2. This program reads bucket contents and generates HTML
directory listings based on templates.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("s3-directory-listing - use --help for available commands")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags can be added here
}
