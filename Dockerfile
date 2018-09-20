FROM golang:1.11

# Setup dir
WORKDIR /go/src/github.com/yagotome/gcp-consumer
ADD . .

# Setting environment variables
ARG env=development
ARG gcp_credentials
ENV GOBIN /go/bin
ENV GOENV $env
ENV GOOGLE_APPLICATION_CREDENTIALS $gcp_credentials

# Install dependencies
RUN go get -u github.com/golang/dep/cmd/dep
RUN rm -rf vendor/
RUN dep ensure

# Build
RUN go build -o binary cmd/worker/main.go
RUN mv binary ${GOBIN}/gcp-consumer

CMD ["gcp-consumer"]