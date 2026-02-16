package main

import (
	"io"
	"os"

	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

var isWindows bool = false

func readWindowsFile(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	if utf8.Valid(data) {
		return string(data), nil
	}

	isWindows = true

	decoded, err := charmap.Windows1251.NewDecoder().Bytes(data)
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}

func saveWindowsFile(filename string, text string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if !isWindows {
		_, err = f.Write([]byte(text))
		return err
	}

	encoder := charmap.Windows1251.NewEncoder()
	writer := transform.NewWriter(f, encoder)

	_, err = writer.Write([]byte(text))
	if err != nil {
		return err
	}

	return writer.Close()
}
