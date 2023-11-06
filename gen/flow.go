package gen

import (
	"fmt"
)

func (s *Session) CreateAccount(emailAddress string) error {

	if err := s.randomProfile(emailAddress); err != nil {
		return fmt.Errorf("error generating random profile: %s", err)
	}

	if err := s.initSession(); err != nil {
		return fmt.Errorf("error initializing session: %s", err)
	}

	if err := s.initRegistration(); err != nil {
		return fmt.Errorf("error initializing registration: %s", err)
	}

	otpCode, err := s.getOTPCode(true)
	if err != nil {
		return fmt.Errorf("error getting OTP code: %s", err)
	}

	if err := s.submitOTP(otpCode); err != nil {
		return fmt.Errorf("error submitting OTP: %s", err)
	}

	if err := s.updateUserProfile(); err != nil {
		return fmt.Errorf("error updating user profile: %s", err)
	}

	if err := s.updateUserPassword(); err != nil {
		return fmt.Errorf("error updating user password: %s", err)
	}

	if err := s.agreeTerms(); err != nil {
		return fmt.Errorf("error agreeing to terms: %s", err)
	}

	if err := s.initSecondaryEmail(); err != nil {
		return fmt.Errorf("error initializing secondary email: %s", err)
	}

	otpCode, err = s.getOTPCode(false)
	if err != nil {
		return fmt.Errorf("error getting OTP code: %s", err)
	}

	if err := s.submitSecondaryOTP(otpCode); err != nil {
		return fmt.Errorf("error submitting secondary OTP: %s", err)
	}

	if err := s.completeRegistration(); err != nil {
		return fmt.Errorf("error completing registration: %s", err)
	}

	if err := s.updateAccounts(emailAddress); err != nil {
		return fmt.Errorf("error updating accounts: %s", err)
	}

	s.Log.Info("account created successfully")

	return nil
}
