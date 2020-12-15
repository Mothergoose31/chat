package main

import (
	"database/sql"
	"net"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/tideland/golib/redis"
)

type Bans struct {
	users    map[Userid]time.Time
	userlock sync.RWMutex
	ips      map[string]time.Time
	userips  map[Userid][]string
	iplock   sync.RWMutex // protects both ips/userips
}
var (
	ipv6mask = net.CIDRMask(64, 128)
	bans     = Bans{
		make(map[Userid]time.Time),
		sync.RWMutex{},
		make(map[string]time.Time),
		make(map[Userid][]string),
		sync.RWMutex{},
	}
)

func getMaskedIP(s string) string {
	ip := net.ParseIP(s)
	if ip.To4() == nil {
		return ip.Mask(ipv6mask).String()
	} else {
		return s
	}
}

func initBans(redisdb int64) {
	go ban_user.run(redisdb)
}

func (b *Bans) run(redisdb int64) {
	b.loadActive()

	go b.runRefresh(redisdb)
	go b.runUnban(redisdb)

	t := time.NewTicker(time.Minute)

	for {
		select {
		case <-t.C:
			b.clean()
		}
	}
}

func (b *Bans) runRefresh(redisdb int64) {
	setupRedisSubscription("refreshbans", redisdb, func(result *redis.PublishedValue) {
		D("Refreshing bans")
		b.loadActive()
	})
}
