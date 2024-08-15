package vcs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

type BitbucketComment struct {
	Content *BitbucketContent `json:"content,omitempty"`
	Id      *int64            `json:"id,omitempty"`
	Links   *BitbucketLinks   `json:"links,omitempty"`
}

type BitbucketContent struct {
	Raw string `json:"raw,omitempty"`
}

type BitbucketComments struct {
	Values []BitbucketComment `json:"values,omitempty"`
}

type BitbucketLinks struct {
	Html BitbucketHtml `json:"html,omitempty"`
}

type BitbucketHtml struct {
	Href string `json:"href,omitempty"`
}

type BitbucketClient struct {
	*http.Client
	BaseURL   *url.URL
	UserAgent string
}

func (c *BitbucketClient) NewRequest(method string, url string, body any) (*http.Request, error) {
	u, err := c.BaseURL.Parse(url)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)

		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
}

func (c *BitbucketClient) Do(ctx context.Context, req *http.Request, v any) (*http.Response, error) {
	resp, err := c.Client.Do(req.WithContext(ctx))
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, fmt.Errorf("error reading response body: %w", err)
	}
	if v != nil {
		decErr := json.Unmarshal(body, v)
		if decErr != nil {
			return resp, fmt.Errorf("error parsing response body to %s: %w", reflect.TypeOf(v), decErr)
		}
	}
	if resp.StatusCode >= 300 {
		return resp, fmt.Errorf("API Error: %s %s", resp.Status, body)
	}
	return resp, nil
}

func addOptions(originalUrl string, opts any) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return originalUrl, nil
	}

	u, err := url.Parse(originalUrl)
	if err != nil {
		return originalUrl, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return originalUrl, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
