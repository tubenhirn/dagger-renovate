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

	renovate := client.Container().From("renovate/renovate:" + renovateVersion)
	renovate = renovate.WithSecretVariable("RENOVATE_TOKEN", client.Secret(accessToken))
	renovate = renovate.WithSecretVariable("GITHUB_COM_TOKEN", client.Secret(githubToken))
	renovate = renovate.WithEnvVariable("RENOVATE_PLATFORM", platform)
	renovate = renovate.WithEnvVariable("RENOVATE_EXTENDS", "github>whitesource/merge-confidence:beta")
	renovate = renovate.WithEnvVariable("RENOVATE_REQUIRE_CONFIG", "true")
	renovate = renovate.WithEnvVariable("RENOVATE_GIT_AUTHOR", "Renovate Bot <bot@renovateapp.com>")
	renovate = renovate.WithEnvVariable("RENOVATE_PIN_DIGEST", "true")
	renovate = renovate.WithEnvVariable("RENOVATE_DEPENDENCY_DASHBOARD", "false")
	renovate = renovate.WithEnvVariable("RENOVATE_LABELS", "renovate")
	renovate = renovate.WithEnvVariable("RENOVATE_AUTODISCOVER", autodiscover)
	renovate = renovate.WithEnvVariable("RENOVATE_AUTODISCOVER_FILTER", autodiscoverFilter)
	renovate = renovate.WithEnvVariable("RENOVATE_REPOSITORIES", repositories)
	renovate = renovate.WithEnvVariable("LOG_LEVEL", logLevel)

	renovate.Exec().Stdout(ctx)

	return nil
}
