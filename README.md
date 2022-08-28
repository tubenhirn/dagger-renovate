# dagger-renovate

do renovate dependency updates with https://dagger.io

## run renovate

export required env and run renovate job

``` shell
export RENOVATE_REPOSITORIES="a-namespace/a-project"
export ACCESS_TOKEN="a-access-token"
export GITHUB_TOKEN="a-github-token"
dagger do renovate-<github|gitlab>
```

