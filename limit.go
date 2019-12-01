package main

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	fillInterval time.Duration
	capacity     int64
	Bucket       chan struct{}
	mu           sync.Mutex
}

func (t *TokenBucket) fillToken() {
	c := time.NewTicker(t.fillInterval)
	for {
		select {
		case <-c.C:
			select {
			case t.Bucket <- struct{}{}:
			default:
			}
			fmt.Printf("token count %d in %v\n", len(t.Bucket), time.Now().UTC())
		}
	}
}

func (t *TokenBucket) Take() {}

// take function in internal
func (t *TokenBucket) take() bool {
	select {
	case <-t.Bucket:
		return true
	default:
		return false
	}
}

func main() {
	done := make(chan struct{})
	tb := &TokenBucket{
		fillInterval: time.Millisecond * 10,
		capacity:     100,
	}
	tb.Bucket = make(chan struct{}, tb.capacity)

	go tb.fillToken()
	<-done
}
