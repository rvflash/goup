# Go Up

[![GoDoc](https://godoc.org/github.com/rvflash/goup?status.svg)](https://godoc.org/github.com/rvflash/goup)
[![Build Status](https://img.shields.io/travis/rvflash/goup.svg)](https://travis-ci.org/rvflash/goup)
[![Code Coverage](https://img.shields.io/codecov/c/github/rvflash/goup.svg)](http://codecov.io/github/rvflash/goup?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/goup?)](https://goreportcard.com/report/github.com/rvflash/goup)


`goup` checks if there are any updates for imports in your module.
It parses `go.mod` files to get dependencies with their version, uses [go-git](https://github.com/src-d/go-git) 
to retrieve the list of remote tags and performs comparisons to advise to update if necessary.

> For now, it doesn't update the go mod file. An option to force it is planned but for the moment,
> the main purpose is using it as a linter in continuous integration or in development process.


## Features

1. No dependency. Pure Go tool designed to be used as a linter. Zero call to `go` or `git` command line tools.
1. As `go.mod` uses the semantic versioning for module version, `goup` does the same and provides 3 modes: major, 
major+minor and by default, path. 
1. Takes care of each part of a mod file: `require`, `exclude` and `replace`.
1. Allows the capacity to force some modules to only use release tag, no prerelease.
1. Manages one or more `go.mod` files, for example with `./...` as parameter. 
1. As with go1.14, you can use the `GOINSECURE` environment variable to skip certificate validation and do
not require an HTTPS connection. Since version `v0.3.0`, `GOPRIVATE` has the same behavior. 
1. Can amend on demand `go.mod` files with deprecated dependencies to update them.

## Demo

```shell script
$ goup -v -m ./...
$ goup: github.com/rvflash/goup: golang.org/x/mod v0.2.1-0.20200121190230-accd165b1659 is up to date
$ goup: github.com/rvflash/goup: github.com/matryer/is v1.1.0 must be updated with v1.2.0
$ goup: github.com/rvflash/goup: github.com/golang/mock v1.4.0 is up to date
$ goup: github.com/rvflash/goup: gopkg.in/src-d/go-git.v4 v4.13.1 is up to date
```

## Installation

It's important to have reproducible CI, so it's recommended to install a specific version of `goup` available
on the [releases page](https://github.com/rvflash/goup/releases).


### Go

```shell script
GO111MODULE=on go get github.com/rvflash/goup@v0.3.0
```

### Docker

```shell script1
docker run --rm -v $(pwd):/pkg rvflash/goup:v0.3.0 goup -V
```

## Usage

```shell script
goup [flags] [modfiles]
```

The `goup` command is used to check updates of any dependencies in the go.mod file.
It supports the following flags:

* `-M`: ensures to have the latest major version. By default: only the path is challenged.
* `-m`: ensures to have the latest couple major with minor version. By default: only the path is challenged.
* `-V`: prints the version of the tool.
* `-f`: force the update of the go.mod file as advised
* `-i`: allows to exclude indirect modules.
* `-r`: it's a comma-separated list of glob patterns to match the repository paths where to force tag usage.
For example with `github.com/group/*` as value, any modules in this repository group must have a release tag,
no prerelease. 
* `-s`: forces the process to exit on first error occurred.
* `-t`: defines the maximum time duration to perform the check. By default: 10s. 
* `-v`: verbose output

`[modfiles]` can be one or more direct path to `go.mod` files, `.` or `./...` to get all those in the tree.

Using example with an auto-signed local git repository:

```shell script
GOINSECURE="gitlab.example.lan/*/*" goup -v .
```
