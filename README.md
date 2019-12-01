# ratelimiter

--- 

# Ratelimit 服务流量限制

## 限流方法

1. 漏桶

    漏桶是指我们有一个一直装满了水的桶，每过固定的一段时间即向外漏一滴水。如果你接到了这滴水，那么你就可以继续服务请求，如果没有接到，那么就需要等待下一滴水。

2. 令牌桶

    令牌桶则是指匀速向桶中添加令牌，服务请求时需要从桶中获取令牌，令牌的数目可以按照需要消耗的资源进行相应的调整。如果没有令牌，可以选择等待，或者放弃。

漏桶流出速率固定，而令牌痛可以储存一定量的令牌，可以承受一定量的并发， 但消耗过多的情况下，令牌桶也会变成漏桶模型。

## 令牌桶

令牌桶看上去像一个全局加减的计数器，可以配合读写锁进行计数，在GO中，也可以用一个
channel 来完成加减的token的操作。

每隔fillInterval时间就往令牌桶bucket 添加token，超过capacity则放弃

```go
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
			fmt.Printf("token count %d in %v\n", len(t.Bucket), time.Now().UTC())
		}
	}
}

func main()  {
	done := make(chan struct{})
	tb := &TokenBucket{
		fillInterval: time.Millisecond * 10,
		capacity: 100,
	}
	tb.Bucket = make(chan struct{}, tb.capacity)

	go tb.fillToken()
	<- done
}

```

可以看到输出
```bash
token count 95 in 2019-12-01 06:24:51.036007 +0000 UTC
token count 96 in 2019-12-01 06:24:51.045262 +0000 UTC
token count 97 in 2019-12-01 06:24:51.056807 +0000 UTC
token count 98 in 2019-12-01 06:24:51.065552 +0000 UTC
token count 99 in 2019-12-01 06:24:51.075493 +0000 UTC
token count 100 in 2019-12-01 06:24:51.086583 +0000 UTC
token count 100 in 2019-12-01 06:24:51.097242 +0000 UTC
token count 100 in 2019-12-01 06:24:51.106018 +0000 UTC
```

