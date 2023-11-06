package imap

import (
	"context"
	"fmt"
	"io"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

type EmailData struct {
	FromAddress string
	ToAddress   string
	SubjectData string
	BodyStr     string
}

// InitClient initializes an IMAP client and returns it
func InitClient() (*client.Client, error) {

	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// LoginClient logs into an IMAP client
func LoginClient(c *client.Client, emailAddress string, appPassword string) error {

	return c.Login(emailAddress, appPassword)
}

// LogoutClient logs out of an IMAP client
func LogoutClient(c *client.Client) error {

	return c.Logout()
}

func KeepLoggedIn(c *client.Client, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := c.Noop(); err != nil {
				return
			}
		}
	}
}

var section imap.BodySectionName

func GetMailRange(c *client.Client, numberOfMails uint32) ([]EmailData, error) {
	inbox, err := c.Select("INBOX", false)
	if err != nil {
		return nil, fmt.Errorf("error selecting INBOX: %w", err)
	}

	if inbox.Messages == 0 || numberOfMails == 0 {
		return []EmailData{}, nil
	}

	start := uint32(1)
	if numberOfMails < inbox.Messages {
		start = inbox.Messages - numberOfMails + 1
	}
	end := inbox.Messages

	seqset := new(imap.SeqSet)
	seqset.AddRange(start, end)

	messages := make(chan *imap.Message, numberOfMails)
	done := make(chan error, 1)
	defer close(done)

	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{section.FetchItem()}, messages)
	}()

	var emails []EmailData
	for msg := range messages {
		if msg == nil {
			continue
		}

		emailData, err := extractEmailData(msg)
		if err != nil {
			return nil, fmt.Errorf("error extracting email data: %w", err)
		}

		emails = append(emails, emailData)
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("error fetching emails: %w", err)
	}

	return emails, nil
}

func extractEmailData(msg *imap.Message) (EmailData, error) {
	var emailData EmailData
	r := msg.GetBody(&section)
	if r == nil {
		return emailData, fmt.Errorf("no body for message")
	}

	mr, err := mail.CreateReader(r)
	if err != nil {
		return emailData, err
	}
	defer mr.Close()

	header := mr.Header
	if from, err := header.AddressList("From"); err == nil && len(from) > 0 {
		emailData.FromAddress = from[0].Address
	}
	if to, err := header.AddressList("To"); err == nil && len(to) > 0 {
		emailData.ToAddress = to[0].Address
	}
	if subject, err := header.Subject(); err == nil {
		emailData.SubjectData = subject
	}

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return emailData, err
		}

		if _, ok := p.Header.(*mail.InlineHeader); ok {
			b, err := io.ReadAll(p.Body)
			if err != nil {
				return emailData, err
			}
			emailData.BodyStr += string(b) // Concatenate if there are multiple parts
		}
	}

	return emailData, nil
}
