package content

import (
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sparetimecoders/pulumi-preview-commenter/config"
)

func TestProcess(t *testing.T) {
	tests := []struct {
		name             string
		jobLink          string
		maxCommentLength int
		expectTruncated  bool
	}{
		{
			name:             "truncated",
			maxCommentLength: 200,
			expectTruncated:  true,
		},
		{
			name:             "unchanged",
			jobLink:          "href://abc.com",
			maxCommentLength: math.MaxInt,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewProcessor(config.Config{
				TagId: "tagId",
				File:  "testdata/unchanged.txt",
			}, tt.maxCommentLength, tt.jobLink)
			content, err := p.Process()
			require.NoError(t, err)
			require.Contains(t, content.Content, "**pulumi diff for tagId")
			if tt.jobLink != "" {
				require.Contains(t, content.Content, fmt.Sprintf("[Job Link](%s)", tt.jobLink))
			}
			if tt.expectTruncated {
				require.Contains(t, content.Content, "**Warning**: Truncated")
			}
		})
	}
}

func Test_diffHasChanges(t *testing.T) {
	tests := []struct {
		name    string
		changes bool
	}{
		{
			name:    "creates_and_updates",
			changes: true,
		},
		{
			name:    "creates",
			changes: true,
		},
		{
			name:    "deletions",
			changes: true,
		},
		{
			name: "unchanged",
		},
		{
			name:    "updates",
			changes: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := fmt.Sprintf("testdata/%s.txt", tt.name)
			data, err := os.ReadFile(file)
			require.NoError(t, err)
			require.Equal(t, tt.changes, diffHasChanges(data))
		})
	}
}
