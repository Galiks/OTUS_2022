package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	env := make(Environment)

	dir, err := os.Getwd()
	require.Nil(t, err)
	testCmd := []string{}
	switch runtime.GOOS {
	case "windows":
		testCmd = []string{"bash", filepath.Join(dir, "/testdata/echo.sh"), "arg1=1", "arg2=2"}
	case "linux":
		testCmd = []string{filepath.Join(dir, "/testdata/echo.sh"), "arg1=1", "arg2=2"}
	}

	t.Run("simple test", func(t *testing.T) {
		env.Clear()
		expectedCode := 0
		cmdCode := RunCmd(testCmd, env)
		require.Equal(t, expectedCode, cmdCode)
	})

	t.Run("test on bash", func(t *testing.T) {
		env.Clear()
		expectedCode := 127
		fmt.Printf("expectedCode: %v\n", expectedCode)
		testCmd := []string{"/bin/bash", "xxx"}
		cmdCode := RunCmd(testCmd, env)
		require.Equal(t, expectedCode, cmdCode)
	})

	t.Run("test on error 'ErrEmptyCmd'", func(t *testing.T) {
		env.Clear()
		expectedCode := 1
		testCmd := make([]string, 0)
		cmdCode := RunCmd(testCmd, env)
		require.Equal(t, expectedCode, cmdCode)
	})

	t.Run("test on set environment", func(t *testing.T) {
		env.Clear()
		expectedCode := 0
		envName := "TEST_ENV_NAME"
		envValue := "TEST_ENV_VALUE"
		env[envName] = EnvValue{envValue, false}
		cmdCode := RunCmd(testCmd, env)
		require.Equal(t, expectedCode, cmdCode)
		val, ok := os.LookupEnv(envName)
		require.True(t, ok)
		require.Equal(t, val, envValue)
	})

	t.Run("test on unset environment", func(t *testing.T) {
		env.Clear()
		expectedCode := 0
		envName := "TEST_ENV_NAME"
		envValueForUnset := "TEST_ENV_VALUE"
		envValue := ""
		err := os.Setenv(envName, envValueForUnset)
		require.Nil(t, err)
		env[envName] = EnvValue{envValue, true}
		cmdCode := RunCmd(testCmd, env)
		require.Equal(t, expectedCode, cmdCode)
		val, ok := os.LookupEnv(envName)
		require.False(t, ok)
		require.Equal(t, val, envValue)
	})
}
