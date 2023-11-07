package gen

import (
	imap "avion-acc-gen/pkg"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/emersion/go-imap/client"
)

const (
	clientMetadata   = "{\"appCode\":\"RH90\",\"appOrg\":\"com.rbc.fg.intnet.bnk.public\",\"appVersion\":\"v1.2.3\",\"assetId\":\"YSMN943L\",\"channelIdentifier\":\"olb\",\"ipAddress\":\"10.40.50.60\",\"language\":\"en\",\"physicalLocationId\":\"00075\",\"legacyId\":\"223456789\",\"requestUniqueId\":\"D83A54AC-7DDD-44C3-AB00-9D31128224EC\",\"timeZoneOffset\":\"-5\",\"originatingComputerCentreCode\":\"GCC\"}"
	securityMetadata = "{\"ipAddress\":\"\",\"locale\":\"en\",\"userAgent\":\"Mozilla\\/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit\\/537.36 (KHTML, like Gecko) Chrome\\/119.0.0.0 Safari\\/537.36\",\"screenHeight\":1080,\"screenWidth\":1920,\"timeZoneOffset\":\"+08:00\",\"everCookies\":\"0\",\"webRTCMedium\":\"0\",\"mouseAndKeystrokeDynamics\":\"0\",\"countingHostsBehindNAT\":\"0\",\"canvasFingerprinting\":\"c796223f9c313be6c6d2854d507a8983\"}"
)

func (s *Session) initSession() error {

	s.Log.Info("init session...")
	res, err := s.Client.R().
		SetHeaders(map[string]string{
			"host":                      "www1.avionrewards.com",
			"sec-ch-ua-mobile":          "?0",
			"upgrade-insecure-requests": "1",
			"user-agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
			"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
			"sec-fetch-site":            "same-origin",
			"sec-fetch-mode":            "navigate",
			"sec-fetch-dest":            "document",
			"accept-language":           "en-US,en;q=0.9",
		}).
		SetQueryParam("PolicyId", "urn:ibm:security:authentication:asf:avionMpolicy").
		Get("https://www1.avionrewards.com/mga/sps/authsvc")
	if err != nil {
		return err
	}
	if res.StatusCode() != 200 {
		return fmt.Errorf("status code %d", res.StatusCode())
	}

	re := regexp.MustCompile(`url=(https:\/\/[^\s"]+)`)
	matches := re.FindStringSubmatch(string(res.Body()))
	if len(matches) != 2 {
		return fmt.Errorf("failed to find url")
	}
	s.state.RefererURL = matches[1]

	re = regexp.MustCompile(`PARM2=([^&]+)`)
	matches = re.FindStringSubmatch(s.state.RefererURL)
	if len(matches) != 2 {
		return fmt.Errorf("failed to find PARM2")
	}
	s.state.RobState = matches[1]

	res, err = s.Client.R().
		SetHeaders(map[string]string{
			"host":             "ssoa.rbc.com",
			"cookie":           "_gcl_au=1.1.1544628081.1699173102",
			"accept":           "application/json, text/plain, */*",
			"sec-ch-ua-mobile": "?0",
			"user-agent":       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
			"sec-fetch-site":   "same-origin",
			"sec-fetch-mode":   "cors",
			"sec-fetch-dest":   "empty",
			"referer":          s.state.RefererURL,
			"accept-language":  "en-US,en;q=0.9",
		}).
		Get(fmt.Sprintf("https://ssoa.rbc.com/riam/ui/v2/rob/states/%s", s.state.RobState))
	if err != nil {
		return err
	}
	if res.StatusCode() != 200 {
		return fmt.Errorf("status code %d", res.StatusCode())
	}

	var data RobTokenRes
	if err := json.Unmarshal(res.Body(), &data); err != nil {
		return err
	}
	s.state.RobToken = data.RequestParameters.AccessToken

	res, err = s.Client.R().
		SetHeaders(map[string]string{
			"host":              "ssoa.rbc.com",
			"locale":            "en",
			"sec-ch-ua-mobile":  "?0",
			"user-agent":        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
			"client-metadata":   clientMetadata,
			"rob-token":         s.state.RobToken,
			"accept":            "application/json, text/plain, */*",
			"rob-state":         s.state.RobState,
			"security-metadata": securityMetadata,
			"trace-id":          "3835416600142",
			"sec-fetch-site":    "same-origin",
			"sec-fetch-mode":    "cors",
			"sec-fetch-dest":    "empty",
			"referer":           s.state.RefererURL,
			"accept-language":   "en-US,en;q=0.9",
		}).
		Get("https://ssoa.rbc.com/riam/ui/v2/idp/policy")
	if err != nil {
		return err
	}
	if res.StatusCode() != 200 {
		return fmt.Errorf("status code %d", res.StatusCode())
	}

	s.state.SessionID = res.Header().Get("session-id")
	s.state.SessionState = res.Header().Get("session-state")

	res, err = s.Client.R().
		SetHeaders(map[string]string{
			"host":              "ssoa.rbc.com",
			"session-id":        s.state.SessionID,
			"locale":            "en",
			"sec-ch-ua-mobile":  "?0",
			"session-state":     s.state.SessionState,
			"user-agent":        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
			"client-metadata":   clientMetadata,
			"content-type":      "application/json",
			"accept":            "application/json, text/plain, */*",
			"security-metadata": securityMetadata,
			"trace-id":          "4849850101864",
			"origin":            "https://ssoa.rbc.com",
			"sec-fetch-site":    "same-origin",
			"sec-fetch-mode":    "cors",
			"sec-fetch-dest":    "empty",
			"referer":           s.state.RefererURL,
			"accept-language":   "en-US,en;q=0.9",
		}).
		SetBody(`{"requestType":"registration"}`).
		Post("https://ssoa.rbc.com/riam/ui/v2/idp/request-type")
	if err != nil {
		return err
	}
	if res.StatusCode() != 200 {
		return fmt.Errorf("status code %d", res.StatusCode())
	}

	s.state.SessionID = res.Header().Get("session-id")
	s.state.SessionState = res.Header().Get("session-state")

	return nil
}

func (s *Session) initRegistration() error {

	s.Log.Info("beginning registration...")
	res, err := s.Client.R().
		SetHeaders(map[string]string{
			"host":              "ssoa.rbc.com",
			"session-id":        s.state.SessionID,
			"client-metadata":   clientMetadata,
			"locale":            "en",
			"sec-ch-ua-mobile":  "?0",
			"session-state":     s.state.SessionState,
			"user-agent":        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
			"content-type":      "application/json",
			"accept":            "application/json, text/plain, */*",
			"security-metadata": securityMetadata,
			"trace-id":          "3835416600142",
			"origin":            "https://ssoa.rbc.com",
			"sec-fetch-site":    "same-origin",
			"sec-fetch-mode":    "cors",
			"sec-fetch-dest":    "empty",
			"referer":           s.state.RefererURL,
			"accept-language":   "en-US,en;q=0.9",
		}).
		SetBody(fmt.Sprintf(`{"credentialType":"email","credential":"%s"}`, s.state.Email)).
		Post("https://ssoa.rbc.com/riam/ui/v2/idp/primary-credential/code")
	if err != nil {
		return err
	}
	if res.StatusCode() != 200 {
		fmt.Println(string(res.Body()))
		return fmt.Errorf("status code %d", res.StatusCode())
	}

	var data ResStatus
	if err := json.Unmarshal(res.Body(), &data); err != nil {
		return err
	}

	if !strings.EqualFold(data.Message, "OTP Sent Successfully.") {
		return fmt.Errorf("unexpected message: %s", data.Message)
	}

	s.state.SessionID = res.Header().Get("session-id")
	s.state.SessionState = res.Header().Get("session-state")

	return nil
}

func (s *Session) submitOTP(otpCode string) error {

	s.Log.Info("submitting OTP...")
	res, err := s.Client.R().
		SetHeaders(map[string]string{
			"host":              "ssoa.rbc.com",
			"session-id":        s.state.SessionID,
			"client-metadata":   clientMetadata,
			"locale":            "en",
			"sec-ch-ua-mobile":  "?0",
			"session-state":     s.state.SessionState,
			"user-agent":        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
			"content-type":      "application/json",
			"accept":            "application/json, text/plain, */*",
			"security-metadata": securityMetadata,
			"trace-id":          "4849850101864",
			"origin":            "https://ssoa.rbc.com",
			"sec-fetch-site":    "same-origin",
			"sec-fetch-mode":    "cors",
			"sec-fetch-dest":    "empty",
			"referer":           s.state.RefererURL,
			"accept-language":   "en-US,en;q=0.9",
		}).
		SetBody(fmt.Sprintf(`{"code":"%s"}`, otpCode)).
		Post("https://ssoa.rbc.com/riam/ui/v2/idp/primary-credential/code/verify")
	if err != nil {
		return err
	}
	if res.StatusCode() != 200 {
		fmt.Println(string(res.Body()))
		return fmt.Errorf("status code %d", res.StatusCode())
	}

	var data SubmitOTPRes
	if err := json.Unmarshal(res.Body(), &data); err != nil {
		return err
	}

	if !data.Verified {
		return fmt.Errorf("OTP not verified")
	}

	s.state.SessionID = res.Header().Get("session-id")
	s.state.SessionState = res.Header().Get("session-state")

	return nil
}

func (s *Session) updateUserProfile() error {

	s.Log.Info("updating user profile...")
	res, err := s.Client.R().
		SetHeaders(map[string]string{
			"host":              "ssoa.rbc.com",
			"session-id":        s.state.SessionID,
			"client-metadata":   clientMetadata,
			"locale":            "en",
			"sec-ch-ua-mobile":  "?0",
			"session-state":     s.state.SessionState,
			"user-agent":        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
			"content-type":      "application/json",
			"accept":            "application/json, text/plain, */*",
			"security-metadata": securityMetadata,
			"trace-id":          "4849850101864",
			"origin":            "https://ssoa.rbc.com",
			"sec-fetch-site":    "same-origin",
			"sec-fetch-mode":    "cors",
			"sec-fetch-dest":    "empty",
			"referer":           s.state.RefererURL,
			"accept-language":   "en-US,en;q=0.9",
		}).
		SetBody(fmt.Sprintf(`{"firstName":"%s","lastName":"%s"}`, s.state.FirstName, s.state.LastName)).
		Post("https://ssoa.rbc.com/riam/ui/v2/idp/user/profile")
	if err != nil {
		return err
	}
	if res.StatusCode() != 200 {
		return fmt.Errorf("status code %d", res.StatusCode())
	}

	var data ResStatus
	if err := json.Unmarshal(res.Body(), &data); err != nil {
		return err
	}

	if !strings.EqualFold(data.Status, "Success") {
		return fmt.Errorf("failed to update profile: %s", data.Message)
	}

	s.state.SessionID = res.Header().Get("session-id")
	s.state.SessionState = res.Header().Get("session-state")

	return nil
}

func (s *Session) updateUserPassword() error {

	s.Log.Info("getting public key...")
	res, err := s.Client.R().
		SetHeaders(map[string]string{
			"host":              "ssoa.rbc.com",
			"session-id":        s.state.SessionID,
			"client-metadata":   clientMetadata,
			"locale":            "en",
			"sec-ch-ua-mobile":  "?0",
			"session-state":     s.state.SessionState,
			"user-agent":        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
			"accept":            "application/json, text/plain, */*",
			"security-metadata": securityMetadata,
			"trace-id":          "4849850101864",
			"sec-fetch-site":    "same-origin",
			"sec-fetch-mode":    "cors",
			"sec-fetch-dest":    "empty",
			"referer":           s.state.RefererURL,
			"accept-language":   "en-US,en;q=0.9"},
		).
		Get("https://ssoa.rbc.com/riam/ui/v2/idp/publickey")
	if err != nil {
		return err
	}
	if res.StatusCode() != 200 {
		return fmt.Errorf("status code %d", res.StatusCode())
	}
	var publicKeyRes PublicKeyRes
	if err := json.Unmarshal(res.Body(), &publicKeyRes); err != nil {
		return err
	}

	s.state.SessionID = res.Header().Get("session-id")
	s.state.SessionState = res.Header().Get("session-state")

	encryptedPassword, err := encryptPassword(s.state.Password, publicKeyRes.PublicKey)
	if err != nil {
		return err
	}

	s.Log.Info("updating user pass...")
	res, err = s.Client.R().
		SetHeaders(map[string]string{
			"host":              "ssoa.rbc.com",
			"session-id":        s.state.SessionID,
			"client-metadata":   clientMetadata,
			"locale":            "en",
			"sec-ch-ua-mobile":  "?0",
			"session-state":     s.state.SessionState,
			"user-agent":        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
			"content-type":      "application/json",
			"accept":            "application/json, text/plain, */*",
			"security-metadata": securityMetadata,
			"trace-id":          "4849850101864",
			"origin":            "https://ssoa.rbc.com",
			"sec-fetch-site":    "same-origin",
			"sec-fetch-mode":    "cors",
			"sec-fetch-dest":    "empty",
			"referer":           s.state.RefererURL,
			"accept-language":   "en-US,en;q=0.9",
		}).
		SetBody(fmt.Sprintf(`{"password":"%s"}`, encryptedPassword)).
		Post("https://ssoa.rbc.com/riam/ui/v2/idp/user/password")
	if err != nil {
		return err
	}
	if res.StatusCode() != 200 {
		return fmt.Errorf("status code %d", res.StatusCode())
	}

	var updatePasswordStatusRes ResStatus
	if err := json.Unmarshal(res.Body(), &updatePasswordStatusRes); err != nil {
		return err
	}

	if !strings.EqualFold(updatePasswordStatusRes.Status, "Success") {
		return fmt.Errorf("failed to update password: %s", updatePasswordStatusRes.Message)
	}

	s.state.SessionID = res.Header().Get("session-id")
	s.state.SessionState = res.Header().Get("session-state")

	return nil
}

func (s *Session) agreeTerms() error {

	s.Log.Info("agreeing to terms...")
	res, err := s.Client.R().
		SetHeaders(map[string]string{
			"host":              "ssoa.rbc.com",
			"session-id":        s.state.SessionID,
			"client-metadata":   clientMetadata,
			"locale":            "en",
			"sec-ch-ua-mobile":  "?0",
			"session-state":     s.state.SessionState,
			"user-agent":        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
			"content-type":      "application/json",
			"accept":            "application/json, text/plain, */*",
			"security-metadata": securityMetadata,
			"trace-id":          "4849850101864",
			"origin":            "https://ssoa.rbc.com",
			"sec-fetch-site":    "same-origin",
			"sec-fetch-mode":    "cors",
			"sec-fetch-dest":    "empty",
			"referer":           s.state.RefererURL,
			"accept-language":   "en-US,en;q=0.9",
		}).
		SetBody(`{"riamUuid":"","consentType":"AVION"}`).
		Post("https://ssoa.rbc.com/riam/ui/v2/idp/consent/termsAndConditions")
	if err != nil {
		return err
	}
	if res.StatusCode() != 200 {
		fmt.Println(string(res.Body()))
		return fmt.Errorf("status code %d", res.StatusCode())
	}

	s.state.SessionID = res.Header().Get("session-id")
	s.state.SessionState = res.Header().Get("session-state")

	return nil
}

func (s *Session) initSecondaryEmail() error {

	s.Log.Info("adding secondary email...")
	res, err := s.Client.R().
		SetHeaders(map[string]string{
			"host":              "ssoa.rbc.com",
			"session-id":        s.state.SessionID,
			"client-metadata":   clientMetadata,
			"locale":            "en",
			"sec-ch-ua-mobile":  "?0",
			"session-state":     s.state.SessionState,
			"user-agent":        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
			"content-type":      "application/json",
			"accept":            "application/json, text/plain, */*",
			"security-metadata": securityMetadata,
			"trace-id":          "4849850101864",
			"origin":            "https://ssoa.rbc.com",
			"sec-fetch-site":    "same-origin",
			"sec-fetch-mode":    "cors",
			"sec-fetch-dest":    "empty",
			"referer":           s.state.RefererURL,
			"accept-language":   "en-US,en;q=0.9",
		}).
		SetBody(fmt.Sprintf(`{"credentialType":"email","credential":"%s"}`, s.state.RecoveryEmail)).
		Post("https://ssoa.rbc.com/riam/ui/v2/idp/secondary-credential/code")
	if err != nil {
		return err
	}
	if res.StatusCode() != 200 {
		fmt.Println(string(res.Body()))
		return fmt.Errorf("status code %d", res.StatusCode())
	}

	var data ResStatus
	if err := json.Unmarshal(res.Body(), &data); err != nil {
		return err
	}

	if !strings.EqualFold(data.Message, "OTP Sent Successfully.") {
		return fmt.Errorf("unexpected message: %s", data.Message)
	}

	s.state.SessionID = res.Header().Get("session-id")
	s.state.SessionState = res.Header().Get("session-state")

	return nil
}

func (s *Session) submitSecondaryOTP(otpCode string) error {

	s.Log.Info("submitting secondary email otp...")
	res, err := s.Client.R().
		SetHeaders(map[string]string{
			"host":              "ssoa.rbc.com",
			"session-id":        s.state.SessionID,
			"client-metadata":   clientMetadata,
			"locale":            "en",
			"sec-ch-ua-mobile":  "?0",
			"session-state":     s.state.SessionState,
			"user-agent":        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
			"content-type":      "application/json",
			"accept":            "application/json, text/plain, */*",
			"security-metadata": securityMetadata,
			"trace-id":          "4849850101864",
			"origin":            "https://ssoa.rbc.com",
			"sec-fetch-site":    "same-origin",
			"sec-fetch-mode":    "cors",
			"sec-fetch-dest":    "empty",
			"referer":           s.state.RefererURL,
			"accept-language":   "en-US,en;q=0.9",
		}).
		SetBody(fmt.Sprintf(`{"code":"%s"}`, otpCode)).
		Post("https://ssoa.rbc.com/riam/ui/v2/idp/secondary-credential/code/verify")
	if err != nil {
		return err
	}
	if res.StatusCode() != 200 {
		return fmt.Errorf("status code %d", res.StatusCode())
	}

	var data SubmitOTPRes
	if err := json.Unmarshal(res.Body(), &data); err != nil {
		return err
	}

	if !data.Verified {
		return fmt.Errorf("OTP not verified")
	}

	s.state.SessionID = res.Header().Get("session-id")
	s.state.SessionState = res.Header().Get("session-state")

	return nil
}

func (s *Session) completeRegistration() error {

	s.Log.Info("completing registration...")
	res, err := s.Client.R().
		SetHeaders(map[string]string{
			"host":              "ssoa.rbc.com",
			"session-id":        s.state.SessionID,
			"client-metadata":   clientMetadata,
			"locale":            "en",
			"sec-ch-ua-mobile":  "?0",
			"session-state":     s.state.SessionState,
			"user-agent":        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
			"content-type":      "application/json",
			"accept":            "application/json, text/plain, */*",
			"security-metadata": securityMetadata,
			"trace-id":          "4849850101864",
			"origin":            "https://ssoa.rbc.com",
			"sec-fetch-site":    "same-origin",
			"sec-fetch-mode":    "cors",
			"sec-fetch-dest":    "empty",
			"referer":           s.state.RefererURL,
			"accept-language":   "en-US,en;q=0.9",
		}).
		SetBody(`{"userEmailConsents":false}`).
		Post("https://ssoa.rbc.com/riam/ui/v2/idp/user")
	if err != nil {
		return err
	}
	if res.StatusCode() != 200 {
		fmt.Println(string(res.Body()))
		return fmt.Errorf("status code %d", res.StatusCode())
	}

	var data CompleteRegistrationRes
	if err := json.Unmarshal(res.Body(), &data); err != nil {
		return err
	}

	if !strings.EqualFold(data.Result, "success") {
		return fmt.Errorf("failed to complete registration: %s", data.RedirectURL)
	}

	return nil
}

func (s *Session) getOTPCode(isPrimary bool) (string, error) {
	s.Log.Info("awaiting OTP code...")
	var targetEmail string
	if isPrimary {
		targetEmail = s.state.Email
	} else {
		targetEmail = s.state.RecoveryEmail
	}
	timeout := time.After(15 * time.Minute)
	for {
		select {
		case <-timeout:
			return "", fmt.Errorf("timed out waiting for OTP code")
		case emailData := <-s.EmailCh:
			if strings.EqualFold(emailData.ToAddress, targetEmail) && strings.Contains(emailData.SubjectData, "Security Code") {
				re := regexp.MustCompile(`security code is (\d{6})`)
				matches := re.FindStringSubmatch(emailData.BodyStr)
				if len(matches) == 2 {
					return matches[1], nil
				}
			}
		}
	}
}

// EmailReader reads emails from the IMAP server and sends them to the email channel
func EmailReader(client *client.Client, emailCh chan<- imap.EmailData) {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		emails, err := imap.GetMailRange(client, 10)
		if err != nil {
			return
		}
		for _, emailData := range emails {
			emailCh <- emailData
		}
	}
}
