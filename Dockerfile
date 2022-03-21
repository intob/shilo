############################
# STEP 1 build the binary
############################

FROM golang:alpine as builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR $GOPATH/src/github.com/intob/shilo/
COPY . .

# Fetch dependencies
RUN go get -d -v

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' -a \
    -o /go/bin/shilo .

############################
# STEP 2 build a small image
############################

FROM alpine:latest

RUN apk update && apk add ffmpeg

COPY --from=builder /go/bin/shilo /etc/shilo

ENV PORT "9000"

EXPOSE $PORT

ENTRYPOINT ["/etc/shilo"]
