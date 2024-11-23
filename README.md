# gRPC Microservices with GraphQL API Gateway

## Overview

This project demonstrates a modern microservices architecture using gRPC for internal communication and GraphQL as the API gateway. It features services for account management, product catalog, and order processing, designed for scalability and flexibility.

**Note:** While this project was initially developed using the latest Go version, it uses an older Go version in the final code for enhanced stability and easier maintenance. This approach is particularly useful for projects with multiple interconnected components like gRPC, GraphQL, PostgreSQL, Elasticsearch, and Docker Compose.

## Features

- **Microservices for modularity and scalability:**
  - Account Service: Manages user accounts and data
  - Catalog Service: Handles product catalog powered by Elasticsearch
  - Order Service: Processes orders and manages transactions
- GraphQL Gateway for seamless API integration across services
- PostgreSQL and Elasticsearch for robust data handling
- Fully containerized with Docker Compose for easy deployment

## Project Structure

### Services

- **Account:** Manages user accounts
- **Catalog:** Manages products (powered by Elasticsearch)
- **Order:** Manages customer orders
- **GraphQL:** API gateway for client interaction

### Databases

- PostgreSQL for Account and Order services
- Elasticsearch for the Catalog service

## Getting Started

### Prerequisites

- Docker & Docker Compose installed
- Go programming environment set up

### Setup and Run

1. Clone the repository:

   ```bash
   git clone <repository-url>
   cd <project-directory>
   ```

2. Start services:

   ```bash
   docker-compose up -d --build
   ```

3. Access GraphQL Playground:
   Open your browser and go to: `http://localhost:8001/playground`

## gRPC Protobuf Setup

### Install protoc

```bash
wget https://github.com/protocolbuffers/protobuf/releases/download/v23.0/protoc-23.0-linux-x86_64.zip
unzip protoc-23.0-linux-x86_64.zip -d protoc
sudo mv protoc/bin/protoc /usr/local/bin/
```

### Install Go plugins for protoc

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Add Go plugins to system path

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
source ~/.bashrc
```

## GraphQL API Usage

### Sample Queries and Mutations

#### Query Accounts

```graphql
query {
  accounts {
    id
    name
  }
}
```

#### Create an Account

```graphql
mutation {
  createAccount(account: {name: "New Account"}) {
    id
    name
  }
}
```

#### Query Products

```graphql
query {
  products {
    id
    name
    price
  }
}
```

#### Create a Product

```graphql
mutation {
  createProduct(product: {name: "New Product", description: "A new product", price: 19.99}) {
    id
    name
    price
  }
}
```

#### Create an Order

```graphql
mutation {
  createOrder(order: {accountId: "account_id", products: [{id: "product_id", quantity: 2}]}) {
    id
    totalPrice
    products {
      name
      quantity
    }
  }
}
```

#### Query Account with Orders

```graphql
query {
  accounts(id: "account_id") {
    name
    orders {
      id
      createdAt
      totalPrice
      products {
        name
        quantity
        price
      }
    }
  }
}
```

### Advanced GraphQL Queries

#### Pagination and Filtering

```graphql
query {
  products(pagination: {skip: 0, take: 5}, query: "search_term") {
    id
    name
    description
    price
  }
}
```

#### Calculate Total Spent by an Account

```graphql
query {
  accounts(id: "account_id") {
    name
    orders {
      totalPrice
    }
  }
}
```

## Technologies Used

- **Go:** Backend language for microservices
- **gRPC:** Inter-service communication
- **GraphQL:** API gateway
- **PostgreSQL:** Database for Account and Order services
- **Elasticsearch:** Search and analytics engine for the Catalog service
- **Docker Compose:** Orchestration of services
