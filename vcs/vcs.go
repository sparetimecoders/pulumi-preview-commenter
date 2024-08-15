package vcs

import (
	"context"
	"fmt"

	"github.com/sparetimecoders/pulumi-preview-commenter/config"
)

const (
	userAgent = "pulumi-preview-commenter"
)

type Comment struct {
	Id   int64
	Body string
	Link string
}

// Provider interface for public methods actions required for handling PullRequest comments
type Provider interface {
	CreateComment(ctx context.Context, content string) (*Comment, error)
	UpdateComment(ctx context.Context, id int64, content string) (*Comment, error)
	DeleteComment(ctx context.Context, id int64) error
	ListComments(ctx context.Context) ([]Comment, error)
	MaxCommentSize() int
}

func CreateProvider(c config.Config) (Provider, error) {
	switch c.Vcs {
	case config.VcsBitbucket:
		return NewBitbucketProvider(c), nil
	default:
		return nil, fmt.Errorf("unspported Version Control System: %s", c.Vcs)
	}
}
