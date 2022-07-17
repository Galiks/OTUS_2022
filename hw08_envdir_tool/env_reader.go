package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	var (
		invalidSymbol = "="
	)
	environments := make(Environment)

	dirEntry, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, de := range dirEntry {
		if de.IsDir() || !de.Mode().IsRegular() || strings.Contains(de.Name(), invalidSymbol) {
			continue
		}
		f, err := os.OpenFile(filepath.Join(dir, de.Name()), os.O_RDONLY, 0644)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		reader := bufio.NewReader(f)
		line, _, err := reader.ReadLine()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
		}
		firstLine := strings.TrimRight(string(bytes.Replace(line, []byte{0x00}, []byte{'\n'}, -1)), "\t ")

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
