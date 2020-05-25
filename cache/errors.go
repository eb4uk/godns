package cache

import "fmt"

type KeyNotFound struct {
	key string
}

func (e KeyNotFound) Error() string {
	return e.key + " " + "not found"
}

type KeyExpired struct {
	Key string
}

func (e KeyExpired) Error() string {
	return e.Key + " " + "expired"
}

type CacheIsFull struct {
}

func (e CacheIsFull) Error() string {
	return "Cache is Full"
}

type SerializerError struct {
	err error
}

func (e SerializerError) Error() string {
	return fmt.Sprintf("Serializer error: got %v", e.err)
}
