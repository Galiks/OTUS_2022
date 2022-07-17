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
}
