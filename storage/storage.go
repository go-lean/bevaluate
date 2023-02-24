package storage

import (
	"bytes"
	"fmt"
	"io"
)

type (
	FileReadOpener interface {
		OpenRead(path string) (io.ReadCloser, error)
	}

	FileCreateOpener interface {
		OpenCreate(path string) (io.WriteCloser, error)
	}
)

func CreateFileWithText(path, text string, opener FileCreateOpener) error {
	file, errOpen := opener.OpenCreate(path)
	if errOpen != nil {
		return fmt.Errorf("could not create file: %w", errOpen)
	}

	defer func() {
		if errClose := file.Close(); errClose != nil {
			fmt.Printf("could not close file: %v\n", errClose)
		}
	}()

	_, errWrite := io.WriteString(file, text)
	if errWrite != nil {
		return fmt.Errorf("could not write to file: %w", errWrite)
	}

	return nil
}

func ReadModuleName(path string, opener FileReadOpener) (string, error) {
	file, errOpen := opener.OpenRead(path)
	if errOpen != nil {
		return "", fmt.Errorf("could not open file: %w", errOpen)
	}

	defer func() {
		if errClose := file.Close(); errClose != nil {
			fmt.Printf("could not close file: %v", errClose)
		}
	}()

	data, errRead := io.ReadAll(file)
	if errRead != nil {
		return "", fmt.Errorf("could not read from file: %w", errRead)
	}

	i := bytes.IndexByte(data, '\n')
	if i < 0 {
		return string(data[7:]), nil
	}

	return string(data[7:i]), nil
}
