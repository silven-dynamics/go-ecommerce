FROM golang:1.23 AS build
RUN apt-get update && apt-get install -y gcc g++ make ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /go/src/github.com/silven-dynamics/go-ecommerce
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/bin/app ./order/cmd/order

FROM debian:bookworm-slim
WORKDIR /usr/bin
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=build /app/bin/app .
EXPOSE 8001
CMD ["./app"]