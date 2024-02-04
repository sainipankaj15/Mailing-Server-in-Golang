# Golang API and gRPC Server

This project is a Golang-based API and gRPC server, offering two approaches for building and serving APIs: a traditional JSON server and a gRPC server. The project includes a proto file that needs to be compiled to generate the necessary code.

## Features

- **JSON API Server:** Traditional JSON-based APIs.
- **gRPC Server:** gRPC-based APIs for efficient communication.

## Getting Started

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/sainipankaj15/Mailing-Server-in-Golang
   ```

## Compile the Proto File:

Run the following command to compile the proto file and generate the required code:

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/mail.proto
```

# How to run 
- Clone the repo
- Run the command ```go mod tidy ```
- Go to server folder and Run ``` go run . ```
- Both server will start.