# ------------------------------------------------------------------------------
# General
# ------------------------------------------------------------------------------

export SRC_DIR := "./src"
export SRC_RECURSIVE := "./src/..."

build:
    go build -o github-pr-analyser ${SRC_DIR}

run:
    go run ${SRC_DIR}

alias fmt := lint-fix
alias fmt-check := lint

lint:
    golangci-lint run ${SRC_RECURSIVE}

lint-fix:
    golangci-lint run --fix ${SRC_RECURSIVE}

vulncheck:
    govulncheck ${SRC_RECURSIVE}

test:
    go test -coverprofile=coverage.out ${SRC_RECURSIVE}

# ------------------------------------------------------------------------------
# Docker
# ------------------------------------------------------------------------------

docker-build:
    docker buildx build -t github-pr-analyser:latest .

# ------------------------------------------------------------------------------
# Prettier
# ------------------------------------------------------------------------------

# Check all files with prettier
prettier-check:
    prettier . --check

# Format all files with prettier
prettier-format:
    prettier . --check --write

# ------------------------------------------------------------------------------
# Justfile
# ------------------------------------------------------------------------------

# Format Justfile
format:
    just --fmt --unstable

# Check Justfile formatting
format-check:
    just --fmt --check --unstable

# ------------------------------------------------------------------------------
# Gitleaks
# ------------------------------------------------------------------------------

# Run gitleaks detection
gitleaks-detect:
    gitleaks detect --source .

# ------------------------------------------------------------------------------
# Lefthook
# ------------------------------------------------------------------------------

# Validate lefthook config
lefthook-validate:
    lefthook validate

# ------------------------------------------------------------------------------
# Zizmor
# ------------------------------------------------------------------------------

# Run zizmor checking
zizmor-check:
    uvx zizmor . --persona=auditor

# ------------------------------------------------------------------------------
# Pinact
# ------------------------------------------------------------------------------

# Run pinact
pinact-run:
    pinact run

# Run pinact checking
pinact-check:
    pinact run --verify --check

# Run pinact update
pinact-update:
    pinact run --update

# ------------------------------------------------------------------------------
# Git Hooks
# ------------------------------------------------------------------------------

# Install pre commit hook to run on all commits
install-git-hooks:
    lefthook install -f
