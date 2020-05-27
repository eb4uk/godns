package cache

import (
	"github.com/eb4uk/godns/models"
	"github.com/hoisie/redis"
	"github.com/miekg/dns"
)

type RedisCache struct {
	Backend *redis.Client
	Expire  int64
}

var cachePath = "godns:cache"

func (r *RedisCache) Get(key string) (*dns.Msg, error) {
	var msg dns.Msg
	item, err := r.Backend.Get(r.buildKeyPath(key))
	if err != nil {
		err = KeyNotFound{r.buildKeyPath(key)}
		return &msg, err
	}
	err = msg.Unpack(item)
	if err != nil {
		err = SerializerError{err}
	}
	return &msg, err
}

func (r *RedisCache) Set(key string, msg *dns.Msg) error {
	var val []byte
	var err error

	// handle cases for negacache where it sets nil values
	if msg == nil {
		val = []byte("nil")
	} else {
		val, err = msg.Pack()
	}
	if err != nil {
		err = SerializerError{err}
	}
	return r.Backend.Setex(r.buildKeyPath(key), r.Expire, val)
}

func (r *RedisCache) Exists(key string) bool {
	_, err := r.Backend.Get(r.buildKeyPath(key))
	if err != nil {
		return true
	}
	return false
}

func (r *RedisCache) Remove(key string) error {
	_, err := r.Backend.Del(r.buildKeyPath(key))
	return err
}

func (r *RedisCache) Full() bool {
	return false
}
func (r *RedisCache) buildKeyPath(key string) string {
	return cachePath + ":" + key
}

func NewRedisCache(rs models.RedisSettings, expire int64) *RedisCache {
	rc := &redis.Client{Addr: rs.Addr(), Db: rs.DB, Password: rs.Password}
	return &RedisCache{
		Backend: rc,
		Expire:  expire,
	}
}
