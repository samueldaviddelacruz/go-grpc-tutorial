# gRPC Tutorial Project

This repository contains my personal code-along and notes from the [TECH SCHOOL YouTube gRPC tutorial](https://www.youtube.com/@TECHSCHOOLGURU) series. The goal of this project is to deepen my understanding of gRPC, Protocol Buffers, and how to build efficient, scalable microservices using Go.

## 📚 Tutorial Overview

The tutorial covers key concepts and best practices, including:

- Protocol Buffers (`.proto` files)
- Generating Go code from `.proto`
- Implementing gRPC services in Go
- Unary and streaming RPCs
- Authentication and middleware (later in the series)
- Connecting clients to gRPC servers
- Bonus: Docker, database integration, and production tips

## 💡 What You'll Find in This Repo

- `proto/`: All `.proto` files defining the services and messages
- `pb/`: Auto-generated Go code from the `.proto` definitions
- `server/`: gRPC server implementation
- `client/`: Example client to test services
- `docs/`: Notes and references from the tutorial

## 🧠 Why I'm Doing This

I'm following this tutorial to:

- Get hands-on experience with gRPC
- Learn how gRPC fits into modern backend systems
- Practice writing idiomatic Go code
- Prepare for building real-world microservices

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
