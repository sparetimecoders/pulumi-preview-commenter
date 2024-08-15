package vcs

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/sparetimecoders/pulumi-preview-commenter/config"
)

const (
	BitbucketMaxCommentLength = 32768
	BitbucketDefaultBaseURL   = "https://api.bitbucket.org/2.0/"
)

type BitbucketProvider struct {
	*BitbucketClient
	config config.Config
}

func NewBitbucketProvider(config config.Config) *BitbucketProvider {
	u, _ := url.Parse(BitbucketDefaultBaseURL)
	c := http.DefaultClient
	c.Transport = &transport{
		underlyingTransport: http.DefaultTransport,
		token:               config.AuthToken,
	}
	return &BitbucketProvider{
		BitbucketClient: &BitbucketClient{
			Client:    http.DefaultClient,
			BaseURL:   u,
			UserAgent: userAgent,
		},
		config: config,
	}
}

func (b BitbucketProvider) MaxCommentSize() int {
	return BitbucketMaxCommentLength
}

func (b BitbucketProvider) CreateComment(ctx context.Context, content string) (*Comment, error) {
	u := fmt.Sprintf("repositories/%s/%s/pullrequests/%d/comments", b.config.RepoOwner, b.config.RepoName, b.config.PullRequestId)
	req, err := b.NewRequest(http.MethodPost, u, &BitbucketComment{Content: &BitbucketContent{Raw: content}})
	if err != nil {
		return nil, err
	}
	comment := &BitbucketComment{}
	if _, err = b.Do(ctx, req, comment); err != nil {
		return nil, err
	}
	return comment.convert(), nil
}

func (b BitbucketProvider) UpdateComment(ctx context.Context, id int64, content string) (*Comment, error) {
	u := fmt.Sprintf("repositories/%s/%s/pullrequests/%d/comments/%d", b.config.RepoOwner, b.config.RepoName, b.config.PullRequestId, id)
	req, err := b.NewRequest(http.MethodPut, u, &BitbucketComment{Content: &BitbucketContent{Raw: content}})
	if err != nil {
		return nil, err
	}
	comment := &BitbucketComment{}
	if _, err = b.Do(ctx, req, comment); err != nil {
		return nil, err
	}
	return comment.convert(), nil
}

func (b BitbucketProvider) DeleteComment(ctx context.Context, id int64) error {
	u := fmt.Sprintf("repositories/%s/%s/pullrequests/%d/comments/%d", b.config.RepoOwner, b.config.RepoName, b.config.PullRequestId, id)
	req, err := b.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return err
	}
	if _, err = b.Do(ctx, req, nil); err != nil {
		return err
	}
	return nil
}

func (b BitbucketProvider) ListComments(ctx context.Context) ([]Comment, error) {
	u := fmt.Sprintf("repositories/%s/%s/pullrequests/%d/comments", b.config.RepoOwner, b.config.RepoName, b.config.PullRequestId)
	u, err := addOptions(u, struct {
		Query      string `url:"q,omitempty"`
		Fields     string `url:"fields,omitempty"`
		PageLength int    `url:"pagelen,omitempty"`
	}{
		// filter out deleted
		Query: "deleted=false",
		// just fetch relevant fields
		Fields: "values.id,values.content.raw,values.links.html.href",
		// set maximum page length, don't  bother with paging
		PageLength: 100,
	})
	if err != nil {
		return nil, err
	}
	req, err := b.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	comments := &BitbucketComments{}
	if _, err = b.Do(ctx, req, comments); err != nil {
		return nil, err
	}
	var result []Comment

	if comments != nil && len(comments.Values) > 0 {
		for _, c := range comments.Values {
			result = append(result, *c.convert())
		}
	}
	return result, nil
}

func (c *BitbucketComment) convert() *Comment {
	comment := &Comment{}
	if c == nil {
		return comment
	}
	if c.Id != nil {
		comment.Id = *c.Id
	}
	if c.Content != nil {
		comment.Body = c.Content.Raw
	}
	if c.Links != nil {
		comment.Link = c.Links.Html.Href
	}
	return comment
}

type transport struct {
	underlyingTransport http.RoundTripper
	token               string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))
	req.Header.Add("Accept", "application/json")

	return t.underlyingTransport.RoundTrip(req)
}
