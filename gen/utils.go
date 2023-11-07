package gen

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/go-faker/faker/v4"
)

func (s *Session) randomProfile() error {
	parts := strings.Split(s.state.Email, "@")
	if len(parts) != 2 {
		return errors.New("invalid email")
	}
	randomNum := rand.Intn(1100) - 100
	s.state.RecoveryEmail = fmt.Sprintf("%s+%d@%s", parts[0], randomNum, parts[1])
	s.state.FirstName = faker.FirstName()
	s.state.LastName = faker.LastName()
	s.state.Password = "P@ssword1"
	return nil
}
