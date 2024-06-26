FROM golang:1.20.4-alpine3.18 AS builder

WORKDIR /src
COPY . .
RUN go mod tidy
RUN go build ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /src/main .
COPY .env .
CMD ["./main"]