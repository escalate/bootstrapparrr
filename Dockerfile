ARG GOLANG_VERSION
FROM golang:${GOLANG_VERSION}-bullseye AS builder

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
WORKDIR /go/src

ARG GOOS
ARG GOARCH
ENV GOOS=${GOOS} \
    GOARCH=${GOARCH}

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN make build


FROM debian:bullseye-slim

COPY --from=builder /go/src/dist/bootstrapparrr /bootstrapparrr
ARG GIN_MODE
ENV GIN_MODE=${GIN_MODE}
EXPOSE 5000
ENTRYPOINT ["/bootstrapparrr"]
