package gen

import (
	"bufio"
	"os"
	"sync"
)

type AccountList struct {
	mu     sync.Mutex
	Emails []string
}

func (a *AccountList) LoadEmails(filename string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		email := scanner.Text()
		a.Emails = append(a.Emails, email)
	}

	return scanner.Err()
}

func (a *AccountList) removeEmail(emailToRemove string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for i, email := range a.Emails {
		if email == emailToRemove {
			a.Emails = append(a.Emails[:i], a.Emails[i+1:]...)
			return
		}
	}
}

func (a *AccountList) saveEmails(filename string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, email := range a.Emails {
		_, err := writer.WriteString(email + "\n")
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func (a *AccountList) addEmail(email string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.Emails = append(a.Emails, email)
}

func (s *Session) updateAccounts(email string) error {

	s.UnregisteredAccs.removeEmail(email)
	if err := s.UnregisteredAccs.saveEmails(s.UnregisteredPath); err != nil {
		return err
	}

	s.RegisteredAccs.addEmail(email)
	if err := s.RegisteredAccs.saveEmails(s.RegisteredPath); err != nil {
		return err
	}

	s.Log.Info("txt file updated")

	return nil
}
