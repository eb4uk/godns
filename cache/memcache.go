package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/miekg/dns"
)

func NewMemcachedCache(servers []string, expire int32) *MemcachedCache {
	c := memcache.New(servers...)
	return &MemcachedCache{
		backend: c,
		expire:  expire,
	}
}

type MemcachedCache struct {
	backend *memcache.Client
	expire  int32
}

func (m *MemcachedCache) Set(key string, msg *dns.Msg) error {
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
	return m.backend.Set(&memcache.Item{Key: key, Value: val, Expiration: m.expire})
}

func (m *MemcachedCache) Get(key string) (*dns.Msg, error) {
	var msg dns.Msg
	item, err := m.backend.Get(key)
	if err != nil {
		err = KeyNotFound{key}
		return &msg, err
	}
	err = msg.Unpack(item.Value)
	if err != nil {
		err = SerializerError{err}
	}
	return &msg, err
}

func (m *MemcachedCache) Exists(key string) bool {
	_, err := m.backend.Get(key)
	if err != nil {
		return true
	}
	return false
}

func (m *MemcachedCache) Remove(key string) error {
	return m.backend.Delete(key)
}

func (m *MemcachedCache) Full() bool {
	// memcache is never full (LRU)
	return false
}
