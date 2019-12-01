package main

import (
	"fmt"
	"time"
)

type TokenBucket struct {
	fillInterval time.Duration
	capacity int64
	Bucket chan struct{}
}

func (t *TokenBucket) fillToken()  {
	c := time.NewTicker(t.fillInterval)
	for {
		select {
		case <- c.C:
			select {
			case t.Bucket <- struct{}{}:
			default:
			}
			fmt.Println("token count %d in %v", len(t.Bucket), time.Now().UTC())
		}
	}
}


func main()  {
	done := make(chan struct{})
	tb := &TokenBucket{
		fillInterval: time.Millisecond * 1000,
		capacity: 100,
	}
	tb.Bucket = make(chan struct{}, tb.capacity)

	go tb.fillToken()
	<- done
}
