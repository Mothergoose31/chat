package main

import (
	"bytes"
	"encoding/gob"
	_ "expvar"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	//_ "github.com/mkevac/debugcharts"
	conf "github.com/msbranco/goconfig"
)

type State struct {
	mutes   map[Userid]time.Time
	submode bool
	sync.RWMutex
}

var (
	state = &State{
		mutes: make(map[Userid]time.Time),
	}
)

const (
	WRITETIMEOUT         = 10 * time.Second
	READTIMEOUT          = time.Minute
	PINGINTERVAL         = 10 * time.Second
	PINGTIMEOUT          = 30 * time.Second
	MAXMESSAGESIZE       = 6144 // 512 max chars in a message, 8bytes per chars possible, plus factor in some protocol overhead
	SENDCHANNELSIZE      = 16
	BROADCASTCHANNELSIZE = 256
	DEFAULTBANDURATION   = time.Hour
	DEFAULTMUTEDURATION  = 10 * time.Minute
)

var (
	debuggingenabled = false
	DELAY            = 300 * time.Millisecond
	MAXTHROTTLETIME  = 5 * time.Minute
)



func main(){
	
	c, err := conf.ReadConfigFile("settings.cfg")
	if err != nil {
		nc := conf.NewConfigFile()
		nc.AddOption("default", "debug", "false")
		nc.AddOption("default", "listenaddress", ":9998")
		nc.AddOption("default", "maxprocesses", "0")
		nc.AddOption("default", "chatdelay", fmt.Sprintf("%d", 300*time.Millisecond))
		nc.AddOption("default", "maxthrottletime", fmt.Sprintf("%d", 5*time.Minute))
		nc.AddOption("default", "allowedoriginhost", "localhost")

		nc.AddSection("redis")
		nc.AddOption("redis", "address", "localhost:6379")
		nc.AddOption("redis", "database", "0")
		nc.AddOption("redis", "password", "")

		nc.AddSection("database")
		nc.AddOption("database", "type", "mysql")
		nc.AddOption("database", "dsn", "username:password@tcp(localhost:3306)/***website***?loc=UTC&parseTime=true&timeout=1s&time_zone=\"+00:00\"")

		nc.AddSection("api")
		nc.AddOption("api", "url", "============website needs to be set up =====================")
		nc.AddOption("api", "key", "changeme")

		if err := nc.WriteConfigFile("settings.cfg", 0644, "ChatBackend"); err != nil {
			log.Fatal("Unable to create settings.cfg: ", err)
		}
		if c, err = conf.ReadConfigFile("settings.cfg"); err != nil {
			log.Fatal("Unable to read settings.cfg: ", err)
		}
	}

	
	debuggingenabled, _ = c.GetBool("default", "debug")
	addr, _ := c.GetString("default", "listenaddress")
	processes, _ := c.GetInt64("default", "maxprocesses")
	delay, _ := c.GetInt64("default", "chatdelay")
	maxthrottletime, _ := c.GetInt64("default", "maxthrottletime")
	allowedoriginhost, _ := c.GetString("default", "allowedoriginhost")
	apiurl, _ := c.GetString("api", "url")
	apikey, _ := c.GetString("api", "key")
	DELAY = time.Duration(delay)
	MAXTHROTTLETIME = time.Duration(maxthrottletime)

	redisaddr, _ := c.GetString("redis", "address")
	redisdb, _ := c.GetInt64("redis", "database")
	redispw, _ := c.GetString("redis", "password")

	dbtype, _ := c.GetString("database", "type")
	dbdsn, _ := c.GetString("database", "dsn")

	if processes <= 0 {
		processes = int64(runtime.NumCPU())
	}
	runtime.GOMAXPROCS(int(processes))

	state.load()

	initApi(apiurl, apikey)
	initRedis(redisaddr, redisdb, redispw)

	initNamesCache()
	initHub()
	initDatabase(dbtype, dbdsn)

	initBroadcast(redisdb)
	initBans(redisdb)
	initUsers(redisdb)

	upgrader := websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header["Origin"]
			if len(origin) == 0 {
				return true
			}

			u, err := url.Parse(origin[0])
			if err != nil {
				return false
			}

			return allowedoriginhost == u.Host
		},
	}

}