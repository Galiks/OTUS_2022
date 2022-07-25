package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	invalidSymbol string = "="
)

type Environment map[string]EnvValue

func (env Environment) Clear() {
	for k := range env {
		delete(env, k)
	}
}

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	environments := make(Environment)

	dirEntry, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, de := range dirEntry {
		if de.IsDir() || !de.Type().IsRegular() || strings.Contains(de.Name(), invalidSymbol) {
			continue
		}
		line, err := readLine(dir, de)
		if err != nil {
			return nil, err
		}
		firstLine := strings.TrimRight(string(bytes.ReplaceAll(line, []byte{0x00}, []byte{'\n'})), "\t ")

		eVal := EnvValue{
			Value:      firstLine,
			NeedRemove: false,
		}
		if firstLine == "" {
			eVal.NeedRemove = true
		}
		environments[de.Name()] = eVal
	}

	return environments, nil
}

func readLine(dir string, de fs.DirEntry) ([]byte, error) {
	f, err := os.OpenFile(filepath.Join(dir, de.Name()), os.O_RDONLY, 0o644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	line, _, err := reader.ReadLine()
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, err
		}
	}
	return line, nil
}
