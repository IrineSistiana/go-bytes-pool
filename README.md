# go-bytes-pool

Simple `[]byte` buffer pool for go backend by `sync.Pool`. At most 50% memory waste for small buffers (`len(b) <= 65536`) and 20% for large buffers.

[![Go Reference](https://pkg.go.dev/badge/github.com/IrineSistiana/go-bytes-pool.svg)](https://pkg.go.dev/github.com/IrineSistiana/go-bytes-pool)

```go
package main

import "github.com/IrineSistiana/go-bytes-pool"

func main() {
	b := bytesPool.Get(1024)
	bytesPool.Release(b)
}
```
