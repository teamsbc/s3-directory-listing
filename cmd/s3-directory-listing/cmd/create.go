package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create <bucket> <template>",
	Short: "Create directory listings for an S3 bucket",
	Long: `Generate directory listings for the specified S3 bucket using the provided template.
Credentials and endpoint configuration are read from environment variables.`,
	Args: cobra.ExactArgs(2),
	RunE: runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func runCreate(cmd *cobra.Command, args []string) error {
	bucket := args[0]
	template := args[1]

	if err := validateBucket(bucket); err != nil {
		return fmt.Errorf("invalid bucket: %w", err)
	}

	if err := validateTemplate(template); err != nil {
		return fmt.Errorf("invalid template: %w", err)
	}

	fmt.Printf("Creating directory listing for bucket: %s\n", bucket)
	fmt.Printf("Using template: %s\n", template)

	// TODO: Implement S3 querying and template rendering
	return nil
}

func validateBucket(bucket string) error {
	if bucket == "" {
		return fmt.Errorf("bucket name cannot be empty")
	}
	return nil
}

func validateTemplate(template string) error {
	if template == "" {
		return fmt.Errorf("template path cannot be empty")
	}

	// Check if template file exists
	if _, err := os.Stat(template); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("template file does not exist: %s", template)
		}
		return fmt.Errorf("cannot access template file: %w", err)
	}

	return nil
}
