package bytesPool

import (
	"math/bits"
	"sync"
)

type Pool struct {
	bitLen int
	ps     []sync.Pool
}

// NewPool Creates a pool to reuse []byte buffers if 
// the length of those buffers are smaller than (1<<bitLen)-1 .
func NewPool(bitLen int) *Pool {
	if bitLen < 0 {
		panic("bytesPool: negative bit length")
	}
	return &Pool{
		bitLen: bitLen,
		ps:     make([]sync.Pool, bitLen+1),
	}
}

// Get returns a *[]byte from pool with most appropriate cap.
// It returns a *[]byte with cap of (2^n)-1.
// If size is too big, Get simply calls make([]byte, size).
func (p *Pool) Get(size int) *[]byte {
	if size < 0 {
		panic("bytesPool: negative buffer size")
	}
	bit := bits.Len(uint(size))
	if bit > p.bitLen {
		b := make([]byte, size)
		return &b
	}

	bp, ok := p.ps[bit].Get().(*[]byte)
	if !ok {
		b := make([]byte, size, (1<<bit)-1)
		bp = &b
	}
	*bp = (*bp)[:size]
	return bp
}

// Release releases b to the pool.
// If cap(b) is too big, Release is noop.
// b should come from Get(). Release will panic if
// cap(b) is not suitable for the pool.
func (p *Pool) Release(b *[]byte) {
	c := cap(*b)
	bit := bits.Len(uint(c))
	if bit > p.bitLen {
		return
	}
	if c != (1<<bit)-1 { // this buffer has a invalid cap size
		panic("bytesPool: invalid buf")
	}
	p.ps[bit].Put(b)
}
