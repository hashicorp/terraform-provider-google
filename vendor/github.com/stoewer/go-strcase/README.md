[![Build Status](https://travis-ci.org/stoewer/go-strcase.svg?branch=master)](https://travis-ci.org/stoewer/go-strcase)
[![Coverage Status](https://coveralls.io/repos/github/stoewer/go-strcase/badge.svg?branch=master)](https://coveralls.io/github/stoewer/go-strcase?branch=master)
[![GoDoc](https://godoc.org/github.com/stoewer/go-strcase?status.svg)](https://godoc.org/github.com/stoewer/go-strcase)
---

Go strcase
==========

The package `strcase` converts between different kinds of naming formats such as camel case 
(`CamelCase`), snake case (`snake_case`) or kebab case (`kebab-case`).
The package is designed to work only with strings consisting of standard ASCII letters. 
Unicode is currently not supported.

Versioning and stability
------------------------

Although the master branch is supposed to remain always backward compatible, the repository
contains version tags in order to support vendoring tools.
The tag names follow semantic versioning conventions and have the following format `v1.0.0`.
This package supports Go modules introduced with version 1.11.

Example
-------

```go
import "github.com/stoewer/go-strcase"

var snake = strcase.SnakeCase("CamelCase")
```

Dependencies
------------

### Build dependencies

* none

### Test dependencies

* `github.com/stretchr/testify`

Run linters and unit tests
-------------------------- 

Since some of the linters ran by gometalinter don't support go modules yet, test dependencies have to be
loaded to the vendor directory first and gometalinter itself must run with disabled module support:

```
go mod vendor
GO111MODULE=off gometalinter --config=.gometalinter.json --deadline=10m .
```

To run the test use the following commands:

```
go test .
```