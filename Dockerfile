FROM golang:1.26-alpine AS builder

WORKDIR /app

# copy go mod
COPY go.mod go.sum ./

# SSH mount
RUN --mount=type=ssh \
    go mod download

# copy source
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/api

# -------- RUN STAGE --------
FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/app .

USER nonroot:nonroot

ENTRYPOINT ["/app/app"]