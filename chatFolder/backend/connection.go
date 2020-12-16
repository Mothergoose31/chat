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
