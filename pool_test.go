package bytesPool

import (
	"testing"
)

func TestGet(t *testing.T) {
	bl := 16 // 65535
	p := NewPool(bl)

	tests := []struct {
		name    string
		size    int
		wantCap int
	}{
		{"", 0, 0},
		{"", 1, 1},
		{"", 2, 3},
		{"", 3, 3},
		{"", 4, 7},
		{"", 127, 127},
		{"", 128, 255},
		{"", 65530, 65535},
		{"", 65535, 65535},
		{"", 65535 + 1, 65535 + 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bp := p.Get(tt.size)
			b := *bp
			if len(b) != tt.size {
				t.Fatalf("want size %d, got %v", tt.size, len(b))
			}
			if cap(b) != tt.wantCap {
				t.Fatalf("want cap %d, got %v", tt.wantCap, cap(b))
			}
			p.Release(bp)
		})
	}
}
