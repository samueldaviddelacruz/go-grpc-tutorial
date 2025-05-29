# gRPC Tutorial Project

This repository contains my personal code-along and notes from the [TECH SCHOOL YouTube gRPC tutorial](https://www.youtube.com/@TECHSCHOOLGURU) series. The goal of this project is to deepen my understanding of gRPC, Protocol Buffers, and how to build efficient, scalable microservices using Go.

## 📚 Tutorial Overview

The tutorial covers key concepts and best practices, including:

- Protocol Buffers (.proto files)
- Generating Go code from .proto definitions
- Implementing gRPC services in Go
- Unary and streaming RPCs (client-side, server-side, bidirectional)
- Context timeouts and cancellation
- Authentication and middleware
- TLS encryption and mutual TLS authentication
- Connecting clients to gRPC servers
- gRPC-Gateway: Adding a RESTful HTTP server that proxies to the gRPC backend

- OpenAPI v2: Auto-generating an OpenAPI specification from Protobufs

## 💡 What You'll Find in This Repo

- `proto/`: All `.proto` files defining the services and messages
- `pb/`: Auto-generated Go code from the `.proto` definitions
- `server/`: gRPC server implementation
- `client/`: Example client to test services

## 🧠 Why I'm Doing This

I'm following this tutorial to:

- Gain practical experience with gRPC and Protobuf
- Understand gRPC's role in modern backend systems
- Write idiomatic, maintainable Go code
- Prepare to build secure, scalable microservices in real-world projects

## 🚀 Getting Started

To run the project locally:

```bash
# Clone the repo
git clone https://github.com/samueldaviddelacruz/go-grpc-tutorial.git
cd grpc-tutorial

# Install dependencies
go mod tidy

# Generate code from .proto
make gen

# Run the server
go run server/main.go

# Run the client
go run client/main.go
