package s3

import (
	"testing"
)

func TestBuildDirectoryStructure(t *testing.T) {
	tests := []struct {
		name     string
		objects  []Object
		wantDirs int
	}{
		{
			name:     "empty bucket",
			objects:  []Object{},
			wantDirs: 1, // root directory always exists
		},
		{
			name: "files in root only",
			objects: []Object{
				{Key: "file1.txt", Size: 100},
				{Key: "file2.txt", Size: 200},
			},
			wantDirs: 1,
		},
		{
			name: "single directory",
			objects: []Object{
				{Key: "dir1/file1.txt", Size: 100},
				{Key: "dir1/file2.txt", Size: 200},
			},
			wantDirs: 2, // root + dir1
		},
		{
			name: "nested directories",
			objects: []Object{
				{Key: "dir1/dir2/file1.txt", Size: 100},
				{Key: "dir1/dir2/file2.txt", Size: 200},
			},
			wantDirs: 3, // root + dir1 + dir1/dir2
		},
		{
			name: "multiple directories",
			objects: []Object{
				{Key: "dir1/file1.txt", Size: 100},
				{Key: "dir2/file2.txt", Size: 200},
				{Key: "file3.txt", Size: 300},
			},
			wantDirs: 3, // root + dir1 + dir2
		},
		{
			name: "directory markers ignored",
			objects: []Object{
				{Key: "dir1/", Size: 0},
				{Key: "dir1/file1.txt", Size: 100},
			},
			wantDirs: 2, // root + dir1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dirs := BuildDirectoryStructure(tt.objects)
			if len(dirs) != tt.wantDirs {
				t.Errorf("BuildDirectoryStructure() got %d directories, want %d", len(dirs), tt.wantDirs)
			}

			// Verify root directory exists
			if _, exists := dirs[""]; !exists {
				t.Error("BuildDirectoryStructure() root directory does not exist")
			}
		})
	}
}

func TestBuildDirectoryStructure_FilesAndDirectories(t *testing.T) {
	objects := []Object{
		{Key: "file1.txt", Size: 100},
		{Key: "dir1/file2.txt", Size: 200},
		{Key: "dir1/file3.txt", Size: 300},
		{Key: "dir1/subdir/file4.txt", Size: 400},
	}

	dirs := BuildDirectoryStructure(objects)

	// Check root directory
	root := dirs[""]
	if len(root.Files) != 1 {
		t.Errorf("root directory has %d files, want 1", len(root.Files))
	}
	if len(root.Directories) != 1 {
		t.Errorf("root directory has %d subdirectories, want 1", len(root.Directories))
	}
	if root.Files[0].Name != "file1.txt" {
		t.Errorf("root file name = %s, want file1.txt", root.Files[0].Name)
	}
	if root.Directories[0].Name != "dir1" {
		t.Errorf("root subdirectory name = %s, want dir1", root.Directories[0].Name)
	}

	// Check dir1
	dir1 := dirs["dir1"]
	if len(dir1.Files) != 2 {
		t.Errorf("dir1 has %d files, want 2", len(dir1.Files))
	}
	if len(dir1.Directories) != 1 {
		t.Errorf("dir1 has %d subdirectories, want 1", len(dir1.Directories))
	}

	// Check dir1/subdir
	subdir := dirs["dir1/subdir"]
	if len(subdir.Files) != 1 {
		t.Errorf("dir1/subdir has %d files, want 1", len(subdir.Files))
	}
	if len(subdir.Directories) != 0 {
		t.Errorf("dir1/subdir has %d subdirectories, want 0", len(subdir.Directories))
	}
}

func TestBuildDirectoryStructure_FileSizes(t *testing.T) {
	objects := []Object{
		{Key: "file1.txt", Size: 1234},
		{Key: "dir1/file2.txt", Size: 5678},
	}

	dirs := BuildDirectoryStructure(objects)

	root := dirs[""]
	if root.Files[0].Size != 1234 {
		t.Errorf("file size = %d, want 1234", root.Files[0].Size)
	}

	dir1 := dirs["dir1"]
	if dir1.Files[0].Size != 5678 {
		t.Errorf("file size = %d, want 5678", dir1.Files[0].Size)
	}
}
