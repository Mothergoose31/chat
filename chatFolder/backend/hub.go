  
package main

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/tideland/golib/redis"
)

type Hub struct {
	connections map[*Connection]bool
	broadcast   chan *message
	privmsg     chan *PrivmsgOut
	register    chan *Connection
	unregister  chan *Connection
	bans        chan Userid
	ipbans      chan string
	getips      chan useridips
	users       map[Userid]*User
	refreshuser chan Userid
}

type useridips struct {
	userid Userid
	c      chan []string
}

var hub = Hub{
	connections: make(map[*Connection]bool),
	broadcast:   make(chan *message, BROADCASTCHANNELSIZE),
	privmsg:     make(chan *PrivmsgOut, BROADCASTCHANNELSIZE),
	register:    make(chan *Connection, 256),
	unregister:  make(chan *Connection),
	bans:        make(chan Userid, 4),
	ipbans:      make(chan string, 4),
	getips:      make(chan useridips),
	users:       make(map[Userid]*User),
	refreshuser: make(chan Userid, 4),
}

func initHub() {
	go hub.run()
}
