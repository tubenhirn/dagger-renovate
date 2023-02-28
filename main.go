package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/tubenhirn/dagger-ci-modules/v2/renovate"
)

func main() {
	ctx := context.Background()

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}

	defer client.Close()

	// parse the input flag `platform` to decide where we wanna run the renovate
	platform := flag.String("platform", "github", "the string of the platform to run renovate on.")
	flag.Parse()

	fmt.Println("running renovate on " + *platform)

	var renovateTokenId dagger.SecretID
	var repositories string
	if *platform == "github" {
	repositories = 	os.Getenv("GITHUB_RENOVATE_REPOSITORIES")
		renovateTokenId, err = client.Host().EnvVariable("GITHUB_ACCESS_TOKEN").Secret().ID(ctx)
		if err != nil {
			panic(err)
		}
	} else {
		repositories = os.Getenv("GITLAB_RENOVATE_REPOSITORIES")
		renovateTokenId, err = client.Host().EnvVariable("GITLAB_ACCESS_TOKEN").Secret().ID(ctx)
		if err != nil {
			panic(err)
		}
	}

	githubTokenId, err := client.Host().EnvVariable("GITHUB_COM_TOKEN").Secret().ID(ctx)
	if err != nil {
		panic(err)
	}

	options := renovate.RenovateOpts{
		Platform:           *platform,
		Autodiscover:       false,
		AutodiscoverFilter: "",
		Repositories:       repositories,
		Env:                map[string]string{},
		Secret: map[string]dagger.SecretID{
			"RENOVATE_TOKEN": renovateTokenId,
			"GITHUB_COM_TOKEN":   githubTokenId,
		},
		LogLevel: "debug",
	}

	renovate.Renovate(ctx, *client, options)
}
