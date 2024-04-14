package util

import (
	"bytes"
	"testing"
)

func TestBufferPool(t *testing.T) {
	bufPool.New = func() interface{} {
		return make([]byte, 512)
	}
	bufPool1k.New = func() interface{} {
		return make([]byte, 1*1024)
	}
	bufPool2k.New = func() interface{} {
		return make([]byte, 2*1024)
	}
	bufPool5k.New = func() interface{} {
		return make([]byte, 5*1024)
	}
	bufPool16k.New = func() interface{} {
		return make([]byte, 16*1024)
	}

	tests := []struct {
		name       string
		size       int
		wantLength int
	}{
		{"Get 16k", 16 * 1024, 16 * 1024},
		{"Get 5k", 5 * 1024, 5 * 1024},
		{"Get 2k", 2 * 1024, 2 * 1024},
		{"Get 1k", 1 * 1024, 1 * 1024},
		{"Get Less than 1k", 512, 512},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := GetBuf(tt.size)
			if len(buf) != tt.wantLength {
				t.Errorf("Got buf length %v, want %v", len(buf), tt.wantLength)
			}
			capacityBefore := cap(buf)
			PutBuf(buf)
			buf2 := GetBuf(tt.size)
			if cap(buf2) != capacityBefore {
				t.Errorf("Buffer capacity changed after put/get, got %v, want %v", cap(buf2), capacityBefore)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	c1 := &MockReadWriteCloser{Buffer: bytes.NewBuffer(make([]byte, 1024))}
	c2 := &MockReadWriteCloser{Buffer: bytes.NewBuffer(make([]byte, 1024))}

	inCount, outCount := Join(c1, c2)
	if inCount == 0 || outCount == 0 {
		t.Fatalf("Transfer failed, got inCount = %d, outCount = %d", inCount, outCount)
	}

	t.Logf("c1 = %v, c2 = %v\n", inCount, outCount)
}

type MockReadWriteCloser struct {
	*bytes.Buffer
}

func (m *MockReadWriteCloser) Close() error {
	return nil
}
