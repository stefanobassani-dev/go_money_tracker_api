package tink

import (
	"context"
	"net/http"
	"time"

	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type accountResponse struct {
	Accounts      []account `json:"accounts"`
	NextPageToken string    `json:"nextPageToken"`
}

type account struct {
	ID                     string      `json:"id"`
	Name                   string      `json:"name"`
	Type                   string      `json:"type"`
	Balances               balances    `json:"balances"`
	Identifiers            identifiers `json:"identifiers"`
	Dates                  dates       `json:"dates"`
	FinancialInstitutionID string      `json:"financialInstitutionId"`
	CustomerSegment        string      `json:"customerSegment"`
	CredentialID           string      `json:"credentialId"`
}

type balances struct {
	Booked    balanceDetails `json:"booked"`
	Available balanceDetails `json:"available"`
}

type balanceDetails struct {
	Amount amount `json:"amount"`
}

type value struct {
	UnscaledValue string `json:"unscaledValue"`
	Scale         string `json:"scale"`
}

type identifiers struct {
	Iban                 iban                 `json:"iban"`
	FinancialInstitution financialInstitution `json:"financialInstitution"`
}

type iban struct {
	Iban string `json:"iban"`
	Bban string `json:"bban"`
}

type dates struct {
	LastRefreshed time.Time `json:"lastRefreshed"`
}

func (c *Client) ListAccounts(ctx context.Context, externalUserID string) ([]domain.Account, error) {
	scopes := []string{"accounts:read"}
	accessToken, err := c.ExchangeUserToken(ctx, externalUserID, scopes)
	if err != nil {
		return []domain.Account{}, err
	}

	var res = struct {
		Accounts []account
	}{}
	props := Props{
		ctx:         ctx,
		httpClient:  c.httpClient,
		baseUrl:     c.Cfg.BaseUrl,
		path:        string(PathListAccounts),
		method:      http.MethodGet,
		contentType: ContentTypeJSON,
		body:        nil,
		result:      &res,
		token:       accessToken,
	}

	err = call(props)
	if err != nil {
		return []domain.Account{}, err
	}

	return ToDomainAccountList(res.Accounts, externalUserID), nil
}
