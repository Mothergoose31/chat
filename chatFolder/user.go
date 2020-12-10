package main

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tideland/golib/redis"
)
// pkg.go.dev

type userTools struct {
	nicklookup  map[string]*uidprot
	nicklock    sync.RWMutex
	featurelock sync.RWMutex
	features    map[uint64][]string
}

var (
	usertools = userTools{
		nicklookup: make(map[string]*uidprot),
		nicklock : sync.RWMutex{},
		featurelock : sync.RWMutex{},
		features : make(map[uint64][]string),
	}
)

const (
	ISADMIN = 1 << iota
	ISMODERATOR = 1 << iota
	ISVIP = 1 << iota
	ISPROTECTED = 1 << iota
	ISSUBSCRIBER  = 1 << iota
	ISBOT = 1 << iota

)	

type uidprot struct {
	id Userid
	protected bool
}

func initUsers(redisdb init64){
	go runRefresh(redisdb)
}

func(ut *UserTools) getUseridForNick(nick string) (userid, bool){
	ut.nicklock.RLock()
	d, ok := ut.nicklookup[strings.ToLower(nick)]
	if !ok {
		uid, protected := db.getUser(nick)
		if uid != 0 {
			ut.nicklock.RUnlock()
			ut.nicklock.Lock()
			ut.nicklookup[strings.ToLower(nick)] = &uidprot{uid, protected}
			ut.nicklock.Unlock()
			return uid, protected
		}
		ut.nicklock.RUnlock()
		return 0, false
	}
	ut.nicklock.RUnlock()
	return d.id, d.protected

}
func (ut *userTools) addUser(u *User, force bool) {
	lowernick := strings.ToLower(u.nick)
	if !force {
		ut.nicklock.RLock()
		_, ok := ut.nicklookup[lowernick]
		ut.nicklock.RUnlock()
		if ok {
			return
		}
	}
	ut.nicklock.Lock()
	defer ut.nicklock.Unlock()
	ut.nicklookup[lowernick] = &uidprot{u.id, u.isProtected()}
}

func runRefreshUser(redisdb int64) {
	setupRedisSubscription("refreshuser", redisdb, func(result *redis.PublishedValue) {
		user := userfromSession(result.Value.Bytes())
		namescache.refresh(user)
		hub.refreshuser <- user.id
	})
}


type Userid int32


type User struct {
	id              Userid
	nick            string
	features        uint64
	lastmessage     []byte
	lastmessagetime time.Time
	delayscale      uint8
	simplified      *SimplifiedUser
	connections     int32
	sync.RWMutex
}

type sessionuser struct {
	Username string   `json:"username"`
	UserId   string   `json:"userId"`
	Features []string `json:"features"`
}

func userfromSession(m []byte) (u *User) {
	
}