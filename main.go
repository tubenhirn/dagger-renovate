package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	if err := renovate(context.Background()); err != nil {
		fmt.Println(err)
	}
}

func renovate(ctx context.Context) error {
	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}

	defer client.Close()

	var accessToken dagger.SecretID
	var githubToken dagger.SecretID
	var platform string
	var repositories string

	renovateVersion := "latest"
	autodiscover := "false"
	logLevel := "debug"
	autodiscoverFilter := ""

	platform, err = client.Host().EnvVariable("PLATFORM").Value(ctx)
	if err != nil {
		panic(err)
	}

	if platform == "github" {
		accessToken, err = client.Host().EnvVariable("GITHUB_ACCESS_TOKEN").Secret().ID(ctx)
		repositories, err = client.Host().EnvVariable("RENOVATE_REPOSITORIES_GITHUB").Value(ctx)
		if err != nil {
			panic(err)
		}
	} else {
		accessToken, err = client.Host().EnvVariable("GITLAB_ACCESS_TOKEN").Secret().ID(ctx)
		repositories, err = client.Host().EnvVariable("RENOVATE_REPOSITORIES_GITLAB").Value(ctx)
		if err != nil {
			panic(err)
		}
	}

	githubToken, err = client.Host().EnvVariable("GITHUB_TOKEN").Secret().ID(ctx)
	if err != nil {
		panic(err)
	}

	golang := client.Container().From("renovate/renovate:" + renovateVersion)
	golang = golang.WithSecretVariable("RENOVATE_TOKEN", client.Secret(accessToken))
	golang = golang.WithSecretVariable("GITHUB_COM_TOKEN", client.Secret(githubToken))
	golang = golang.WithEnvVariable("RENOVATE_PLATFORM", platform)
	golang = golang.WithEnvVariable("RENOVATE_EXTENDS", "github>whitesource/merge-confidence:beta")
	golang = golang.WithEnvVariable("RENOVATE_REQUIRE_CONFIG", "true")
	golang = golang.WithEnvVariable("RENOVATE_GIT_AUTHOR", "Renovate Bot <bot@renovateapp.com>")
	golang = golang.WithEnvVariable("RENOVATE_PIN_DIGEST", "true")
	golang = golang.WithEnvVariable("RENOVATE_DEPENDENCY_DASHBOARD", "false")
	golang = golang.WithEnvVariable("RENOVATE_LABELS", "renovate")
	golang = golang.WithEnvVariable("RENOVATE_AUTODISCOVER", autodiscover)
	golang = golang.WithEnvVariable("RENOVATE_AUTODISCOVER_FILTER", autodiscoverFilter)
	golang = golang.WithEnvVariable("RENOVATE_REPOSITORIES", repositories)
	golang = golang.WithEnvVariable("LOG_LEVEL", logLevel)

	golang.Exec().Stdout(ctx)

	return nil
}
