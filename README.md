# go-bytes-pool

Simple `[]byte` buffer pool for go backend by `sync.Pool`.

[![Go Reference](https://pkg.go.dev/badge/github.com/IrineSistiana/go-bytes-pool.svg)](https://pkg.go.dev/github.com/IrineSistiana/go-bytes-pool)

```go
package main

import "github.com/IrineSistiana/go-bytes-pool"

func main() {
	p := NewPool(16) 
	b := p.Get(1024)
	p.Release(b)
}
```
