package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrSourceFileNotFound      = errors.New("source file not found")
	ErrDestinationFileNotFound = errors.New("destination file not found") // destination unknown 🎷
	ErrInvalidCopyParams       = errors.New("invalid copy params")
	ErrTryCopyNothing          = errors.New("offset is equal size of file. you try copy nothing")
	ErrUnsupportedFile         = errors.New("unsupported file")
	ErrOffsetExceedsFileSize   = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 || limit < 0 {
		return ErrInvalidCopyParams
	}
	sourceInfo, err := os.Stat(fromPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrSourceFileNotFound
		}
		return err
	}
	if sourceInfo.Size() == 0 || !sourceInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	if sourceInfo.Size() < offset {
		return ErrOffsetExceedsFileSize
	}
	if sourceInfo.Size() == offset {
		return ErrTryCopyNothing
	}
	source, err := os.OpenFile(fromPath, os.O_RDONLY, 0o644)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrSourceFileNotFound
		}
		return err
	}
	defer source.Close()
	if _, err := os.Stat(toPath); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(toPath), os.ModePerm)
		if err != nil {
			return err
		}
	}
	destination, err := os.OpenFile(toPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return ErrDestinationFileNotFound
	}
	defer destination.Close()

	if offset != 0 {
		_, err = source.Seek(offset, 0)
		if err != nil {
			return err
		}
	}

	var bytesLimit int64 // limit of bytes to read from file
	if limit == 0 {
		bytesLimit = sourceInfo.Size()
	} else {
		bytesLimit = limit
	}

	bar := pb.Full.Start64(bytesLimit - offset)
	defer bar.Finish()
	reader := io.LimitReader(source, bytesLimit)
	barReader := bar.NewProxyReader(reader)

	if _, err := io.Copy(destination, barReader); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	if err = bar.Err(); err != nil {
		return err
	}

	return nil
}
