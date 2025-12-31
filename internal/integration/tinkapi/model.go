package tinkapi

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

type TinkUserResponse struct {
	AppID          string      `json:"appId"`
	ID             string      `json:"id"`
	Created        int64       `json:"created"`
	ExternalUserID string      `json:"externalUserId"`
	NationalID     interface{} `json:"nationalId"`
	Flags          []string    `json:"flags"`
	Profile        UserProfile `json:"profile"`
	Username       interface{} `json:"username"`
}

type UserProfile struct {
	Currency             string               `json:"currency"`
	Locale               string               `json:"locale"`
	Market               string               `json:"market"`
	NotificationSettings NotificationSettings `json:"notificationSettings"`
	PeriodAdjustedDay    int                  `json:"periodAdjustedDay"`
	PeriodMode           string               `json:"periodMode"`
	TimeZone             string               `json:"timeZone"`
	PeriodSettings       PeriodSettings       `json:"periodSettings"`
}

type NotificationSettings struct {
	Balance         bool `json:"balance"`
	Budget          bool `json:"budget"`
	DoubleCharge    bool `json:"doubleCharge"`
	Income          bool `json:"income"`
	LargeExpense    bool `json:"largeExpense"`
	SummaryMonthly  bool `json:"summaryMonthly"`
	SummaryWeekly   bool `json:"summaryWeekly"`
	Transaction     bool `json:"transaction"`
	UnusualCategory bool `json:"unusualCategory"`
	UnusualAccount  bool `json:"unusualAccount"`
	Einvoices       bool `json:"einvoices"`
	Fraud           bool `json:"fraud"`
	LeftToSpend     bool `json:"leftToSpend"`
	LoanUpdate      bool `json:"loanUpdate"`
}

type PeriodSettings struct {
	Mode              string `json:"mode"`
	AdjustedPeriodDay int    `json:"adjustedPeriodDay"`
}

type CreateUserRequest struct {
	ExternalUserID string `json:"external_user_id,omitempty"`
	Market         string `json:"market"`
	Locale         string `json:"locale,omitempty"`
	RetentionClass string `json:"retention_class,omitempty"` // "permanent" o "temporary"
}

type CreateUserResponse struct {
	UserID         string `json:"user_id"`
	ExternalUserID string `json:"external_user_id"`
}

type Credential struct {
	ID                string            `json:"id"`
	ProviderName      string            `json:"providerName"`
	Type              string            `json:"type"`
	Status            string            `json:"status"`
	StatusUpdated     int64             `json:"statusUpdated"`
	StatusPayload     string            `json:"statusPayload"`
	Updated           int64             `json:"updated"`
	Fields            map[string]string `json:"fields"`
	SessionExpiryDate int64             `json:"sessionExpiryDate"`
	UserID            string            `json:"userId"`
}

type CredentialResponse struct {
	Credential []Credential `json:"credentials"`
}
