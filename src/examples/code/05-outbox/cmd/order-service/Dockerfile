FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/order-service ./cmd/order-service

FROM gcr.io/distroless/static:nonroot
COPY --from=build /bin/order-service /order-service
USER nonroot
EXPOSE 8080
ENTRYPOINT ["/order-service"]
