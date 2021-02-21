package cache

import (
	"shiny_session/session"
	"sync"
)

type cache struct {
	sync.Mutex
	sessions map[string]*session.Session
}
