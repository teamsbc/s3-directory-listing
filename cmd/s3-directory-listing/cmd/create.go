package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/teamsbc/s3-directory-listing/internal/s3"
	"github.com/teamsbc/s3-directory-listing/internal/template"
)

var (
	outputDir string
	profile   string
)

var createCmd = &cobra.Command{
	Use:   "create <bucket> <template>",
	Short: "Create directory listings for an S3 bucket",
	Long: `Generate directory listings for the specified S3 bucket using the provided template.
Credentials and endpoint configuration are read from environment variables or AWS config files.`,
	Args: cobra.ExactArgs(2),
	RunE: runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory for generated listings")
	createCmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS profile to use from config file")
}

func runCreate(cmd *cobra.Command, args []string) error {
	bucket := args[0]
	templatePath := args[1]

	if err := validateBucket(bucket); err != nil {
		return fmt.Errorf("invalid bucket: %w", err)
	}

	if err := validateTemplate(templatePath); err != nil {
		return fmt.Errorf("invalid template: %w", err)
	}

	ctx := context.Background()

	// Create S3 client
	fmt.Printf("Connecting to bucket: %s\n", bucket)
	if profile != "" {
		fmt.Printf("Using AWS profile: %s\n", profile)
	}
	client, err := s3.NewClient(ctx, bucket, profile)
	if err != nil {
		return fmt.Errorf("failed to create S3 client: %w", err)
	}

	// List all objects in the bucket
	fmt.Println("Listing objects...")
	objects, err := client.ListAllObjects(ctx)
	if err != nil {
		return fmt.Errorf("failed to list objects: %w", err)
	}
	fmt.Printf("Found %d objects\n", len(objects))

	// Build directory structure
	fmt.Println("Building directory structure...")
	directories := s3.BuildDirectoryStructure(objects)
	fmt.Printf("Found %d directories\n", len(directories))

	// Create template renderer
	renderer, err := template.NewRenderer(templatePath)
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	// Render templates for each directory
	fmt.Printf("Rendering templates to: %s\n", outputDir)
	for dirPath, listing := range directories {
		// Create output path
		var outPath string
		if dirPath == "" {
			outPath = filepath.Join(outputDir, "index.html")
		} else {
			outPath = filepath.Join(outputDir, dirPath, "index.html")
		}

		// Ensure output directory exists
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		// Render template
		if err := renderer.RenderToFile(listing, outPath); err != nil {
			return fmt.Errorf("failed to render template for %s: %w", dirPath, err)
		}

		fmt.Printf("  Created: %s\n", outPath)
	}

	fmt.Println("Done!")
	return nil
}

func validateBucket(bucket string) error {
	if bucket == "" {
		return fmt.Errorf("bucket name cannot be empty")
	}
	return nil
}

func validateTemplate(templatePath string) error {
	if templatePath == "" {
		return fmt.Errorf("template path cannot be empty")
	}

	// Check if template file exists
	if _, err := os.Stat(templatePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("template file does not exist: %s", templatePath)
		}
		return fmt.Errorf("cannot access template file: %w", err)
	}

	return nil
}
