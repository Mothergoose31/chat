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

}