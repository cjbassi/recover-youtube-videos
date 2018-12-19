FROM golang:latest as builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build .

FROM alpine
RUN apk add -U --no-cache ca-certificates
WORKDIR /webapp
COPY --from=builder /build .
CMD ["./backend"]
