package mini

import (
	"errors"
	"time"
)

type UserType string

const (
	Student UserType = "student"
	School  UserType = "school"
)

type Session struct {
	Session string
	Typ     UserType
	Id      int
	Expires time.Time
}

var Sessions []Session

func GetSession(sessionStr string) (*Session, error) {
	updateSessions()

	for i, s := range Sessions {
		if s.Session == sessionStr {
			// Prüfen, ob Session abgelaufen ist
			if time.Now().After(s.Expires) {
				return nil, errors.New("Session abgelaufen")
			}
			return &Sessions[i], nil
		}
	}

	return nil, errors.New("Session nicht gefunden")
}

func updateSessions() {
	var newSessions []Session

	for _, session := range Sessions {
		if session.Expires.After(time.Now()) {
			newSessions = append(newSessions, session)
		}
	}

	Sessions = newSessions
}
