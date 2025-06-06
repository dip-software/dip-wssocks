FROM golang:1.24.4 AS builder
WORKDIR /build
COPY go.mod .
COPY go.sum .
# We vendor dependencies to speed up the build
# RUN go mod download

# Build
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
  CGO_ENABLED=0 go build -o app -ldflags="-s -w" -trimpath

FROM alpine:latest 
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=builder /build/app /app
EXPOSE 1088
CMD ["/app/app"]
