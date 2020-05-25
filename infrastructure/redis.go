package infrastructure

import (
	"fmt"
	"github.com/eb4uk/godns/models"
	"github.com/gomodule/redigo/redis"
)

func NewRedisConnection(settings models.RedisSettings) redis.Conn {
	var opts = []redis.DialOption{
		redis.DialDatabase(settings.DB),
		redis.DialPassword(settings.Password),
	}

	dial, err := redis.Dial("tcp", settings.Addr(), opts...)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return dial
}
