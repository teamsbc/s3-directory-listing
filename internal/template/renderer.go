package template

import (
	"bytes"
	"html/template"
	"os"

	"github.com/teamsbc/s3-directory-listing/internal/s3"
)

type TemplateData struct {
	Path        string
	Directories []s3.DirectoryEntry
	Files       []s3.DirectoryEntry
}

type Renderer struct {
	tmpl *template.Template
}

func NewRenderer(templatePath string) (*Renderer, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	return &Renderer{tmpl: tmpl}, nil
}

func (r *Renderer) Render(listing *s3.DirectoryListing) (string, error) {
	data := TemplateData{
		Path:        listing.Path,
		Directories: listing.Directories,
		Files:       listing.Files,
	}

	var buf bytes.Buffer
	if err := r.tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (r *Renderer) RenderToFile(listing *s3.DirectoryListing, outputPath string) error {
	content, err := r.Render(listing)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, []byte(content), 0644)
}
