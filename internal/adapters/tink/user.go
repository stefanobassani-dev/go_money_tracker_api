package tink

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type createUserRequest struct {
	ExternalUserID string `json:"external_user_id,omitempty"`
	Market         string `json:"market"`
	Locale         string `json:"locale,omitempty"`
	RetentionClass string `json:"retention_class,omitempty"`
}

type createUserResponse struct {
	UserID         string `json:"user_id"`
	ExternalUserID string `json:"external_user_id"`
}

func (c *Client) CreateUser(ctx context.Context, externalID string) (string, error) {
	clientToken, err := c.TokenManager.GetToken(ctx)
	if err != nil {
		return "", err
	}

	req := &createUserRequest{
		ExternalUserID: externalID,
		Market:         c.Cfg.Market,
		Locale:         c.Cfg.Locale,
		RetentionClass: "permanent",
	}

	body, err := toJSONReader(req)
	if err != nil {
		return "", err
	}

	var res createUserResponse

	props := Props{
		ctx:         ctx,
		httpClient:  c.httpClient,
		baseUrl:     c.Cfg.BaseUrl,
		path:        string(PathUserCreate),
		method:      http.MethodPost,
		contentType: ContentTypeJSON,
		body:        body,
		result:      &res,
		token:       clientToken,
	}
	err = call(props)
	if err != nil {
		var tErr *TinkError
		if errors.As(err, &tErr) && tErr.StatusCode == http.StatusConflict {
			return "", domain.ErrUserAlreadyExists
		}
		slog.Error("tink error while creating a new tink user", "userID", externalID)
		return "", domain.ErrInternal
	}

	return res.UserID, nil
}

type user struct {
	AppID          string   `json:"appId"`
	Created        string   `json:"created"`
	ExternalUserID string   `json:"externalUserId"`
	Flags          []string `json:"flags"`
	ID             string   `json:"id"`
	NationalID     string   `json:"nationalId"`
	Profile        profile  `json:"profile"`
	Username       string   `json:"username"`
}

type profile struct {
	Currency             string               `json:"currency"`
	Locale               string               `json:"locale"`
	Market               string               `json:"market"`
	NotificationSettings notificationSettings `json:"notificationSettings"`
	PeriodAdjustedDay    int                  `json:"periodAdjustedDay"`
	PeriodMode           string               `json:"periodMode"`
	TimeZone             string               `json:"timeZone"`
}

type notificationSettings struct {
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

func (c *Client) GetUserByExternalID(ctx context.Context, externalID string) (string, error) {
	scopes := []string{"accounts:read"}
	accessToken, err := c.ExchangeUserToken(ctx, externalID, scopes)
	if err != nil {
		return "", err
	}

	var res user
	props := Props{
		ctx:         ctx,
		httpClient:  c.httpClient,
		baseUrl:     c.Cfg.BaseUrl,
		path:        string(PathGetUser),
		method:      http.MethodGet,
		contentType: ContentTypeJSON,
		body:        nil,
		result:      &res,
		token:       accessToken,
	}
	err = call(props)
	if err != nil {
		var tErr *TinkError
		if errors.As(err, &tErr) && tErr.StatusCode == http.StatusNotFound {
			return "", domain.ErrUserNotFound
		}
		slog.Error("tink error while retrieving tink user", "userID", externalID)
		return "", domain.ErrInternal
	}

	return res.ID, nil
}
