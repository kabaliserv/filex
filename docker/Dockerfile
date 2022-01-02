FROM golang:1.17.5-alpine AS builder
LABEL stage=builder
WORKDIR /go/src/github.com/kabaliserv/filex

# Add gcc and libc-dev early so it is cached
RUN set -xe \
	&& apk add --no-cache gcc libc-dev

# Install dependencies earlier so they are cached between builds
COPY ../go.mod go.sum ./
RUN set -xe \
	&& go mod download

# Copy the source code, because directories are special, there are separate layers
COPY ../cmd ./cmd/

RUN set -xe \
	&& GOOS=linux GOARCH=amd64 go build \
	    -tags prod \
        -o /go/bin/kbs-filex github.com/kabaliserv/filex/cmd/kbs-filex


FROM alpine:3.15.0

COPY --from=builder /go/bin/ks-share /usr/local/bin/kbs-share

EXPOSE 8080

ENTRYPOINT ["kbs-share"]