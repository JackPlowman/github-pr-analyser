# Development Guide

## Table of Contents

- [Development Guide](#development-guide)
  - [Table of Contents](#table-of-contents)
  - [Development](#development)
    - [Prerequisites](#prerequisites)
    - [First time setup](#first-time-setup)

## Development

### Prerequisites

- Go - Version can be found in the [go.mod](../github-pr-analyser/go.mod) file
  - Install Go from [here](https://go.dev/doc/install)
- Just
  - Install Just from [here](https://github.com/casey/just#installation)

### First time setup

1. Clone the repository and cd into the repository
2. Run `cd github-pr-analyser`
3. Run `go mod download` to download the dependencies
4. Run `go mod tidy` to ensure that the dependencies are up-to-date