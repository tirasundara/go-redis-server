# go-redis-server 🚀  
*A Redis-compatible server built from scratch in Go*  

This project is a Redis server implementation in Go, built step-by-step following the [Codecrafters Redis Challenge](https://app.codecrafters.io/courses/redis/overview). It supports basic Redis commands, RDB persistence, and replication.

## Features Implemented ✅  

### **Core Redis Commands**  
✔ Bind to a port and accept connections  
✔ Handle multiple concurrent clients  
✔ Implement the following Redis commands:  
  - `PING`
  - `ECHO`
  - `SET` & `GET`  

### **RDB Persistence**  
✔ Parse and load RDB files  
✔ Support for:  
  - Reading a key and its value  
  - Handling multiple keys and values  
  - Expiry handling for keys  

### **Replication**  
✔ Configure the server to act as a replica  
✔ Implemented the `INFO` command for replication status  
✔ Initial replication setup:  
  - Send handshake (3/3)  
  - Receive handshake (1/2) _(in progress)_  

## Upcoming Roadmap 🛠️  

🔹 **Finish Replication:**  
  - Complete Receive handshake (2/2)  
  - Full master-replica synchronization  

🔹 **Support Redis Streams**  

🔹 **Implement Redis Transactions (`MULTI`, `EXEC`, `DISCARD`)**  

## Running the Server

### **Prerequisites**  
- Go 1.20+ installed  

### **Start the Redis server**  
```sh
go run app/main.go
```

### **Interact Using `redis-cli`**
```sh
redis-cli -p 6379
```

```sh
127.0.0.1:6379> PING
PONG

127.0.0.1:6379> SET key1 "hello"
OK

127.0.0.1:6379> GET key1
"hello"
```

## Contributing 🤝
Contributions are welcome! Feel free to open issues and PRs.
