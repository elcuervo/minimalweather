package db

import (
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	server   = os.Getenv("REDIS_URL")
	password = ""
)

var Pool = &redis.Pool{
	MaxIdle:     10,
	IdleTimeout: 240 * time.Second,
	Dial: func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", server)
		if err != nil {
			return nil, err
		}
		if len(password) > 0 {
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
		}
		return c, err
	},

	TestOnBorrow: func(c redis.Conn, t time.Time) error {
		_, err := c.Do("PING")
		return err
	},
}
