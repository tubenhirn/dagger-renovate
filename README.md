# dagger-renovate

do renovate dependency updates with dagger

## run renovate

export required env and run renovate job

``` shell
export RENOVATE_REPOSITORIES="a-namespace/a-project"
export GITLAB_TOKEN="a-gitlab-token"
export GITHUB_TOKEN="a-github-token"
dagger do renovate
```

