FROM golang:1.25.0

ENV GO111MODULE="on"

WORKDIR /
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

ENTRYPOINT ["go", "run", "main.go"]
