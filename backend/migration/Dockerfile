FROM golang:1.22.5-alpine3.20 AS builder

ARG MODE
ARG GIT_COMMIT_ID

ENV APP_MODE=${MODE:-"production"}
ENV APP_REVISION=${GIT_COMMIT_ID:-"unknown"}

ENV HOME=/usr/src/app
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV COMMON_MODULE_PATH=./common
ENV MIGRATION_MODULE_PATH=./migration

WORKDIR ${HOME}
COPY ${COMMON_MODULE_PATH}/go.mod ${COMMON_MODULE_PATH}/go.sum ${COMMON_MODULE_PATH}/
COPY ${MIGRATION_MODULE_PATH}/go.mod ${MIGRATION_MODULE_PATH}/go.sum ${MIGRATION_MODULE_PATH}/

WORKDIR ${COMMON_MODULE_PATH}
RUN go mod download && go mod verify

WORKDIR ${MIGRATION_MODULE_PATH}
RUN go mod download && go mod verify

WORKDIR ${HOME}
COPY ${COMMON_MODULE_PATH} ${COMMON_MODULE_PATH}/
COPY ${MIGRATION_MODULE_PATH} ${MIGRATION_MODULE_PATH}/

WORKDIR ${MIGRATION_MODULE_PATH}
RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w \
    -X 'x-operation/common/core/config.AppVersion=$(cat VERSION)' \
    -X 'x-operation/common/core/config.CompatibleVersion=$(cat COMPATIBLE_VERSION)' \
    -X 'x-operation/common/core/config.AppMode=${APP_MODE}' \
    -X 'x-operation/common/core/config.AppRevision=${APP_REVISION}' \
    " \
    -o "/usr/bin/migration" .

FROM alpine:3.20.1

COPY --from=builder "/usr/bin/migration" /usr/bin/migration

ENTRYPOINT ["/usr/bin/migration"]
