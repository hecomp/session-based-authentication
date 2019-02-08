package redisclient

import "github.com/gomodule/redigo/redis"

// Store the redis connection as a package level variable
var Cache redis.Conn

func InitCache() {
	// Iniitalize the redis connection to a redis instance running on your local machine
	conn, err := redis.DialURL("redis://localhost")
	if err != nil {
		panic(err)
	}
	// Assin the connection to the package level `cache` variable
	Cache = conn
}
