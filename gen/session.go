package gen

import (
	"context"

	"github.com/emersion/go-imap/client"
	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

type Session struct {
	Client     *resty.Client
	Ctx        context.Context
	Cancel     context.CancelFunc
	Log        *logrus.Entry
	IMAPClient *client.Client

	// Accounts
	UnregisteredPath string
	RegisteredPath   string
	RegisteredAccs   *AccountList
	UnregisteredAccs *AccountList

	state *state
}

type state struct {
	Email         string
	Password      string
	FirstName     string
	LastName      string
	RecoveryEmail string

	RobState     string
	RobToken     string
	RefererURL   string
	SessionID    string
	SessionState string
}

func NewSession() *Session {

	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: false,
		FullTimestamp:    true,
		TimestampFormat:  "02 Jan 06 15:04:05",
	})

	s := &Session{
		Client: resty.New(),
		Log:    logrus.NewEntry(log),
		state:  &state{},
	}

	return s
}
