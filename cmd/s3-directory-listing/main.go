package main

import (
	"os"

	"github.com/teamsbc/s3-directory-listing/cmd/s3-directory-listing/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
