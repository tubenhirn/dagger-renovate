# dagger-renovate

do renovate dependency updates with https://dagger.io

## run renovate

export required env and run renovate job.

can be run for github or gitlab repositories.

``` shell
export RENOVATE_GITHUB_REPOSITORIES="a-namespace/a-project"
export GITHUB_ACCESS_TOKEN="a-access-token"
make github
```

