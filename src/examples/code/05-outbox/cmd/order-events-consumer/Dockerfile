FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/order-events-consumer ./cmd/order-events-consumer

FROM gcr.io/distroless/static:nonroot
COPY --from=build /bin/order-events-consumer /order-events-consumer
USER nonroot
ENTRYPOINT ["/order-events-consumer"]
