package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	pb "github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFile, err := openFile(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	fromFileSize, err := checkOffset(fromFile, offset)
	if err != nil {
		return err
	}

	toFile, err := createFile(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	err = copyFile(toFile, fromFile, limit, fromFileSize)
	if err != nil {
		return err
	}

	return nil
}

func initBar(limit, fromFileSize int64) *pb.ProgressBar {
	barLength := fromFileSize
	if limit > 0 {
		barLength = limit
	}
	reader := io.LimitReader(rand.Reader, barLength)
	writer := ioutil.Discard
	bar := pb.Full.Start64(barLength)
	barReader := bar.NewProxyReader(reader)
	io.Copy(writer, barReader)

	return bar
}

func copyFile(toFile, fromFile *os.File, limit, fromFileSize int64) error {
	bar := initBar(limit, fromFileSize)
	if limit > 0 {
		_, err := io.CopyN(toFile, fromFile, limit)
		if err != nil {
			return err
		}
	} else {
		_, err := io.Copy(toFile, fromFile)
		if err != nil {
			return err
		}
	}
	bar.Finish()
	fmt.Println("The file was copied.")

	return nil
}

func checkOffset(fromFile *os.File, offset int64) (int64, error) {
	var err error
	fromFileStat, _ := fromFile.Stat()
	fromFileSize := fromFileStat.Size()
	if offset > 0 {
		if offset > fromFileSize {
			err = ErrOffsetExceedsFileSize
		}
		fromFile.Seek(offset, io.SeekStart)
	}

	return fromFileSize, err
}

func openFile(path string) (*os.File, error) {
	fromFile, err := os.Open(path)
	if err != nil {
		err = ErrUnsupportedFile
	}
	return fromFile, err
}

func createFile(path string) (*os.File, error) {
	toFile, err := os.Create(path)

	return toFile, err
}
