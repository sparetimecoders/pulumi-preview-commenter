package config

import "github.com/alecthomas/kong"

const (
	VcsBitbucket = "bitbucket"
)

// TODO Determine config from CI Env variables?
type Config struct {
	Vcs           string `name:"vcs" env:"VERSION_CONTROL_SYSTEM" help:"The Version Control System to use" enum:"bitbucket" required:""`
	AuthToken     string `name:"token" env:"TOKEN" help:"The token used to authenticate to the VCS" required:""`
	RepoName      string `name:"repository" env:"REPO_NAME" help:"The name of the VCS repository" required:""`
	RepoOwner     string `name:"owner" env:"REPO_OWNER" help:"The owner of the repository" required:""`
	PullRequestId int    `name:"pull-request-id" env:"PULL_REQUEST_ID" help:"The pull request ID used to create/update/delete the comment" required:""`
	TagId         string `name:"tag-id" env:"TAG_ID" help:"Unique identifier for the comment" required:""`
	File          string `name:"file" env:"FILE" help:"The path to the pulumi preview diff file" default:"preview"`
	LogLevel      string `name:"log-level" env:"LOG_LEVEL" help:"The level of logging to use" enum:"debug,info,warn,error" default:"info"`
	Version       kong.VersionFlag
}
