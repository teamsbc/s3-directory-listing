package s3

import (
	"path"
	"strings"
)

type DirectoryEntry struct {
	Name string
	Size int64
}

type DirectoryListing struct {
	Path        string
	Directories []DirectoryEntry
	Files       []DirectoryEntry
}

// BuildDirectoryStructure organizes S3 objects into directory listings
func BuildDirectoryStructure(objects []Object) map[string]*DirectoryListing {
	directories := make(map[string]*DirectoryListing)

	// Ensure root directory exists
	directories[""] = &DirectoryListing{
		Path:        "",
		Directories: []DirectoryEntry{},
		Files:       []DirectoryEntry{},
	}

	// Process all objects to build directory structure
	for _, obj := range objects {
		// Skip if it's just a directory marker (ends with /)
		if strings.HasSuffix(obj.Key, "/") {
			continue
		}

		// Get the directory path
		dir := path.Dir(obj.Key)
		if dir == "." {
			dir = ""
		}

		// Ensure all parent directories exist
		ensureParentDirectories(directories, dir)

		// Add file to its directory
		if _, exists := directories[dir]; !exists {
			directories[dir] = &DirectoryListing{
				Path:        dir,
				Directories: []DirectoryEntry{},
				Files:       []DirectoryEntry{},
			}
		}

		directories[dir].Files = append(directories[dir].Files, DirectoryEntry{
			Name: path.Base(obj.Key),
			Size: obj.Size,
		})
	}

	// Build directory entries (subdirectories within each directory)
	for dirPath := range directories {
		if dirPath == "" {
			continue
		}

		parentDir := path.Dir(dirPath)
		if parentDir == "." {
			parentDir = ""
		}

		if parent, exists := directories[parentDir]; exists {
			dirName := path.Base(dirPath)
			// Check if not already added
			found := false
			for _, d := range parent.Directories {
				if d.Name == dirName {
					found = true
					break
				}
			}
			if !found {
				parent.Directories = append(parent.Directories, DirectoryEntry{
					Name: dirName,
					Size: 0,
				})
			}
		}
	}

	return directories
}

func ensureParentDirectories(directories map[string]*DirectoryListing, dirPath string) {
	if dirPath == "" || dirPath == "." {
		return
	}

	if _, exists := directories[dirPath]; exists {
		return
	}

	directories[dirPath] = &DirectoryListing{
		Path:        dirPath,
		Directories: []DirectoryEntry{},
		Files:       []DirectoryEntry{},
	}

	parent := path.Dir(dirPath)
	if parent != "." {
		ensureParentDirectories(directories, parent)
	}
}
