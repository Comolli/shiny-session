package config

import (
	"math"
	"net/http"
	"shiny_session/session"
	"time"
)

var (
	Persistence             session.PersistenceLayer = session.ExtendablePersistenceLayer{}
	SessionExpiry           time.Duration            = math.MaxInt64
	SessionIDExpiry                                  = time.Hour
	SessionIDGracePeriod                             = 5 * time.Minute
	AcceptRemoteIP                                   = 1
	AcceptChangingUserAgent                          = false
	SessionCookie                                    = "id"
	NewSessionCookie                                 = func() *http.Cookie {
		return &http.Cookie{ // Default lifetime is 10 years (i.e. forever).
			Expires:  time.Now().Add(10 * 365 * 24 * time.Hour), // For IE, other browsers will use MaxAge.
			MaxAge:   10 * 365 * 24 * 60 * 60,
			HttpOnly: true,

			// Uncomment and edit the following fields for production use:
			//Domain: "www.example.com",
			//Path:   "/",
			//Secure: true,
		}
	}
	MaxSessionCacheSize = 1024 * 1024
	SessionCacheExpiry  = time.Hour

	MutexMaxCacheSize     = 1024 * 1024
	MutexCleanupFrequency = 10 * time.Minute
	MutexStaleMutexes     = time.Hour
)
