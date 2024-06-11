# Strict mode
set shell := ["bash", "-euo", "pipefail", "-c"]

# Globals
export GIT_SHA := 'git rev-parse --short HEAD'

build:
    go build -o bin/repo-checker cmd/checker/main.go