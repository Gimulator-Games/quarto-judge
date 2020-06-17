FROM golang:alpine as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o judge-bin cmd/quarto/main.go


FROM alpine

WORKDIR /app

COPY --from=builder /build/judge-bin judge
