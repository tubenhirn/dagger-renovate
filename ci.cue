package ci

import (
	"dagger.io/dagger"
	"tubenhirn.com/ci/renovate"
)

dagger.#Plan & {
	client: env: {
		ACCESS_TOKEN: dagger.#Secret
		GITHUB_TOKEN: dagger.#Secret
		// repositories is a list of git repositories seperated by ","
		// e.g. "mynamespace/myproject"
		RENOVATE_REPOSITORIES: string
	}

	actions: {
		"renovate-gitlab": renovate.#Run & {
			repositories: client.env.RENOVATE_REPOSITORIES
			version:      "32.159.0"
			platform:     "gitlab"
			accessToken:  client.env.ACCESS_TOKEN
			githubToken:  client.env.GITHUB_TOKEN
		}
		"renovate-github": renovate.#Run & {
			repositories: client.env.RENOVATE_REPOSITORIES
			version:      "32.159.0"
			platform:     "github"
			accessToken:  client.env.ACCESS_TOKEN
			githubToken:  client.env.GITHUB_TOKEN
		}

	}
}
