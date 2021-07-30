# Build the manager binary
FROM golang:1.15 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# Copy DBaaSProvider config
COPY config/dbaasprovider/dbaas_provider.yaml dbaas_provider.yaml

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source & git info
COPY main.go main.go
COPY .git/ .git/
COPY pkg/ pkg/

ARG VERSION
ENV PRODUCT_VERSION=${VERSION}

# Build
RUN if [ -z $PRODUCT_VERSION ]; then PRODUCT_VERSION=$(git describe --tags); fi; \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on \
    go build -a -ldflags="-X main.version=$PRODUCT_VERSION" -o manager main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
COPY --from=builder /workspace/dbaas_provider.yaml .

USER nonroot:nonroot

ENTRYPOINT ["/manager"]
