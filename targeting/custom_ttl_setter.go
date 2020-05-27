package targeting

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"sync"
	"time"
)

type CustomTtlSetter interface {
	GetTtl(domain string, currentTtl uint32) uint32
}

var ttlKeyPath = "godns:custom_ttl"

type RedisCustomTtlSetter struct {
	c          redis.Conn
	inMemCache map[string]int
	mu         sync.RWMutex
}

func NewRedisCustomTtlSetter(c redis.Conn) *RedisCustomTtlSetter {
	r := &RedisCustomTtlSetter{c: c, mu: sync.RWMutex{},
		inMemCache: map[string]int{}}
	go r.startRefresher()
	return r
}

func (s *RedisCustomTtlSetter) GetTtl(domain string, currentTtl uint32) uint32 {
	s.mu.RLock()
	ttl, ok := s.inMemCache[domain]
	if ok {
		fmt.Println("setting custom ttl", domain, ttl)
		currentTtl = uint32(ttl)
	}
	s.mu.RUnlock()

	return currentTtl
}

func (s *RedisCustomTtlSetter) refresh() {
	intMap, err := redis.IntMap(s.c.Do("HGETALL", ttlKeyPath))
	if err != nil {
		fmt.Println("error on updating custom ttl settings. error -", err)
		return
	}
	s.mu.Lock()
	s.inMemCache = intMap
	s.mu.Unlock()
	return
}

func (s *RedisCustomTtlSetter) startRefresher() {
	ticker := time.NewTicker(time.Second * 10)
	for {
		s.refresh()
		<-ticker.C
	}
}
