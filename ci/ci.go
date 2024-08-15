package ci

import (
	"fmt"
	"log/slog"
)

type GetEnvFn = func(string) string

func GetJobLink(getEnv GetEnvFn) string {
	var jobLink string
	if getEnv("BUILDKITE") == "true" {
		jobLink = getEnv("BUILDKITE_BUILD_URL")
	}
	slog.Info(fmt.Sprintf("found job link: %s", jobLink))
	return jobLink
}
