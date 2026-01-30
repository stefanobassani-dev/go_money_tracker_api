package tink

import "time"

type AuthorizationRequest struct {
	ActorClientID string `url:"actor_client_id"`
	ExternalID    string `url:"external_user_id"`
	IDHint        string `url:"id_hint,omitempty"`
	Scope         string `url:"scope"`
}

type AuthorizationResponse struct {
	Code string `json:"code"`
}

type TokenRequest struct {
	GrantType    string `url:"grant_type"`
	ClientID     string `url:"client_id"`
	ClientSecret string `url:"client_secret"`
	Code         string `url:"code,omitempty"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	IdHint      string `json:"id_hint"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type Credential struct {
	ID                      string            `json:"id"`
	ProviderName            string            `json:"providerName"`
	Type                    string            `json:"type"`
	Status                  string            `json:"status"`
	StatusPayload           string            `json:"statusPayload"`
	StatusUpdated           int64             `json:"statusUpdated"`
	Updated                 int64             `json:"updated"`
	SessionExpiryDate       int64             `json:"sessionExpiryDate"`
	UserID                  string            `json:"userId"`
	Fields                  map[string]string `json:"fields"`
	SupplementalInformation interface{}       `json:"supplementalInformation"`
}

type CreateUserRequest struct {
	ExternalUserID string `json:"external_user_id,omitempty"`
	Market         string `json:"market"`
	Locale         string `json:"locale,omitempty"`
	RetentionClass string `json:"retention_class,omitempty"`
}

type CreateUserResponse struct {
	UserID         string `json:"user_id"`
	ExternalUserID string `json:"external_user_id"`
}

type User struct {
	AppID          string   `json:"appId"`
	Created        string   `json:"created"`
	ExternalUserID string   `json:"externalUserId"`
	Flags          []string `json:"flags"`
	ID             string   `json:"id"`
	NationalID     string   `json:"nationalId"`
	Profile        Profile  `json:"profile"`
	Username       string   `json:"username"`
}

type Profile struct {
	Currency             string               `json:"currency"`
	Locale               string               `json:"locale"`
	Market               string               `json:"market"`
	NotificationSettings NotificationSettings `json:"notificationSettings"`
	PeriodAdjustedDay    int                  `json:"periodAdjustedDay"`
	PeriodMode           string               `json:"periodMode"`
	TimeZone             string               `json:"timeZone"`
}

type NotificationSettings struct {
	Balance         bool `json:"balance"`
	Budget          bool `json:"budget"`
	DoubleCharge    bool `json:"doubleCharge"`
	EInvoices       bool `json:"einvoices"`
	Fraud           bool `json:"fraud"`
	Income          bool `json:"income"`
	LargeExpense    bool `json:"largeExpense"`
	LeftToSpend     bool `json:"leftToSpend"`
	LoanUpdate      bool `json:"loanUpdate"`
	SummaryMonthly  bool `json:"summaryMonthly"`
	SummaryWeekly   bool `json:"summaryWeekly"`
	Transaction     bool `json:"transaction"`
	UnusualAccount  bool `json:"unusualAccount"`
	UnusualCategory bool `json:"unusualCategory"`
}

type AccountResponse struct {
	Accounts      []Account `json:"accounts"`
	NextPageToken string    `json:"nextPageToken"`
}

type Account struct {
	ID                     string      `json:"id"`
	Name                   string      `json:"name"`
	Type                   string      `json:"type"`
	Balances               Balances    `json:"balances"`
	Identifiers            Identifiers `json:"identifiers"`
	Dates                  Dates       `json:"dates"`
	FinancialInstitutionID string      `json:"financialInstitutionId"`
	CustomerSegment        string      `json:"customerSegment"`
}

type Balances struct {
	Booked    BalanceDetails `json:"booked"`
	Available BalanceDetails `json:"available"`
}

type BalanceDetails struct {
	Amount Amount `json:"amount"`
}

type Amount struct {
	Value        Value  `json:"value"`
	CurrencyCode string `json:"currencyCode"`
}

type Value struct {
	UnscaledValue string `json:"unscaledValue"`
	Scale         string `json:"scale"`
}

type Identifiers struct {
	Iban                 Iban                 `json:"iban"`
	FinancialInstitution FinancialInstitution `json:"financialInstitution"`
}

type Iban struct {
	Iban string `json:"iban"`
	Bban string `json:"bban"`
}

type FinancialInstitution struct {
	AccountNumber    string                 `json:"accountNumber"`
	ReferenceNumbers map[string]interface{} `json:"referenceNumbers"`
}

type Dates struct {
	LastRefreshed time.Time `json:"lastRefreshed"`
}

type ProviderConsentsResponse struct {
	ProviderConsents []ProviderConsent `json:"providerConsents"`
}

type ProviderConsent struct {
	CredentialsId     string   `json:"credentialsId"`
	ProviderName      string   `json:"providerName"`
	Status            string   `json:"status"`
	SessionExpiryDate int64    `json:"sessionExpiryDate"`
	SessionExtendable bool     `json:"sessionExtendable"`
	AccountIds        []string `json:"accountIds"`
	StatusUpdated     int64    `json:"statusUpdated"`
}
