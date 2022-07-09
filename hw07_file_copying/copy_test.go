package main

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	// Place your code here.
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	destinationPath := filepath.Join(dir, "testoutput/foo/foo/foo/output.txt")
	sourcePath := filepath.Join(dir, "testdata/input.txt")
	err = Copy(sourcePath, destinationPath, 0, 6000)
	require.Nil(t, err)
	// err = os.Remove(destinationPath)
	// require.Nil(t, err)
}
