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
	rdsGetIPCache, err = conn.DoString("SCRIPT", "LOAD", `
	local key = KEYS[1]
	return redis.call("ZRANGEBYSCORE", key, 1, 3)
`)
if err != nil {
	F("Get IP Cache script loading error", err)
}

rdsSetIPCache, err = conn.DoString("SCRIPT", "LOAD", `
	local key, value, maxlength = KEYS[1], ARGV[1], 3
	
	local count = redis.call("ZCOUNT", key, 1, maxlength)
	local existingscore = redis.call("ZSCORE", key, value)
	if existingscore then
		-- renumber all the elements and make this one the last
		local elements = redis.call("ZRANGEBYSCORE", key, 1, maxlength)
		local i = 1
		for _, v in ipairs(elements) do
			if v == value then
				redis.call("ZADD", key, count, v)
			else
				redis.call("ZADD", key, i, v)
				i = i + 1
			end
		end
		return
	end
	
	if count == maxlength then
		-- delete the first element, modify the other elements score down
		-- and add the new one to the end
		redis.call("ZREMRANGEBYSCORE", key, 1, 1)
		local elements = redis.call("ZRANGEBYSCORE", key, 2, maxlength)
		local i = 1
		for _, v in ipairs(elements) do
			redis.call("ZADD", key, i, v)
			i = i + 1
		end
		return redis.call("ZADD", key, count, value)
	else
		-- otherwise just insert it with the next score
		return redis.call("ZADD", key, count + 1, value)
	end
`)
	if err != nil {
	F("Set IP Cache script loading error", err)
	}
}






func cacheIPForUser(userid Userid, ip string) {
	if ip == "127.0.0.1" {
		return
	}

	conn := redisGetConn()
	defer conn.Return()

	_, err := conn.Do("EVALSHA", rdsSetIPCache, 1, fmt.Sprintf("CHAT:userips-%d", userid), ip)
	if err != nil {
		D("cacheIPForUser redis error", err)
	}
}
