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


func initNamesCache() {
}

func (nc *namesCache) getIrcNames() [][]string {
	nc.RLock()
	defer nc.RUnlock()
	return nc.ircnames
}

func (nc *namesCache) marshalNames(updateircnames bool) {
	users := make([]*SimplifiedUser, 0, len(nc.users))
	var allnames []string
	if updateircnames {
		allnames = make([]string, 0, len(nc.users))
	}
	for _, u := range nc.users {
		u.RLock()
		n := atomic.LoadInt32(&u.connections)
		if n <= 0 {
			continue
		}
		users = append(users, u.simplified)
		if updateircnames {
			prefix := ""
			switch {
			case u.featureGet(ISADMIN):
				prefix = "~" // +q
			case u.featureGet(ISBOT):
				prefix = "&" // +a
			case u.featureGet(ISMODERATOR):
				prefix = "@" // +o
			case u.featureGet(ISVIP):
				prefix = "%" // +h
			case u.featureGet(ISSUBSCRIBER):
				prefix = "+" // +v
			}
			allnames = append(allnames, prefix+u.nick)
		}
	}

	if updateircnames {
		l := 0
		var namelines [][]string
		var names []string
		for _, name := range allnames {
			if l+len(name) > 400 {
				namelines = append(namelines, names)
				l = 0
				names = nil
			}
			names = append(names, name)
			l += len(name)
		}
		nc.ircnames = namelines
	}

	n := NamesOut{
		Users:       users,
		Connections: nc.usercount,
	}
	nc.marshallednames, _ = n.MarshalJSON()

	cacheConnectedUsers(nc.marshallednames)

	for _, u := range nc.users {
		u.RUnlock()
	}
}
