# syntax=docker/dockerfile:1

FROM golang:1.25.0 AS builder
WORKDIR /app
COPY . .
RUN go mod download
# Build
RUN CGO_ENABLED=1 GOOS=linux go build -o main .


FROM scratch
COPY --from=builder /app/main /main
# Run
CMD ["/main"]