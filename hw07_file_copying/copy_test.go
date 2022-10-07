package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	fromFile = "testdata/input.txt"
	toFile   = "out.txt"
)

func TestCopy(t *testing.T) {
	t.Run("test no errors", func(t *testing.T) {
		offset := int64(0)
		limit := int64(0)
		err := Copy(fromFile, toFile, offset, limit)

		require.Equal(t, err, nil)
	})

	t.Run("test offset exceed file size", func(t *testing.T) {
		offset := int64(1024 * 1024)
		limit := int64(0)
		err := Copy(fromFile, toFile, offset, limit)

		require.Equal(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("test limit works just fine", func(t *testing.T) {
		offset := int64(0)
		limit := int64(2)
		Copy(fromFile, toFile, offset, limit)
		newFile, _ := os.Open(toFile)
		defer newFile.Close()
		newFileContent := make([]byte, limit)
		newFile.Read(newFileContent)

		require.Len(t, newFileContent, 2)
		require.Equal(t, string(newFileContent), "Go")
	})
}
