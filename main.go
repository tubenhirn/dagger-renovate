package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/tubenhirn/dagger-ci-modules/v5"
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

	var repositories string
	var secrets = make(map[string]dagger.SecretID)
	if *platform == "github" {
		repositories = os.Getenv("GITHUB_RENOVATE_REPOSITORIES")
		token := os.Getenv("GITHUB_ACCESS_TOKEN")
		renovateTokenId, err := client.SetSecret("GITHUB_ACCESS_TOKEN", token).ID(ctx)
		if err != nil {
			panic(err)
		}
		secrets["RENOVATE_TOKEN"] = renovateTokenId
	} else {
		repositories = os.Getenv("GITLAB_RENOVATE_REPOSITORIES")
		gitlabtoken := os.Getenv("GITLAB_ACCESS_TOKEN")
		renovateTokenId, err := client.SetSecret("GITLAB_ACCESS_TOKEN", gitlabtoken).ID(ctx)
		if err != nil {
			panic(err)
		}
		githubtoken := os.Getenv("GITHUB_COM_TOKEN")
		githubTokenId, err := client.SetSecret("GITHUB_COM_TOKEN", githubtoken).ID(ctx)

		if err != nil {
			panic(err)
		}
		secrets["RENOVATE_TOKEN"] = renovateTokenId
		secrets["GITHUB_COM_TOKEN"] = githubTokenId
	}

	options := cimodules.RenovateOpts{
		Platform:           *platform,
		Autodiscover:       false,
		AutodiscoverFilter: "",
		Repositories:       repositories,
		Env:                map[string]string{},
		Secret:             secrets,
		LogLevel:           "warn",
	}

	err = cimodules.Renovate(ctx, *client, options)
	if err != nil {
		panic(err)
	}
}
