package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("Test environment extraction", func(t *testing.T) {
		env, err := ReadDir("testMy")
		require.NoError(t, err, "actual error %q", err)

		expected := Environment{
			"ONE":   EnvValue{Value: "One"},
			"TWO":   EnvValue{Value: "Two"},
			"EMPTY": EnvValue{NeedRemove: true},
		}

		require.Equal(t, env, expected, "actual: %v expected: %v", env, expected)
	})
}
