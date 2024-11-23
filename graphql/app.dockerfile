FROM golang:1.23 AS build
RUN apt-get update && apt-get install -y gcc g++ make ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /go/src/github.com/silven-dynamics/go-ecommerce
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY account account
COPY catalog catalog
COPY order order
RUN go build -o /app/bin/app ./graphql

FROM debian:bookworm-slim
WORKDIR /usr/bin
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=build /app/bin/app .
EXPOSE 8001
CMD ["./app"]

# FROM alpine:3.11
# WORKDIR /usr/bin
# RUN apk --no-cache add ca-certificates
# COPY --from=build /app/bin/app .
# EXPOSE 8001
# CMD ["./app"]