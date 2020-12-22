package main

import (
	"time"
)

type Mutes int

var mutes Mutes

func (m *Mutes) clean() {
	state.Lock()
	defer state.Unlock()

	for uid, unmutetime := range state.mutes {
		if isExpiredUTC(unmutetime) {
			delete(state.mutes, uid)
		}
	}
	state.save()
}
