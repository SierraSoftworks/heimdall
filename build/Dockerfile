FROM golang:1.10

ARG SENTRY_DSN=""
ARG VERSION="development"

ENV GO_PATH="/go"

ADD . $GO_PATH/src/github.com/SierraSoftworks/heimdall
WORKDIR $GO_PATH/src/github.com/SierraSoftworks/heimdall

RUN go get -t ./...
RUN go test -v ./...

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -o bin/heimdall -a -installsuffix cgo -ldflags "-s -X main.version=$VERSION -X main.sentryDSN=$SENTRY_DSN" github.com/SierraSoftworks/heimdall/cmd/heimdall

# Build the actual container
FROM alpine:latest
LABEL maintainer="Sierra Softworks <admin@sierrasoftworks.com>"

COPY --from=0 /go/src/github.com/SierraSoftworks/bender/bin/heimdall /bin/heimdall
ADD configs/* /etc/heimdall/

EXPOSE 8080

LABEL VERSION=$VERSION

WORKDIR /bin
ENTRYPOINT ["/bin/heimdall"]
CMD ["--config", "/etc/heimdall/"]