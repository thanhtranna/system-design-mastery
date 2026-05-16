FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/outbox-relay ./cmd/outbox-relay

FROM gcr.io/distroless/static:nonroot
COPY --from=build /bin/outbox-relay /outbox-relay
USER nonroot
ENTRYPOINT ["/outbox-relay"]
