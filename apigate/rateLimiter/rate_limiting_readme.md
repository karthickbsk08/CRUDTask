
# 🛑 Rate Limiting in Go – Concepts & Implementation Guide

This README covers the most common **rate limiting algorithms** and how to implement them in Go using the built-in `golang.org/x/time/rate` package.

---

## 📌 What is Rate Limiting?

**Rate limiting** is a technique to control how frequently clients can hit your API or access a resource.

Use cases include:

- Preventing abuse (DDoS attacks)
- Enforcing API usage policies (e.g., 100 requests/min)
- Reducing server load

---

## 🧠 Common Rate Limiting Algorithms

| Algorithm        | Description |
|------------------|-------------|
| **Fixed Window** | Counts requests in fixed time intervals (e.g., per minute). Simple but prone to bursts. |
| **Sliding Window** | Smooths out traffic by using a rolling time window. |
| **Leaky Bucket** | Think of a bucket leaking at a fixed rate. Ensures a steady request rate. |
| **Token Bucket** | Bucket fills with tokens at a fixed rate. Each request consumes a token. Allows bursts up to bucket size. |

---

## ⏱️ Go's `golang.org/x/time/rate` Package

This package implements the **Token Bucket algorithm** and is widely used in Go apps.

### 🔧 Installation

```bash
go get golang.org/x/time/rate
```

---

## ✅ Basic Example: Global Rate Limiter

```go
package main

import (
    "fmt"
    "net/http"
    "golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(2, 5) // 2 req/sec with burst of 5

func handler(w http.ResponseWriter, r *http.Request) {
    if !limiter.Allow() {
        http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
        return
    }
    fmt.Fprintln(w, "Request successful!")
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
```

---

## 🔁 Per-Client Rate Limiting (IP-based)

```go
type client struct {
    limiter *rate.Limiter
}

var clients = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getClientLimiter(ip string) *rate.Limiter {
    mu.Lock()
    defer mu.Unlock()

    if l, exists := clients[ip]; exists {
        return l
    }

    l := rate.NewLimiter(1, 3)
    clients[ip] = l
    return l
}
```

---

## 🚿 Simulating Leaky Bucket (Manually)

The `time/rate` package is optimized for token bucket. To simulate **leaky bucket**, you can:

- Queue requests in a channel
- Drain them at a fixed rate using a ticker
- Limit queue size to simulate overflow

---

## 📎 Summary

| Feature | time/rate support |
|--------|--------------------|
| Global Rate Limit | ✅ |
| Per-User/IP Limit | ✅ |
| Token Bucket | ✅ |
| Leaky Bucket | 🚧 (Manual Implementation) |

---

## 🔗 References

- [Go rate package docs](https://pkg.go.dev/golang.org/x/time/rate)
- [Rate Limiting Explained (Cloudflare)](https://developers.cloudflare.com/rate-limiting/about/)
- [Robust Rate Limiting with Go](https://www.alexedwards.net/blog/how-to-rate-limit-http-requests)

---

## 🧪 Tip

Use middlewares to apply rate limit logic cleanly in your API projects.

