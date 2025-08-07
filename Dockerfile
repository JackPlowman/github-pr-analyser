#checkov:skip=CKV_DOCKER_2
#checkov:skip=CKV_DOCKER_3
FROM golang:1.24.6

ENV GO111MODULE="on"

WORKDIR /
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

ENTRYPOINT ["go", "run", "main.go"]
