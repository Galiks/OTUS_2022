package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {

	t.Run("full copy", func(t *testing.T) {
		dir, err := os.Getwd()
		require.Nil(t, err)
		destinationPath := filepath.Join(dir, "testoutput/foo/foo/foo/output.txt")
		sourcePath := filepath.Join(dir, "testdata/input.txt")
		err = Copy(sourcePath, destinationPath, 0, 0)
		require.Nil(t, err)
		ds, err := os.Stat(destinationPath)
		require.Nil(t, err)
		ss, err := os.Stat(sourcePath)
		require.Nil(t, err)
		require.Equal(t, ds.Size(), ss.Size())
		err = os.Remove(destinationPath)
		require.Nil(t, err)
	})

	t.Run("copy with limit", func(t *testing.T) {
		dir, err := os.Getwd()
		require.Nil(t, err)
		destinationPath := filepath.Join(dir, "testoutput/foo/foo/foo/output.txt")
		sourcePath := filepath.Join(dir, "testdata/input.txt")
		var destinationFileSize int64 = 1000
		err = Copy(sourcePath, destinationPath, 0, destinationFileSize)
		require.Nil(t, err)
		ds, err := os.Stat(destinationPath)
		require.Nil(t, err)
		ss, err := os.Stat(sourcePath)
		require.Nil(t, err)
		require.Equal(t, ds.Size(), destinationFileSize)
		require.NotEqual(t, ds.Size(), ss.Size())
		// dest, err := os.OpenFile(destinationPath, os.O_RDONLY, 0644)
		// require.Nil(t, err)
		err = os.Remove(destinationPath)
		require.Nil(t, err)
	})

	t.Run("copy with offset", func(t *testing.T) {

	})

}
