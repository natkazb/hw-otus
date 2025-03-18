package app

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/storage" //nolint
	"github.com/stretchr/testify/require"
)

type LoggerTest struct{}

func (l LoggerTest) Error(msg string) {
	fmt.Fprintf(os.Stdout, "[Error] %s\n", msg)
}

func (l LoggerTest) Info(msg string) {
	fmt.Fprintf(os.Stdout, "[Info] %s\n", msg)
}

func (l LoggerTest) Debug(msg string) {
	fmt.Fprintf(os.Stdout, "[Debug] %s\n", msg)
}

func (l LoggerTest) Warn(msg string) {
	fmt.Fprintf(os.Stdout, "[Warn] %s\n", msg)
}

func TestApp(t *testing.T) {
	l := &LoggerTest{}
	appl := New(l, nil)
	now := time.Now()
	err := appl.UpdateEvent(nil, 1, "", "", now.Add(time.Hour), now.Add(-time.Hour)) //nolint
	require.Error(t, err)
	require.ErrorIs(t, err, storage.ErrDates)
}
