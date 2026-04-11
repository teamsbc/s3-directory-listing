package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/teamsbc/s3-directory-listing/internal/s3"
)

func TestNewRenderer(t *testing.T) {
	// Create a temporary template file
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "test.html")
	templateContent := `<html><body>{{.Path}}</body></html>`
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("failed to create test template: %v", err)
	}

	renderer, err := NewRenderer(templatePath)
	if err != nil {
		t.Errorf("NewRenderer() error = %v, want nil", err)
	}
	if renderer == nil {
		t.Error("NewRenderer() returned nil renderer")
	}
}

func TestNewRenderer_NonExistentFile(t *testing.T) {
	_, err := NewRenderer("/nonexistent/template.html")
	if err == nil {
		t.Error("NewRenderer() with non-existent file should return error")
	}
}

func TestRenderer_Render(t *testing.T) {
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "test.html")
	templateContent := `Path: {{.Path}}
Directories: {{range .Directories}}{{.Name}} {{end}}
Files: {{range .Files}}{{.Name}} {{end}}`
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("failed to create test template: %v", err)
	}

	renderer, err := NewRenderer(templatePath)
	if err != nil {
		t.Fatalf("NewRenderer() error = %v", err)
	}

	listing := &s3.DirectoryListing{
		Path: "test/path",
		Directories: []s3.DirectoryEntry{
			{Name: "dir1", Size: 0},
			{Name: "dir2", Size: 0},
		},
		Files: []s3.DirectoryEntry{
			{Name: "file1.txt", Size: 100},
			{Name: "file2.txt", Size: 200},
		},
	}

	result, err := renderer.Render(listing)
	if err != nil {
		t.Errorf("Render() error = %v, want nil", err)
	}

	if !strings.Contains(result, "Path: test/path") {
		t.Error("Render() output does not contain expected path")
	}
	if !strings.Contains(result, "dir1") || !strings.Contains(result, "dir2") {
		t.Error("Render() output does not contain expected directories")
	}
	if !strings.Contains(result, "file1.txt") || !strings.Contains(result, "file2.txt") {
		t.Error("Render() output does not contain expected files")
	}
}

func TestRenderer_RenderToFile(t *testing.T) {
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "test.html")
	templateContent := `<html><body>{{.Path}}</body></html>`
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("failed to create test template: %v", err)
	}

	renderer, err := NewRenderer(templatePath)
	if err != nil {
		t.Fatalf("NewRenderer() error = %v", err)
	}

	listing := &s3.DirectoryListing{
		Path:        "test",
		Directories: []s3.DirectoryEntry{},
		Files:       []s3.DirectoryEntry{},
	}

	outputPath := filepath.Join(tmpDir, "output.html")
	err = renderer.RenderToFile(listing, outputPath)
	if err != nil {
		t.Errorf("RenderToFile() error = %v, want nil", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("RenderToFile() did not create output file")
	}

	// Verify content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if !strings.Contains(string(content), "test") {
		t.Error("RenderToFile() output does not contain expected content")
	}
}
