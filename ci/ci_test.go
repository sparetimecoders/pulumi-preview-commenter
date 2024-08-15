package ci

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetJobLink(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		expected string
	}{
		{
			name: "none",
		},
		{
			name: "buildkite",
			env: map[string]string{
				"BUILDKITE":           "true",
				"BUILDKITE_BUILD_URL": "url",
			},
			expected: "url",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link := GetJobLink(getEnv(tt.env))
			require.Equal(t, tt.expected, link)
		})
	}
}

func getEnv(env map[string]string) GetEnvFn {
	return func(key string) string {
		return env[key]
	}
}
