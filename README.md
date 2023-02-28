[![Go](https://github.com/go-lean/bevaluate/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/go-lean/bevaluate/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/go-lean/bevaluate/branch/master/graph/badge.svg?token=iOmRU6RDLt)](https://codecov.io/gh/go-lean/bevaluate)

# Go-Lean BEvaluate

## Install
    go install github.com/go-lean/bevaluate@latest

## About
BEvaluate is a tool aiming to drastically optimise the testing and deployment for GO source code
hosted inside a monorepo project. Evaluating the need of packages to be retested
and redeployed after a change to the codebase can make the difference between
having a build that only takes a few minutes and waiting for your entire codebase to be
retested and redeployed.

## Convention
In a monorepo typically there is one root folder containing the go module file
and a number of sub-folders with different packages, some of which could also be
deployed in a certain step of your CI/CD process.

### Expected structure
Some aspects of the structure, like the name of the deployments folder, can be configured,
meaning you can skip having a deployments folder all together, but it is still expected
to have a root folder containing all packages and the module file.
```yaml
root
  - cmd // optional deployments dir
    - service
      Dockerfile
      main.go
  - service
    server.go
    server_test.go
  - response
    response.go
    response_test.go
go.mod
```

## Evaluation
The process of evaluating the packages consists of several steps. First all packages
are being recursively parsed and their dependencies stored are collected. Next the dependencies
of each package are used to build a dependants graph. This graph is the key to determining
whether a change in one package would affect the testing or deployment of dependant packages.

## Init
If you want to use the tool with the default settings you can skip this step.
However, it is quite useful to specify custom scenarios according to your needs, 
so once you have successfully installed to tool you can run the init command.

    bevaluate init
This will create the config file in the root directory under the name of `bevaluate.yaml`
using the default settings.
```yaml
packages:
    ignored_dirs: [build$, vendor$, .*/mocks$]
evaluations:
    deployments_dir: cmd/
    retest_out: bevaluate/retest.out
    redeploy_out: bevaluate/redeploy.out
    special_cases:
        retest_triggers: []
        full_scale_triggers: [go.mod$]

```