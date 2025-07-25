ARG GOLANG_IMAGE=golang:1.24-bookworm
ARG USE_BAZEL_VERSION=6.1.1

# STAGE 1: build...
FROM ${GOLANG_IMAGE} AS build

WORKDIR /app
COPY . ./

RUN go mod download

ARG GOPS_VERSION=v0.3.27
RUN CGO_ENABLED=0 go install -ldflags "-s -w" github.com/google/gops@${GOPS_VERSION}
RUN CGO_ENABLED=0 GOOS=linux go build -C /app/cmd/adapter -o /sgnl/adapter
RUN CGO_ENABLED=0 GOOS=linux go build -C /app/cmd/ldap-adapter -o /sgnl/ldap-adapter

# STAGE 2: run...
FROM gcr.io/distroless/static AS run

# Fixture files are loaded from `pkg/mock/...`, but we copy the files to `/sgnl/pkg/mock/...`
# so we change the working directory to `/sgnl` to make sure the files are found.
WORKDIR /sgnl

COPY --from=build --chown=nonroot:nonroot /go/bin/gops /sgnl/gops
COPY --from=build --chown=nonroot:nonroot /sgnl/adapter /sgnl/adapter
COPY --from=build --chown=nonroot:nonroot /sgnl/ldap-adapter /sgnl/ldap-adapter
COPY --from=build --chown=nonroot:nonroot /app/pkg/mock/servicenow/fixtures/*.yaml /sgnl/pkg/mock/servicenow/fixtures/

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT [ "/sgnl/adapter" ]
