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
