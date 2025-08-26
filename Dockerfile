# trivy:ignore:AVD-DS-0002
FROM golang:1.24.3

ENV GO111MODULE="on"

WORKDIR /
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

ENTRYPOINT ["go", "run", "main.go"]
