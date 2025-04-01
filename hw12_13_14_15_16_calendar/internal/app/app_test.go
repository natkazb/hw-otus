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
	event := storage.Event{
		ID:          1,
		Title:       "Title",
		Description: "Description",
		StartDate:   storage.EventDateTime(now.Add(time.Hour)),
		EndDate:     storage.EventDateTime(now.Add(-time.Hour)),
	}
	err := appl.UpdateEvent(event)
	require.Error(t, err)
	require.ErrorIs(t, err, storage.ErrDates)
}
