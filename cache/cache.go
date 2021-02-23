package cache

import (
	"shiny_session/session"
	"sync"
	"time"
)

var sessions *cache

type cache struct {
	sync.Mutex
	sessions map[string]*session.Session
}

func initCache() {
	sessions = &cache{
		sessions: make(map[string]*session.Session),
	}
}

func (c *cache) Get(id string) (*session.Session, error) {
	c.Lock()
	defer c.Unlock()

	// Do we have a cached session?
	session_, ok := c.sessions[id]
	if !ok {
		// Not cached. Query the persistence layer for a session.
		var err error
		session_, err = session.Persistence.LoadSession(id)
		if err != nil {
			return nil, err
		}

		if session_ != nil {
			// Save it in the cache.
			if session.MaxSessionCacheSize != 0 {
				c.compact(1)
				c.sessions[id] = session_
			}

			// Store ID.
			session_.Lock()
			session_.Id = id
			session_.Unlock()
		}
	}

	return session_, nil
}

// Set inserts or updates a session in the cache. Since this is a write-through
// cache, the persistence layer is also triggered to save the session.
func (c *cache) Set(session_ *session.Session) error {
	c.Lock()
	defer c.Unlock()
	session_.Lock()
	session_.LastAccess = time.Now()
	id := session_.Id
	session_.Unlock()

	// Try to compact the cache.
	var requiredSpace int
	if _, ok := c.sessions[id]; !ok {
		requiredSpace = 1
	}
	c.compact(requiredSpace)

	// Save in cache.
	if session.MaxSessionCacheSize != 0 {
		c.sessions[id] = session_
	}

	// Write through to database.
	if err := session.Persistence.SaveSession(id, session_); err != nil {
		return err
	}

	return nil
}

func (c *cache) compact(requiredSpace int) (int, error) {
	// Check for old sessions.
	for id, session_ := range c.sessions {
		session_.RLock()
		age := time.Since(session_.LastAccess)
		session_.RUnlock()
		if age > session.SessionCacheExpiry {
			if err := session.Persistence.SaveSession(id, session_); err != nil {
				return 0, err
			}
			delete(c.sessions, id)
		}
	}

	// Cache may still grow.
	if session.MaxSessionCacheSize < 0 || len(c.sessions)+requiredSpace <= session.MaxSessionCacheSize {
		return 0, nil
	}

	// Drop the oldest sessions.
	var dropped int
	if requiredSpace > session.MaxSessionCacheSize {
		requiredSpace = session.MaxSessionCacheSize // We can't request more than is allowed.
	}
	for len(c.sessions)+requiredSpace > session.MaxSessionCacheSize {
		// Find oldest sessions and delete them.
		var (
			oldestAccessTime time.Time
			oldestSessionID  string
		)
		for id, session := range c.sessions {
			session.RLock()
			before := session.LastAccess.Before(oldestAccessTime)
			session.RUnlock()
			if oldestSessionID == "" || before {
				oldestSessionID = id
				oldestAccessTime = session.LastAccess
			}
		}
		if err := session.Persistence.SaveSession(oldestSessionID, c.sessions[oldestSessionID]); err != nil {
			return 0, err
		}
		delete(c.sessions, oldestSessionID)
		dropped++
	}

	return dropped, nil
}

func PurgeSessions() {
	sessions.Lock()
	defer sessions.Unlock()

	// Update all sessions in the database.
	for id, session_ := range sessions.sessions {
		session.Persistence.SaveSession(id, session_)
		// We only do this to update the last access time. Errors are not that
		// bad.
	}

	sessions.sessions = make(map[string]*session.Session, session.MaxSessionCacheSize)
}
