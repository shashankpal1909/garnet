<p align="center">
  <h1 align="center">♦️ Garnet</h1>
  <p align="center">
    <strong>A high-performance, purely asynchronous, Redis-compatible in-memory data store written in Go.</strong>
  </p>
  <p align="center">
    <a href="https://golang.org/doc/go1.21"><img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go" alt="Go Version" /></a>
    <a href="https://github.com/shashankpal1909/garnet/actions"><img src="https://img.shields.io/badge/Build-Passing-brightgreen?style=for-the-badge&logo=githubactions" alt="Build Status" /></a>
    <a href="https://hub.docker.com/r/library/garnet"><img src="https://img.shields.io/badge/Docker-Supported-2496ED?style=for-the-badge&logo=docker" alt="Docker Support" /></a>
    <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge" alt="License: MIT" /></a>
    <a href="https://github.com/shashankpal1909/garnet/pulls"><img src="https://img.shields.io/badge/PRs-Welcome-ff69b4.svg?style=for-the-badge" alt="PRs Welcome" /></a>
  </p>
</p>

## 🚀 Overview

> [!WARNING]
> Garnet is a hobby and educational project created for learning purposes. It is **not** intended or expected to be used in production environments.

**Garnet** is a fast, lightweight, and asynchronous in-memory key-value store. Built from the ground up to speak the **RESP (REdis Serialization Protocol)**, Garnet acts as a drop-in replacement for basic Redis caching workloads while taking full advantage of Linux's `epoll` system calls to deliver non-blocking, high-throughput performance.

By sidestepping the overhead of traditional thread-per-connection models and utilizing an event-driven loop, Garnet easily serves tens of thousands of requests per second on minimal hardware.

## ✨ Features

- **Purely Asynchronous:** Powered by a custom non-blocking I/O event loop utilizing Linux `epoll`.
- **Redis Compatible:** Implements RESP, meaning your existing Redis clients (e.g., `redis-cli`, Jedis, Go-Redis) work out of the box!
- **High Performance:** Designed to deliver high throughput and low latency on minimal hardware.
- **Key Expiration Engine:** Native active & passive TTL enforcement for cache invalidation.
- **Lightweight & Containerized:** Provided via an Alpine-based Docker image. No heavy dependencies.

## 🛠️ Installation & Quick Start

Because Garnet leverages the native Linux `epoll` API for its event loop, the server **must** run on a Linux environment. For macOS and Windows users, **Docker** is the recommended way to run Garnet.

### Using Docker Compose (Recommended)

1. Clone the repository:
   ```bash
   git clone https://github.com/shashankpal1909/garnet.git
   cd garnet
   ```

2. Start the server using Docker Compose:
   ```bash
   docker compose up -d
   ```

3. Connect using standard `redis-cli`:
   ```bash
   redis-cli -h 127.0.0.1 -p 6379
   127.0.0.1:6379> PING
   PONG
   127.0.0.1:6379> SET hello "world"
   OK
   127.0.0.1:6379> GET hello
   "world"
   ```

### Building from Source (Linux Only)

If you are on a Linux machine, you can build and run Garnet natively:

```bash
go build -o garnet ./cmd/server/main.go
./garnet --host=0.0.0.0 --port=6379 --max-clients=10000
```

## 📚 Supported Commands

Garnet currently supports the core set of Redis commands needed for caching and basic key-value operations:

* **Connection:** `PING`, `ECHO`
* **Strings:** `SET`, `GET`
* **Keys:** `DEL`, `EXPIRE`, `TTL`

*Note: More commands are constantly being added to the roadmap.*

## 🤝 Contributing

Contributions make the open-source community an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.
