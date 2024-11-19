FROM golang:1.23 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/silven-dynamics/go-ecommerce
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/bin/app ./catalog/cmd/catalog

FROM debian:bullseye-slim
WORKDIR /usr/bin
COPY --from=build /app/bin/app .
EXPOSE 8002
CMD ["./app"]