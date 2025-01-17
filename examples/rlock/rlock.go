package main

import (
	"log"
	"time"

	godisson "github.com/ggdream/go-redisson"
	"github.com/redis/go-redis/v9"
)

func main() {

	// create redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()

	g := godisson.NewGodisson(rdb, godisson.WithWatchDogTimeout(30*time.Second))

	// lock with watchdog without retry
	lock := g.NewRLock("godisson")

	err := lock.Lock()
	if err == godisson.ErrLockNotObtained {
		log.Println("Could not obtain lock")
	} else if err != nil {
		log.Fatalln(err)
	}
	defer lock.Unlock()

	// lock with retry、watchdog
	// leaseTime value is -1, enable watchdog
	lock2 := g.NewRLock("godission-try-watchdog")

	err = lock2.TryLock(20000, -1)
	if err == godisson.ErrLockNotObtained {
		log.Println("Could not obtain lock")
	} else if err != nil {
		log.Fatalln(err)
	}
	time.Sleep(10 * time.Second)
	defer lock.Unlock()
}
