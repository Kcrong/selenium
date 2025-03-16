package fedcm

// LoginState represents the state of a login
type LoginState string

const (
	// SignIn represents a sign-in state
	SignIn LoginState = "SignIn"
	// SignUp represents a sign-up state
	SignUp LoginState = "SignUp"
)

// Account represents an account displayed in a FedCM account list
// See: https://w3c-fedid.github.io/FedCM/#dictdef-identityprovideraccount
//
//	https://w3c-fedid.github.io/FedCM/#webdriver-accountlist
type Account struct {
	accountData map[string]string
}

// NewAccount creates a new Account instance
func NewAccount(accountData map[string]string) *Account {
	return &Account{
		accountData: accountData,
	}
}

// GetAccountID returns the account ID
func (a *Account) GetAccountID() string {
	return a.accountData["accountId"]
}

// GetEmail returns the email address
func (a *Account) GetEmail() string {
	return a.accountData["email"]
}

// GetName returns the account name
func (a *Account) GetName() string {
	return a.accountData["name"]
}

// GetGivenName returns the given name
func (a *Account) GetGivenName() string {
	return a.accountData["givenName"]
}

// GetPictureURL returns the picture URL
func (a *Account) GetPictureURL() string {
	return a.accountData["pictureUrl"]
}

// GetIDPConfigURL returns the IDP config URL
func (a *Account) GetIDPConfigURL() string {
	return a.accountData["idpConfigUrl"]
}

// GetTermsOfServiceURL returns the terms of service URL
func (a *Account) GetTermsOfServiceURL() string {
	return a.accountData["termsOfServiceUrl"]
}

// GetPrivacyPolicyURL returns the privacy policy URL
func (a *Account) GetPrivacyPolicyURL() string {
	return a.accountData["privacyPolicyUrl"]
}

// GetLoginState returns the login state
func (a *Account) GetLoginState() LoginState {
	return LoginState(a.accountData["loginState"])
}
