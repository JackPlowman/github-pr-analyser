# Just variables

export package_directory := "src"

# ------------------------------------------------------------------------------
# General
# ------------------------------------------------------------------------------

build:
    echo TODO

test:
    echo TODO

# ------------------------------------------------------------------------------
# Go
# ------------------------------------------------------------------------------

go-format:
    cd {{ package_directory }} && go fmt ./...

go-staticcheck:
    cd {{ package_directory }} && go staticcheck ./...

# ------------------------------------------------------------------------------
# Prettier
# ------------------------------------------------------------------------------

prettier-check:
    prettier . --check

prettier-format:
    prettier . --check --write

# ------------------------------------------------------------------------------
# Justfile
# ------------------------------------------------------------------------------

format:
    just --fmt --unstable

format-check:
    just --fmt --check --unstable
