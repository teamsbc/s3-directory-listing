package s3

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

func GenerateSHA256SUMS(ctx context.Context, client *Client, listing *DirectoryListing) (string, error) {
	type entry struct {
		filename string
		hash     string
	}

	var entries []entry
	for _, f := range listing.Files {
		if !strings.HasSuffix(f.Name, ".sha256") {
			continue
		}

		content, err := client.GetObjectContent(ctx, f.Key)
		if err != nil {
			return "", fmt.Errorf("failed to fetch %s: %w", f.Key, err)
		}

		hash := strings.TrimSpace(string(content))
		// Handle both "hash  filename" and bare hash formats
		if idx := strings.IndexAny(hash, " \t"); idx != -1 {
			hash = hash[:idx]
		}

		filename := strings.TrimSuffix(f.Name, ".sha256")
		entries = append(entries, entry{filename: filename, hash: hash})
	}

	if len(entries) == 0 {
		return "", nil
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].filename < entries[j].filename
	})

	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "%s  %s\n", e.hash, e.filename)
	}

	return sb.String(), nil
}
