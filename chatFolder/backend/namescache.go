package main

import (
	"sync"
	"sync/atomic"
)

// ffjson: skip
type namesCache struct {
	users           map[Userid]*User
	marshallednames []byte
	usercount       uint32
	ircnames        [][]string
	sync.RWMutex
}

// ffjson: skip
type userChan struct {
	user *User
	c    chan *User
}

type NamesOut struct {
	Users       []*SimplifiedUser `json:"users"`
	Connections uint32            `json:"connectioncount"`
}

var namescache = namesCache{
	users:   make(map[Userid]*User),
	RWMutex: sync.RWMutex{},
}
