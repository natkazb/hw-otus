package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyErrUnsupportedFile(t *testing.T) {
	dst, err := os.CreateTemp(os.TempDir(), "something")
	defer os.Remove(dst.Name())
	require.NoError(t, err)
	err = Copy("/dev/null", dst.Name(), 0, 1000)
	require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual error %q", err)
}

func TestCopyErrOffsetExceedsFileSize(t *testing.T) {
	src, err := os.CreateTemp(os.TempDir(), "src")
	defer os.Remove(src.Name())
	require.NoError(t, err)

	dst, err := os.CreateTemp(os.TempDir(), "dst")
	defer os.Remove(dst.Name())
	require.NoError(t, err)

	err = Copy(src.Name(), dst.Name(), 50, 200)
	require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual error %q", err)
}

func TestCopyOk(t *testing.T) {
	src, err := os.CreateTemp(os.TempDir(), "src")
	defer os.Remove(src.Name())
	require.NoError(t, err)

	data := "something"
	bytes, err := src.WriteString(data)
	require.NoError(t, err)

	dst, err := os.CreateTemp(os.TempDir(), "dst")
	defer os.Remove(dst.Name())
	require.NoError(t, err)

	err = Copy(src.Name(), dst.Name(), 0, int64(bytes))
	require.NoError(t, err)
	output := make([]byte, bytes)
	_, _ = dst.Read(output)
	require.Equal(t, data, string(output))
}
