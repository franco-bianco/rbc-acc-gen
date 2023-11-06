package gen

type SubmitOTPRes struct {
	Verified              bool   `json:"verified"`
	OnetimePasswordStatus string `json:"onetimePasswordStatus"`
	RemainingAttempts     int    `json:"remainingAttempts"`
}

type ResStatus struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type PublicKeyRes struct {
	PublicKey string `json:"publicKey"`
}

type RobTokenRes struct {
	ID                string `json:"id"`
	RootID            any    `json:"rootId"`
	RequestParameters struct {
		AccessToken  string `json:"access_token"`
		Idp          string `json:"idp"`
		RequestType  string `json:"request_type"`
		Scope        string `json:"scope"`
		ResponseType string `json:"response_type"`
		RedirectURI  string `json:"redirect_uri"`
		Locale       string `json:"locale"`
		ClientID     string `json:"client_id"`
		UILocales    string `json:"ui_locales"`
	} `json:"requestParameters"`
	RequestHeaders struct {
	} `json:"requestHeaders"`
	IdpConfiguration struct {
		ID                     string `json:"id"`
		Urn                    string `json:"urn"`
		URL                    string `json:"url"`
		SessionTimeout         string `json:"sessionTimeout"`
		SessionExpiration      int    `json:"sessionExpiration"`
		FormPostRedirect       bool   `json:"formPostRedirect"`
		UserinfoEndpoint       string `json:"userinfoEndpoint"`
		UserinfoDefaultKeyName any    `json:"userinfoDefaultKeyName"`
		AttributesToSurrogate  []any  `json:"attributesToSurrogate"`
		SurrogatedAttributes   []any  `json:"surrogatedAttributes"`
	} `json:"idpConfiguration"`
	ConsentProviderConfiguration any `json:"consentProviderConfiguration"`
	Client                       struct {
		Disabled   bool   `json:"disabled"`
		ClientID   string `json:"clientId"`
		ClientName []struct {
			Value  string `json:"value"`
			Locale string `json:"locale"`
		} `json:"clientName"`
		ClientUrn             string   `json:"clientUrn"`
		GrantTypes            []string `json:"grantTypes"`
		PkceRequired          bool     `json:"pkceRequired"`
		SignedRequestRequired bool     `json:"signedRequestRequired"`
		Scopes                []string `json:"scopes"`
		RedirectUris          []string `json:"redirectUris"`
		IdentityProviders     []struct {
			ID             string `json:"id"`
			Attributes     any    `json:"attributes"`
			SelectorScript any    `json:"selectorScript"`
		} `json:"identityProviders"`
		AccessTokenExpiry                  string `json:"accessTokenExpiry"`
		IDTokenExpiry                      string `json:"idTokenExpiry"`
		ProspectTokenExpiry                string `json:"prospectTokenExpiry"`
		TokenSigningPolicy                 string `json:"tokenSigningPolicy"`
		SsoTimeout                         string `json:"ssoTimeout"`
		SrfRequired                        bool   `json:"srfRequired"`
		AudClaimRequired                   bool   `json:"audClaimRequired"`
		SsoEnabled                         bool   `json:"ssoEnabled"`
		EncSubRequired                     bool   `json:"encSubRequired"`
		SignInRestEnabled                  bool   `json:"signInRestEnabled"`
		ConsentEnabled                     bool   `json:"consentEnabled"`
		ConsentUserEligibilityCheckEnabled bool   `json:"consentUserEligibilityCheckEnabled"`
	} `json:"client"`
	AggregatorApp           any    `json:"aggregatorApp"`
	Authentication          any    `json:"authentication"`
	AdapterDropoff          any    `json:"adapterDropoff"`
	UserConsent             any    `json:"userConsent"`
	State                   string `json:"state"`
	SignInRestConfiguration any    `json:"signInRestConfiguration"`
}

type CompleteRegistrationRes struct {
	Result      string `json:"result"`
	RedirectURL string `json:"redirectUrl"`
}
