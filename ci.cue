package ci

import (
	"dagger.io/dagger"
	"tubenhirn.com/ci/renovate"
)

dagger.#Plan & {
	client: env: {
		GITLAB_TOKEN: dagger.#Secret
		GITHUB_TOKEN: dagger.#Secret
	}

	actions: {
		"renovate": renovate.#Run & {
			repositories: "jstang/semantic-release-gitlab, jstang/rasic"
			version:      "32.131.1"
			platform:     "gitlab"
			gitlabToken:  client.env.GITLAB_TOKEN
			githubToken:  client.env.GITHUB_TOKEN
		}
	}
}
