FROM golang:latest AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o server .

FROM scratch

COPY --from=builder ["/build/server", "/"]

ENTRYPOINT ["/server"]
