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
