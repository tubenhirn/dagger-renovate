package ci

import (
	"dagger.io/dagger"
	"tubenhirn.com/ci/renovate"
)

dagger.#Plan & {
	client: env: {
		GITLAB_TOKEN: dagger.#Secret
		GITHUB_TOKEN: dagger.#Secret
		// repositories is a list of git repositories seperated by ","
		// e.g. "mynamespace/myproject"
		RENOVATE_REPOSITORIES: string
	}

	actions: {
		"renovate": renovate.#Run & {
			repositories: client.env.RENOVATE_REPOSITORIES
			version:      "32.131.1"
			platform:     "gitlab"
			gitlabToken:  client.env.GITLAB_TOKEN
			githubToken:  client.env.GITHUB_TOKEN
		}
	}
}
