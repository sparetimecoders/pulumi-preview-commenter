package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/alecthomas/kong"

	"github.com/sparetimecoders/pulumi-preview-commenter/ci"
	"github.com/sparetimecoders/pulumi-preview-commenter/config"
	"github.com/sparetimecoders/pulumi-preview-commenter/content"
	"github.com/sparetimecoders/pulumi-preview-commenter/vcs"
)

var version = "dev"

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args, os.Getenv); err != nil {
		slog.Default().With("error", err).Error("error")
		os.Exit(1)
	}
}

func run(rootCtx context.Context, out io.Writer, args []string, getEnv func(string) string) error {
	cfg := &config.Config{}
	_ = kong.Parse(cfg, kong.Vars{"version": version})

	var leveler slog.LevelVar

	err := leveler.UnmarshalText([]byte(cfg.LogLevel))

	logger := slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{
		AddSource:   false,
		Level:       leveler.Level(),
		ReplaceAttr: nil,
	}))

	if err != nil {
		logger.With("err", err).Error("Failed to parse log level")
		os.Exit(1)
	}
	slog.SetDefault(logger)

	p, err := vcs.CreateProvider(*cfg)
	if err != nil {
		return err
	}
	c := content.NewProcessor(*cfg, p.MaxCommentSize(), ci.GetJobLink(getEnv))
	comment, err := c.Process()
	if err != nil {
		return err
	}

	if err = ProcessComment(rootCtx, p, *cfg, *comment); err != nil {
		return err
	}
	return nil
}

// ProcessComment contains business logic to create, update and delete comments
// If the comment already exist the content will be updated.
// If there are no differences the comment will be deleted
func ProcessComment(ctx context.Context, provider vcs.Provider, cfg config.Config, diff content.Diff) error {
	logger := slog.Default().With("tagId", cfg.TagId)
	comment, err := findComment(ctx, provider, cfg.TagId)
	if err != nil {
		return err
	}
	if comment != nil {
		logger = logger.With("commentId", comment.Id).With("link", comment.Link)
		if !diff.HasChanges {
			err = provider.DeleteComment(ctx, comment.Id)
			if err != nil {
				logger.With("error", err).Error("delete")
				return err
			}
			logger.Info("Deleted comment because no changes detected")
			return nil
		}
		// if comment exists and there are diff then update existing comment
		_, err = provider.UpdateComment(ctx, comment.Id, diff.Content)
		if err != nil {
			logger.With("error", err).Error("update")
			return err
		}
		logger.Info("Updated comment")
		return nil
	}
	if !diff.HasChanges {
		logger.Info("There is no diff detected. Skip posting diff.")
		return nil
	}
	_, err = provider.CreateComment(ctx, diff.Content)
	if err != nil {
		logger.With("error", err).Error("create")
		return err
	}
	logger.Info("Created comment")
	return nil
}

func findComment(ctx context.Context, ns vcs.Provider, tagId string) (*vcs.Comment, error) {
	logger := slog.Default().With("tagId", tagId)
	comments, err := ns.ListComments(ctx)
	if err != nil {
		return nil, err
	}
	for _, comment := range comments {
		if strings.Contains(comment.Body, getHeaderTagID(tagId)) {
			logger.Debug("Found existing comment")
			return &comment, nil
		}
	}
	return nil, nil
}

func getHeaderTagID(tagId string) string {
	return fmt.Sprintf("%s %s", content.HeaderPrefix, tagId)
}
