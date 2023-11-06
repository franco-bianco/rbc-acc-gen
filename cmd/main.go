package main

import (
	"avion-acc-gen/gen"
	imap "avion-acc-gen/pkg"
	"context"
	"log"

	"github.com/emersion/go-imap/client"
)

var (
	err              error
	imapClient       *client.Client
	registeredAccs   *gen.AccountList = &gen.AccountList{}
	unregisteredAccs *gen.AccountList = &gen.AccountList{}
)

const (
	registeredPath   = "data/registered.txt"
	unregisteredPath = "data/unregistered.txt"
)

func init() {
	imapClient, err = imap.InitClient()
	if err != nil {
		log.Fatalf("error initializing IMAP client: %s", err)
	}
	if err := imapClient.Login("", ""); err != nil { //! ADD LOGIN CREDENTIALS
		log.Fatalf("error logging into IMAP client: %s", err)
	}
	if err := registeredAccs.LoadEmails(registeredPath); err != nil {
		log.Fatalf("error loading registered emails: %s", err)
	}
	if err := unregisteredAccs.LoadEmails(unregisteredPath); err != nil {
		log.Fatalf("error loading unregistered emails: %s", err)
	}
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go imap.KeepLoggedIn(imapClient, ctx)

	for _, email := range unregisteredAccs.Emails {
		s := gen.NewSession()
		s.RegisteredPath = registeredPath
		s.UnregisteredPath = unregisteredPath
		s.RegisteredAccs = registeredAccs
		s.UnregisteredAccs = unregisteredAccs
		s.IMAPClient = imapClient
		s.Log = s.Log.WithField("email", email)
		err := s.CreateAccount(email)
		if err != nil {
			s.Log.Warnf("error creating account: %s", err)
		}
	}

	imapClient.Logout()
}
