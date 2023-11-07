package main

import (
	"avion-acc-gen/gen"
	imap "avion-acc-gen/pkg"
	"context"
	"log"
	"sync"
	"time"

	"github.com/emersion/go-imap/client"
	"golang.org/x/sync/semaphore"
)

var (
	err              error
	imapClient       *client.Client
	registeredAccs   *gen.AccountList = &gen.AccountList{}
	unregisteredAccs *gen.AccountList = &gen.AccountList{}
	emailCh                           = make(chan imap.EmailData)
	proxyList        []string
)

const (
	registeredFilepath   = "data/registered.txt"
	unregisteredFilepath = "data/unregistered.txt"
	proxyListFilepath    = "data/proxies.txt"
)

func init() {
	imapClient, err = imap.InitClient()
	if err != nil {
		log.Fatalf("error initializing IMAP client: %s", err)
	}
	proxyList, err = gen.LoadProxies(proxyListFilepath)
	if err != nil {
		log.Fatalf("error loading proxies: %s", err)
	}
	if err := imapClient.Login("", ""); err != nil { //! enter your imap email and password here
		log.Fatalf("error logging into IMAP client: %s", err)
	}
	if err := registeredAccs.LoadEmails(registeredFilepath); err != nil {
		log.Fatalf("error loading registered emails: %s", err)
	}
	if err := unregisteredAccs.LoadEmails(unregisteredFilepath); err != nil {
		log.Fatalf("error loading unregistered emails: %s", err)
	}
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go gen.EmailReader(imapClient, emailCh)

	sem := semaphore.NewWeighted(2) //! number of tasks to be ran concurrently
	var wg sync.WaitGroup
	for _, email := range unregisteredAccs.Emails {
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Fatalf("error acquiring semaphore: %s", err)
		}
		wg.Add(1)
		go func(email string) {
			defer wg.Done()
			defer sem.Release(1)
			s, err := gen.NewSession(email, unregisteredFilepath, registeredFilepath, proxyList)
			if err != nil {
				log.Fatalf("error creating session: %s", err)
			}
			s.RegisteredAccs = registeredAccs
			s.UnregisteredAccs = unregisteredAccs
			s.EmailCh = emailCh
			err = s.CreateAccount()
			if err != nil {
				s.Log.Warnf("error creating account: %s", err)
			}
		}(email)
		time.Sleep(5 * time.Second)
	}
	wg.Wait()

	imapClient.Logout()
}
