package bytesPool

import (
	"math/bits"
	"sync"
)

const (
	SmallBufSize = 1 << 16
	MaxBufSize   = 1<<32 - 1
)

var (
	sp         smallPool
	lp         largePool
)

// Get returns a *[]byte from pool with most appropriate cap.
// If size <= SmallBufSize, it returns a *[]byte with cap of
// 2^n. (At most 100% larger than the given size, or 50% waste)
// If MaxBufSize >= size > SmallBufSize, the cap of returned *[]byte is
// 2^n + 2^(n-2)*k (0<=k<=3) (At most 25% larger than the given size, or 20% waste)
// If size > MaxBufSize. Get will just call make([]byte, size).
func Get(size int) *[]byte {
	if size <= SmallBufSize {
		return sp.get(size)
	}
	return lp.get(size)
}

// Release releases b to the pool.
// If cap(b) > MaxBufSize, Release is noop.
// b should come from Get(). Release will panic if
// cap(b) is not suitable for the pool.
func Release(b *[]byte) {
	c := cap(*b)
	if c <= SmallBufSize {
		sp.release(b)
	} else {
		lp.release(b)
	}
}

type smallPool struct {
	ps [17]sync.Pool
}

func (p *smallPool) get(size int) *[]byte {
	if size < 0 {
		panic("negative buffer size")
	}

	var bit int
	if size > 1 { // corner case: size == 0 -> bit =0
		bit = bits.Len(uint(size - 1))
	}

	bp, ok := p.ps[bit].Get().(*[]byte)
	if !ok {
		b := make([]byte, size, 1<<bit)
		bp = &b
	}
	*bp = (*bp)[:size]
	return bp
}

func (p *smallPool) release(b *[]byte) {
	c := cap(*b)

	var bit int
	if c > 1 { // corner case: size == 0 -> bit =0
		bit = bits.Len(uint(c - 1))
	}

	if c != 1<<bit { // this buffer has a invalid cap size
		panic("release: invalid buf cap")
	}
	p.ps[bit].Put(b)
}

type largePool struct {
	ps [16][4]sync.Pool
}

func (p *largePool) get(size int) *[]byte {
	if size > MaxBufSize {
		b := make([]byte, size)
		return &b
	}

	ub := bits.Len(uint(size - 1))
	lb := (size - 1<<(ub-1) - 1) >> (ub - 3)
	bp, ok := p.ps[ub-17][lb].Get().(*[]byte)
	if !ok {
		b := make([]byte, size, 1<<(ub-1)+(lb+1)<<(ub-3))
		bp = &b
	}
	*bp = (*bp)[:size]
	return bp
}
func (p *largePool) release(b *[]byte) {
	c := cap(*b)
	if c > MaxBufSize {
		return
	}
	ub := bits.Len(uint(c - 1))
	lb := (c - 1<<(ub-1) - 1) >> (ub - 3)
	if c != 1<<(ub-1)+(lb+1)<<(ub-3) {
		panic("release: invalid buf cap")
	}
	p.ps[ub-17][lb].Put(b)
}
