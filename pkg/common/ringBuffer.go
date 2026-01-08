package common

import (
	"errors"
	"sync"
)

var ErrNotEnoughData = errors.New("ringbuffer: not enough data")

type RingBuffer struct {
	data  []byte
	wHead int
	rHead int
	size  int
	cap   int
	mu    sync.Mutex
}

func NewRingBuffer(capacity int) *RingBuffer {
	if capacity <= 0 {
		capacity = 4096 
	}
	return &RingBuffer{
		data: make([]byte, capacity),
		cap:  capacity,
	}
}

func (r *RingBuffer) Write(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(p) > r.cap {
		p = p[len(p)-r.cap:]
	}

	pSize := len(p)
	totalWrite := pSize 

	for pSize > 0 {
		n1 := copy(r.data[r.wHead:], p)
		pSize -= n1
		r.size += n1
		r.wHead = (r.wHead + n1) % r.cap
		if r.size > r.cap {
			overwritten := r.size - r.cap
			r.rHead = (r.rHead + overwritten) % r.cap
			r.size = r.cap
		}
		p = p[n1:]
	}
    

	return totalWrite, nil
}


func (r *RingBuffer) Read(buf []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	targetSize := len(buf)

	if targetSize == 0 {
		return 0, nil
	}

	if r.size < targetSize {
		return 0, ErrNotEnoughData
	}

	p := buf

	for len(p) > 0 {
		n := copy(p, r.data[r.rHead:])
		r.rHead = (r.rHead + n) % r.cap
		p = p[n:]
	}
	r.size -= targetSize
	return targetSize, nil
}

func (r *RingBuffer) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.size
}


func (r *RingBuffer) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rHead = 0
	r.wHead = 0
	r.size = 0
}