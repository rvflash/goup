# Go Up

[![GoDoc](https://godoc.org/github.com/rvflash/goup?status.svg)](https://godoc.org/github.com/rvflash/goup)
[![Build Status](https://img.shields.io/travis/rvflash/goup.svg)](https://travis-ci.org/rvflash/goup)
[![Code Coverage](https://img.shields.io/codecov/c/github/rvflash/goup.svg)](http://codecov.io/github/rvflash/goup?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/goup)](https://goreportcard.com/report/github.com/rvflash/goup)


`goup` checks if there are any updates for imports in your module.
It parses `go.mod` files to get dependencies with their used version and uses [go-git](https://github.com/src-d/go-git) to list remote tags. 

> For now, it doesn't update the go mod file. An option to force it is planned but for the moment,
> the main purpose is using it as a linter.


## Features

1. No dependency. Pure Go tool designed to be used as a linter. Zero call to `go` or `git` command line tools.
1. As `go.mod` uses the semantic versioning for module version, `goup` does the same and provides 3 modes: major, major+minor and by default, path. 
1. Takes care of each part of a mod file: `require`, `exclude` and `replace`.
1. Allows the capacity to force some module to only use release tag, no prerelease.
1. Manages one or more `go.mod` files, for example with `./...` as parameter. 


## Installation

It's important to have reproducible CI, so it's recommended to install a specific version of `goup` available
on the [releases page](https://github.com/rvflash/goup/releases).


### Go

```shell script
GO111MODULE=on go get github.com/rvflash/goup@v0.1.0
```

### Docker

```shell script
docker run --rm -v $(pwd):/pkg rvflash/goup:v0.1.0 goup -V
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
* `-i`: allows to exclude indirect modules.
* `-r`: it's a comma separated list of repositories (or part of) used to force tag usage.
For example with `gitlab` as value, any modules with this word in the path must have a release tag, no prerelease.  
* `-s`: forces the process to exit on first error occurred.
* `-t`: defines the maximum time duration to perform the check. By default: 10s. 
* `-v`: verbose output

`[modfiles]` can be one or more direct path to `go.mod` files, `.` or `./...` to get all those in the tree.


## Soon

An option `-f` will be available to force updates of the go.mod file as recommended.