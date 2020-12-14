package main

import (
	"fmt"
	"time"

	"github.com/tideland/golib/redis"
)

var (
	rds               *redis.Database
	rdsCircularBuffer string
	rdsGetIPCache     string
	rdsSetIPCache     string
)

// how many log lines to buffer for the scrollback
const CHATLOGLINES = 150

func redisGetConn() *redis.Connection {
again:
	conn, err := rds.Connection()
	if err != nil {
		D("Error getting a redis connection", err)
		if conn != nil {
			conn.Return()
		}
		time.Sleep(500 * time.Millisecond)
		goto again
	}

	return conn
}



