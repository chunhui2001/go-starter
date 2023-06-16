package gsession

import (
	"time"
)

// each session contains the username of the user and the time at which it expires
type session struct {
	username string
	expiry   time.Time
}

// we'll use this method later to determine if the session has expired
func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

var (
	// this map stores the users sessions. For larger scale applications, you can use a database or cache for this purpose
	sessions = map[string]session{}
)
