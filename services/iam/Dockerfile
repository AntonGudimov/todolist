# ---------- Build ----------
FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG BUILD_DIR=/usr/local/bin

RUN CGO_ENABLED=0 GOOS=linux \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o ${BUILD_DIR}/iam-server \
    ./cmd/server

RUN CGO_ENABLED=0 GOOS=linux \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o ${BUILD_DIR}/iam-migrater \
    ./cmd/migrate

# ---------- Runtime ----------
FROM alpine:3.22 AS server

# RUN apk add --no-cache ca-certificates

ARG BUILD_DIR=/usr/local/bin
WORKDIR ${BUILD_DIR}

COPY --from=builder ${BUILD_DIR}/iam-server ./server

RUN apk add --no-cache curl
RUN addgroup -S app && adduser -S app -G app

USER app

EXPOSE 8080

ENTRYPOINT ["./server"]


# ---------- Migration ----------
FROM alpine:3.22 AS migrate

# RUN apk add --no-cache ca-certificates

ARG BUILD_DIR=/usr/local/bin
WORKDIR ${BUILD_DIR}

COPY --from=builder ${BUILD_DIR}/iam-migrater ./migrate

RUN addgroup -S app && adduser -S app -G app

USER app

ENTRYPOINT ["./migrate"]