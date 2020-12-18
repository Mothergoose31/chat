package main

import (
	"bytes"
	"crypto/md5"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"github.com/gorilla/websocket"
)




// regexp to detect three or more consecutive characters intended to be combined
// with another char (like accents, diacritics), if there are more than 5

var invalidmessage = regexp.MustCompile(`\p{M}{5,}|[\p{Zl}\p{Zp}\x{202f}\x{00a0}]`)



type Connection struct {
	socket         *websocket.Conn
	ip             string
	send           chan *message
	sendmarshalled chan *message
	blocksend      chan *message
	banned         chan bool
	stop           chan bool
	user           *User
	ping           chan time.Time
	sync.RWMutex
}

type SimplifiedUser struct {
	Nick     string    `json:"nick,omitempty"`
	Features *[]string `json:"features,omitempty"`
}
ype EventDataIn struct {
	Data      string `json:"data"`
	Extradata string `json:"extradata"`
	Duration  int64  `json:"duration"`
}

type EventDataOut struct {
	*SimplifiedUser
	Targetuserid Userid `json:"-"`
	Timestamp    int64  `json:"timestamp"`
	Data         string `json:"data,omitempty"`
	Extradata    string `json:"extradata,omitempty"`
	Duration     int64  `json:"duration,omitempty"`
}

type BanIn struct {
	Nick        string `json:"nick"`
	BanIP       bool   `json:"banip"`
	Duration    int64  `json:"duration"`
	Ispermanent bool   `json:"ispermanent"`
	Reason      string `json:"reason"`
}

type PingOut struct {
	Timestamp int64 `json:"data"`
}

type message struct {
	msgtyp int
	event  string
	data   interface{}
}

type PrivmsgIn struct {
	Nick string `json:"nick"`
	Data string `json:"data"`
}

type PrivmsgOut struct {
	message
	targetuid Userid
	Messageid int64  `json:"messageid"`
	Timestamp int64  `json:"timestamp"`
	Nick      string `json:"nick,omitempty"`
	Data      string `json:"data,omitempty"`
}


func newConnection(s *websocket.Conn, user *User, ip string) {
	c := &Connection{
		socket:         s,
		ip:             ip,
		send:           make(chan *message, SENDCHANNELSIZE),
		sendmarshalled: make(chan *message, SENDCHANNELSIZE),
		blocksend:      make(chan *message),
		banned:         make(chan bool, 8),
		stop:           make(chan bool),
		user:           user,
		ping:           make(chan time.Time, 2),
		RWMutex:        sync.RWMutex{},
	}

	go c.writePumpText()
	c.readPumpText()
}

func (c *Connection) readPumpText() {
	defer func() {
		namescache.disconnect(c.user)
		c.Quit()
		c.socket.Close()
	}()

	c.socket.SetReadLimit(MAXMESSAGESIZE)
	c.socket.SetReadDeadline(time.Now().Add(READTIMEOUT))
	c.socket.SetPongHandler(func(string) error {
		c.socket.SetReadDeadline(time.Now().Add(PINGTIMEOUT))
		return nil
	})
	c.socket.SetPingHandler(func(string) error {
		c.sendmarshalled <- &message{
			msgtyp: websocket.PongMessage,
			event:  "PONG",
			data:   []byte{},
		}
		return nil
	})

	if c.user != nil {
		c.rlockUserIfExists()
		n := atomic.LoadInt32(&c.user.connections)
		if n > 5 {
			c.runlockUserIfExists()
			c.SendError("toomanyconnections")
			c.stop <- true
			return
		}
		c.runlockUserIfExists()
	} else {
		namescache.addConnection()
	}

	hub.register <- c
	c.Names()
	c.Join() // broadcast to the chat that a user has connected

	// Check mute status.
	muteTimeLeft := mutes.muteTimeLeft(c)
	if muteTimeLeft > time.Duration(0) {
		c.EmitBlock("ERR", NewMutedError(muteTimeLeft))
	}

	for {
		msgtype, message, err := c.socket.ReadMessage()
		if err != nil || msgtype == websocket.BinaryMessage {
			return
		}

		name, data, err := Unpack(string(message))
		if err != nil {
			// invalid protocol message from the client, just ignore it,
			// disconnect the user
			return
		}

		// dispatch
		switch name {
		case "MSG":
			c.OnMsg(data)
		case "MUTE":
			c.OnMute(data)
		case "UNMUTE":
			c.OnUnmute(data)
		case "BAN":
			c.OnBan(data)
		case "UNBAN":
			c.OnUnban(data)
		case "SUBONLY":
			c.OnSubonly(data)
		case "PING":
			c.OnPing(data)
		case "PONG":
			c.OnPong(data)
		case "BROADCAST":
			c.OnBroadcast(data)
		case "PRIVMSG":
			c.OnPrivmsg(data)
		}
	}
}


func (c *Connection) write(mt int, payload []byte) error {
	c.socket.SetWriteDeadline(time.Now().Add(WRITETIMEOUT))
	return c.socket.WriteMessage(mt, payload)
}
func (c *Connection) writePumpText() {
	defer func() {
		// need to write hub.go
		hub.unregister <- c
		c.socket.Close() // Necessary to force reading to stop, will start the cleanup
	}()
	for {
		select {
		case _, ok := <-c.ping:
			if !ok {
				return
			}
			m, _ := time.Now().MarshalBinary()
			if err := c.write(websocket.PingMessage, m); err != nil {
				return
			}
		case <-c.banned:
			c.write(websocket.TextMessage, []byte(`ERR {"description":"banned"}`))
			c.write(websocket.CloseMessage, []byte{})
			return
		case <-c.stop:
			return
		case m := <-c.blocksend:
			c.rlockUserIfExists()
			if data, err := Marshal(m.data); err == nil {
				c.runlockUserIfExists()
				if data, err := Pack(m.event, data); err == nil {
					if err := c.write(websocket.TextMessage, data); err != nil {
						return
					}
				}
			} else {
				c.runlockUserIfExists()
			}
		case m := <-c.send:
			c.rlockUserIfExists()
			if data, err := Marshal(m.data); err == nil {
				c.runlockUserIfExists()
				if data, err := Pack(m.event, data); err == nil {
					typ := m.msgtyp
					if typ == 0 {
						typ = websocket.TextMessage
					}
					if err := c.write(typ, data); err != nil {
						return
					}
				}
			} else {
				c.runlockUserIfExists()
			}
		case message := <-c.sendmarshalled:
			data := message.data.([]byte)
			if data, err := Pack(message.event, data); err == nil {
				typ := message.msgtyp
				if typ == 0 {
					typ = websocket.TextMessage
				}
				if err := c.write(typ, data); err != nil {
					return
				}
			}
		}
	}
}
func (c *Connection) rlockUserIfExists() {
	if c.user == nil {
		return
	}

	c.user.RLock()
}

func (c *Connection) runlockUserIfExists() {
	if c.user == nil {
		return
	}

	c.user.RUnlock()
}

func (c *Connection) Emit(event string, data interface{}) {
	c.send <- &message{
		event: event,
		data:  data,
	}
}

func (c *Connection) EmitBlock(event string, data interface{}) {
	c.blocksend <- &message{
		event: event,
		data:  data,
	}
	return
}


func (c *Connection) Broadcast(event string, data *EventDataOut) {
	c.rlockUserIfExists()
	marshalled, _ := Marshal(data)
	c.runlockUserIfExists()

	m := &message{
		event: event,
		data:  marshalled,
	}
	hub.broadcast <- m
}

func (c *Connection) canModerateUser(nick string) (bool, Userid) {
	if c.user == nil || utf8.RuneCountInString(nick) == 0 {
		return false, 0
	}

	uid, protected := usertools.getUseridForNick(nick)
	if uid == 0 || c.user.id == uid || protected {
		return false, uid
	}

	return true, uid
}

func (c *Connection) getEventDataOut() *EventDataOut {
	out := &EventDataOut{
		Timestamp: unixMilliTime(),
	}
	if c.user != nil {
		out.SimplifiedUser = c.user.simplified
	}
	return out
}

func (c *Connection) Join() {
	if c.user != nil {
		c.rlockUserIfExists()
		defer c.runlockUserIfExists()
		n := atomic.LoadInt32(&c.user.connections)
		if n == 1 {
			c.Broadcast("JOIN", c.getEventDataOut())
		}
	}
}

func (c *Connection) Quit() {
	if c.user != nil {
		c.rlockUserIfExists()
		defer c.runlockUserIfExists()
		n := atomic.LoadInt32(&c.user.connections)
		if n <= 0 {
			c.Broadcast("QUIT", c.getEventDataOut())
		}
	}
}


func (c *Connection) OnBroadcast(data []byte) {
	m := &EventDataIn{}
	if err := Unmarshal(data, m); err != nil {
		c.SendError("protocolerror")
		return
	}

	if c.user == nil {
		c.SendError("needlogin")
		return
	}

	if !c.user.featureGet(ISADMIN) {
		c.SendError("nopermission")
		return
	}

	msg := strings.TrimSpace(m.Data)
	msglen := utf8.RuneCountInString(msg)
	if !utf8.ValidString(msg) || msglen == 0 || msglen > 512 || invalidmessage.MatchString(msg) {
		c.SendError("invalidmsg")
		return
	}

	out := c.getEventDataOut()
	out.Data = msg
	c.Broadcast("BROADCAST", out)

}


func (c *Connection) Join() {
	if c.user != nil {
		c.rlockUserIfExists()
		defer c.runlockUserIfExists()
		n := atomic.LoadInt32(&c.user.connections)
		if n == 1 {
			c.Broadcast("JOIN", c.getEventDataOut())
		}
	}
}

func (c *Connection) Quit() {
	if c.user != nil {
		c.rlockUserIfExists()
		defer c.runlockUserIfExists()
		n := atomic.LoadInt32(&c.user.connections)
		if n <= 0 {
			c.Broadcast("QUIT", c.getEventDataOut())
		}
	}
}
