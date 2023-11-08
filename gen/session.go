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

	s := &Session{
		Client:           resty.New(),
		Log:              log.WithField("email", email),
		UnregisteredPath: unregisteredPath,
		RegisteredPath:   registeredPath,
		state:            &state{Email: email},
	}

	if err := s.setClientProxy(proxyList); err != nil {
		return nil, err
	}

	return s, nil
}

// setClientProxy sets the proxy for the resty client
func (s *Session) setClientProxy(proxyList []string) error {

	proxy := proxyList[rand.Intn(len(proxyList))]
	parts := strings.Split(proxy, ":")

	uri, _ := url.Parse(fmt.Sprintf("http://%s:%s@%s:%s", parts[2], parts[3], parts[0], parts[1]))
	s.Client.SetTransport(&http.Transport{
		Proxy: http.ProxyURL(uri),
	})
	s.Client.SetProxy(uri.String())
	return nil
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

	if len(proxies) == 0 {
		return nil, fmt.Errorf("no proxies found in file")
	}

	return proxies, scanner.Err()
}
