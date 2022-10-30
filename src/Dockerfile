ARG GO_VERSION=${GO_VERSION:-1.19}

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .
# Build the binary.
RUN CGO_ENABLED=0  GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -tags=nomsgpack -o /app .

RUN cat /etc/passwd | grep nobody > /etc/passwd.nobody

# build a small image
FROM alpine

COPY --from=builder /etc/passwd.nobody /etc/passwd
COPY --from=builder /app /app

# Run
USER nobody
ENTRYPOINT ["/app"]
