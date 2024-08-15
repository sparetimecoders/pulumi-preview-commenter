package content

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/sparetimecoders/pulumi-preview-commenter/config"
)

const (
	// HeaderPrefix default prefix for comment message
	HeaderPrefix = "pulumi diff for"
)

type Processor struct {
	file             string
	tagId            string
	maxContentLength int
	jobLink          string
}

type Diff struct {
	HasChanges bool
	Content    string
}

func NewProcessor(cfg config.Config, maxContentSize int, jobLink string) *Processor {
	if jobLink != "" {
		jobLink = fmt.Sprintf("[Job Link](%s)", jobLink)
	}
	return &Processor{
		tagId:            cfg.TagId,
		file:             cfg.File,
		maxContentLength: maxContentSize,
		jobLink:          jobLink,
	}
}

func (p Processor) Process() (*Diff, error) {
	content, err := os.ReadFile(p.file)
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New("commentTemplate").Funcs(sprig.FuncMap()).Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}

	stringWriter := bytes.NewBufferString("")
	data := commentTemplate{
		TagID:        p.tagId,
		Diff:         string(content),
		JobLink:      p.jobLink,
		Backticks:    "```",
		HeaderPrefix: HeaderPrefix,
	}
	err = tmpl.Execute(stringWriter, data)
	if err != nil {
		return nil, err
	}

	res := &Diff{
		HasChanges: diffHasChanges(content),
	}

	runes := bytes.Runes(stringWriter.Bytes())
	maxContentLength := p.maxContentLength
	truncatedCommentSuffix := fmt.Sprintf("\n%s\n\n**Warning**: Truncated output as length greater than max comment size.", data.Backticks)
	maxContentLength = maxContentLength - len([]rune(truncatedCommentSuffix))
	if len(runes) > maxContentLength {
		truncated := string(runes[:maxContentLength])
		truncated += truncatedCommentSuffix
		res.Content = truncated
	} else {
		res.Content = string(runes)
	}
	return res, nil
}

type commentTemplate struct {
	TagID        string
	Diff         string
	JobLink      string
	Backticks    string
	HeaderPrefix string
}

var defaultTemplate = `
**{{ .HeaderPrefix }} {{ .TagID }}** {{ .JobLink }}
{{ .Backticks }}
{{ .Diff }}
{{ .Backticks }}
`

func diffHasChanges(s []byte) bool {
	return regexp.MustCompile(`(?m)Resources:(\r)?\n.*[+-~]{1} (\d+) to.*`).Match(s)
}
