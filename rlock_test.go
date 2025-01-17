package godisson

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"net"
	"testing"
	"time"
)

func TestNewGedisson(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	gedisson := NewGodisson(rdb)
	lock := gedisson.NewRLock("hkn")
	t.Log(lock.TryLock(20000, 40000))
	time.Sleep(10 * time.Second)
	lock.Unlock()

	//time.Sleep(100 * time.Second)
}

func TestNewGedisson1(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	gedisson := NewGodisson(rdb)
	lock1 := gedisson.NewRLock("hkn")

	t.Log(lock1.TryLock(20000, 40000))
	time.Sleep(10 * time.Second)
	lock1.Unlock()

}

func TestNewGedisson2(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	gedisson := NewGodisson(rdb)
	lock1 := gedisson.NewRLock("hkn")

	t.Log(lock1.TryLock(20000, 40000))
	time.Sleep(10 * time.Second)
	lock1.Unlock()

}
func TestNewGedissonWatchdog(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	gedisson := NewGodisson(rdb)

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				for _, i := range gedisson.RenewMap.Items() {
					e := i.(*RenewEntry)
					fmt.Printf("goids %v", e.goroutineIds)
				}
			}

		}
	}()

	lock := gedisson.NewRLock("hkn")
	t.Log(lock.TryLock(40000, -1))

	lock1 := gedisson.NewRLock("hkn")
	t.Log(lock1.TryLock(40000, -1))
	time.Sleep(1 * time.Minute)
	lock1.Unlock()

	time.Sleep(30 * time.Second)
	lock.Unlock()

}

func TestSub(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	sub := rdb.Subscribe(context.Background(), "123-channel")
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()

	message, err := sub.ReceiveMessage(ctx)
	var target *net.OpError
	t.Log(message, err, errors.As(err, &target))
}

func TestPub(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	rdb.Publish(context.Background(), "123-channel", "123")

}

func TestReLock(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	g := NewGodisson(rdb)

	lock1 := g.NewRLock("hkn")
	t.Log(lock1.TryLock(-1, 30000))
	time.Sleep(10 * time.Second)

	lock2 := g.NewRLock("hkn")
	t.Log(lock2.TryLock(-1, -1))

}
