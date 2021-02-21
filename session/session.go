package session

import (
	"fmt"
	"hash/fnv"
	"net/http"
	"sync"
	"time"
)

type Session struct {
	sync.RWMutex
	id                string                 // The session ID. Will not be saved with the session.
	user              User                   // The session user. If nil, no user is attached to this session.
	created           time.Time              // The time when this session was created.
	lastAccess        time.Time              // The last time the session was accessed through this API.
	lastIP            string                 // The remote address (IP:port) of the last request. If empty, it will not be compared.
	lastUserAgentHash uint64                 // A hash of the remote user agent string of the last request. If 0, it will not be compared.
	referenceID       string                 // If this session's ID was replaced, this is the ID of the newer session.
	data              map[string]interface{} // Any custom data stored in the session.
}

func Start(response http.ResponseWriter, request *http.Request, createIfNew bool) (*Session, error) {
	// We may need this hash later.
	var agentHash uint64
	hash := fnv.New64a()
	userAgent := request.Header.Get("User-Agent")
	if userAgent != "" {
		fmt.Fprint(hash, userAgent)
		agentHash = hash.Sum64()
	}
	//todo
	return
}
