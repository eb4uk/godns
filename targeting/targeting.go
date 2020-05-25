package targeting

import (
	"fmt"
	"github.com/eb4uk/godns/cache"
	"github.com/gomodule/redigo/redis"
	"strings"
	"sync"
	"time"
)

type CallerHostProvider interface {
	GetTargetedResponse(key string) (a []string, err error)
}

type RedisCallerHostProvider struct {
	c          redis.Conn
	inMemCache map[string]string
	mu         sync.RWMutex
}

func NewRedisTargetingProvider(client redis.Conn) *RedisCallerHostProvider {
	r := &RedisCallerHostProvider{}
	r.c = client

	go r.startRefreshing()
	return r
}

var targetRedisKey = "godns:target"

func (r *RedisCallerHostProvider) startRefreshing() {
	ticker := time.NewTicker(time.Second * 10)
	for {
		r.refresh()
		<-ticker.C
	}
}
func (r *RedisCallerHostProvider) refresh() {
	stringMap, err := redis.StringMap(r.c.Do("HGETALL", targetRedisKey))
	if err != nil {
		fmt.Println("refresh target responses failed", err)
		return
	}
	r.mu.Lock()
	r.inMemCache = stringMap
	r.mu.Unlock()
}
func (r *RedisCallerHostProvider) GetTargetedResponse(key string) (a []string, err error) {
	r.mu.RLock()
	s, ok := r.inMemCache[key]
	r.mu.RUnlock()
	if !ok {
		err = cache.KeyNotFound{}
		return
	}
	a = strings.Split(s, ",")
	return
}
