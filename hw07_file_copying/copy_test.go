package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	dir, err := os.Getwd()
	require.Nil(t, err)
	destinationPath := filepath.Join(dir, "testoutput/foo/foo/foo/output.txt")
	sourcePath := filepath.Join(dir, "testdata/input.txt")

	t.Run("full copy", func(t *testing.T) {
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
		var fileLimit int64 = 1000
		err = Copy(sourcePath, destinationPath, 0, fileLimit)
		require.Nil(t, err)
		ds, err := os.Stat(destinationPath)
		require.Nil(t, err)
		ss, err := os.Stat(sourcePath)
		require.Nil(t, err)
		require.Equal(t, ds.Size(), fileLimit)
		require.NotEqual(t, ds.Size(), ss.Size())
		err = os.Remove(destinationPath)
		require.Nil(t, err)
	})

	t.Run("copy with offset", func(t *testing.T) {
		var fileOffset int64 = 1000
		err = Copy(sourcePath, destinationPath, fileOffset, 0)
		require.Nil(t, err)
		ds, err := os.Stat(destinationPath)
		require.Nil(t, err)
		ss, err := os.Stat(sourcePath)
		require.Nil(t, err)
		require.Equal(t, ss.Size()-fileOffset, ds.Size())
		require.NotEqual(t, ds.Size(), ss.Size())
		err = os.Remove(destinationPath)
		require.Nil(t, err)
	})

	t.Run("copy with offset and limit", func(t *testing.T) {
		var fileLimit int64 = 1000
		var fileOffset int64 = 1000
		err = Copy(sourcePath, destinationPath, fileOffset, fileLimit)
		require.Nil(t, err)
		ds, err := os.Stat(destinationPath)
		require.Nil(t, err)
		ss, err := os.Stat(sourcePath)
		require.Nil(t, err)
		require.Equal(t, fileLimit, ds.Size())
		require.NotEqual(t, ds.Size(), ss.Size())
		err = os.Remove(destinationPath)
		require.Nil(t, err)
	})

	t.Run("test on error `ErrSourceFileNotFound`", func(t *testing.T) {
		var notExistedFile string
		err = Copy(notExistedFile, destinationPath, 0, 0)
		require.ErrorIs(t, err, ErrSourceFileNotFound)
	})

	t.Run("test on error `ErrDestinationFileNotFound`", func(t *testing.T) {
		var notExistedFile string
		err = Copy(sourcePath, notExistedFile, 0, 0)
		require.ErrorIs(t, err, ErrDestinationFileNotFound)
	})

	t.Run("test on error `ErrInvalidCopyParams`", func(t *testing.T) {
		err = Copy(sourcePath, destinationPath, -1, -1)
		require.ErrorIs(t, err, ErrInvalidCopyParams)

		err = Copy(sourcePath, destinationPath, 1, -1)
		require.ErrorIs(t, err, ErrInvalidCopyParams)

		err = Copy(sourcePath, destinationPath, -1, 1)
		require.ErrorIs(t, err, ErrInvalidCopyParams)
	})

	t.Run("test on error `ErrTryCopyNothing`", func(t *testing.T) {
		ss, err := os.Stat(sourcePath)
		require.Nil(t, err)
		err = Copy(sourcePath, destinationPath, ss.Size(), 0)
		require.ErrorIs(t, err, ErrTryCopyNothing)
	})

	t.Run("test on error `ErrUnsupportedFile`", func(t *testing.T) {
		newSourceFile := "source_test.qwe"
		file, err := os.Create(newSourceFile)
		require.Nil(t, err)
		err = Copy(newSourceFile, destinationPath, 0, 0)
		require.ErrorIs(t, err, ErrUnsupportedFile)
		err = file.Close()
		require.Nil(t, err)
		err = os.Remove(filepath.Join(dir, newSourceFile))
		require.Nil(t, err)
	})

	t.Run("test on error `ErrOffsetExceedsFileSize`", func(t *testing.T) {
		var fileOffset int64 = 1000
		ss, err := os.Stat(sourcePath)
		require.Nil(t, err)
		err = Copy(sourcePath, destinationPath, ss.Size()+fileOffset, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})
}
