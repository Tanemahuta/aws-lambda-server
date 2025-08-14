# Build the manager binary
FROM golang:1.25 as builder

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY buildinfo/ buildinfo/
COPY pkg/ pkg/

ARG VERSION
ARG COMMIT_SHA

# Build
RUN VERSION=${VERSION} COMMIT_SHA=${COMMIT_SHA} go generate ./... && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a  -o aws-lambda-server main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/aws-lambda-server .
USER 65532:65532

ENTRYPOINT ["/aws-lambda-server"]
