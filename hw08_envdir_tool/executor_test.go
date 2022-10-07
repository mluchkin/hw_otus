package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("0", func(t *testing.T) {
		cmd := []string{
			"/bin/bash",
			"testdata/echo.sh",
		}
		returnCode := RunCmd(cmd, Environment{})

		require.Equal(t, 0, returnCode)
	})

	t.Run("127", func(t *testing.T) {
		cmd := []string{
			"/bin/bash",
			"test2",
		}
		returnCode := RunCmd(cmd, Environment{})

		require.Equal(t, 127, returnCode)
	})
}
