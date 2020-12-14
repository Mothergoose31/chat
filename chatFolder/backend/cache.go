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



func initRedis(addr string, db int64, pw string) {
	var err error
	rds, err = redis.Open(
		redis.TcpConnection(addr, 1*time.Second),
		redis.Index(int(db), pw),
		redis.PoolSize(50),
	)
	if err != nil {
		F("Error making the redis pool", err)
	}

	conn := redisGetConn()
	defer conn.Return()

	rdsCircularBuffer, err = conn.DoString("SCRIPT", "LOAD", `
		local key = KEYS[1]
		local maxlength = tonumber(ARGV[1])
		local payload = ARGV[2]
		if not key then
			return {err = "INVALID KEY"}
		end
		if not payload then
			return {err = "INVALID PAYLOAD"}
		end
		if not maxlength then
			return {err = "INVALID MAXLENGTH"}
		end
		-- push the payload onto the end
		redis.call("RPUSH", key, payload)
		local delcount = 0
		-- get rid of excess lines from the front
		local numlines = redis.call("LLEN", key)
		for _ = numlines - 1, maxlength, -1 do
			redis.call("LPOP", key)
			delcount = delcount + 1
		end
		return delcount
	`)
	if err != nil {
		F("Circular buffer script loading error", err)
	}

