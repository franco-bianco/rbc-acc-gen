package gen

import (
	imap "avion-acc-gen/pkg"
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

type Session struct {
	Client  *resty.Client
	Ctx     context.Context
	Cancel  context.CancelFunc
	Log     *logrus.Entry
	EmailCh chan imap.EmailData

	UnregisteredPath string
	RegisteredPath   string
	RegisteredAccs   *AccountList
	UnregisteredAccs *AccountList
	ProxyList        []string

	state *state
}

type state struct {
	Email         string
	Password      string
	FirstName     string
	LastName      string
	RecoveryEmail string

	SecurityMetadata string
	RobState         string
	RobToken         string
	RefererURL       string
	SessionID        string
	SessionState     string
}

func NewSession(email, unregisteredPath, registeredPath string, proxyList []string) (*Session, error) {

	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: false,
		FullTimestamp:    true,
		TimestampFormat:  "02 Jan 06 15:04:05",
	})

	client := resty.New()
	parts := strings.Split(randomProxy(proxyList), ":")
	uri, _ := url.Parse(fmt.Sprintf("http://%s:%s@%s:%s", parts[2], parts[3], parts[0], parts[1]))
	client.SetTransport(&http.Transport{
		Proxy: http.ProxyURL(uri),
	})
	client.SetProxy(uri.String())

	s := &Session{
		Client:           client,
		Log:              log.WithField("email", email),
		UnregisteredPath: unregisteredPath,
		RegisteredPath:   registeredPath,
		state:            &state{Email: email},
	}

	return s, nil
}

func LoadProxies(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var proxies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxies = append(proxies, scanner.Text())
	}

	return proxies, scanner.Err()
}

func randomProxy(proxies []string) string {
	return proxies[rand.Intn(len(proxies))]
}
