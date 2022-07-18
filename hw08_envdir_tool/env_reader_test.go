package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("simple test", func(t *testing.T) {
		expectedLen := 5
		dir, err := os.Getwd()
		require.Nil(t, err)
		dirPath := filepath.Join(dir, "/testdata/env")
		env, err := ReadDir(dirPath)
		require.Nil(t, err)
		require.Equal(t, expectedLen, len(env))
	})

	t.Run("invalid dir", func(t *testing.T) {
		env, err := ReadDir("")
		require.NotNil(t, err)
		require.Nil(t, env)
	})

	t.Run("invalid symbol", func(t *testing.T) {
		dir, err := os.Getwd()
		require.Nil(t, err)
		pathTemp, err := os.MkdirTemp(dir, "test_invalid_symbol")
		require.Nil(t, err)
		tempFileName := "INVALID=FILE"
		expectedEnvLen := 0
		file, err := os.CreateTemp(pathTemp, tempFileName)
		require.Nil(t, err)
		f, err := os.Stat(file.Name())
		require.Nil(t, err)
		require.Contains(t, f.Name(), tempFileName)
		env, err := ReadDir(pathTemp)
		require.Nil(t, err)
		require.Equal(t, expectedEnvLen, len(env))
		err = file.Close()
		require.Nil(t, err)
		err = os.Remove(file.Name())
		require.Nil(t, err)
		err = os.RemoveAll(pathTemp)
		require.Nil(t, err)
	})
}
