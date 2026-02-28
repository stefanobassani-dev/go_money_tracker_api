package tink

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type providerConsentsResponse struct {
	ProviderConsents []providerConsent `json:"providerConsents"`
}

type providerConsent struct {
	CredentialsId     string   `json:"credentialsId"`
	ProviderName      string   `json:"providerName"`
	Status            string   `json:"status"`
	SessionExpiryDate int64    `json:"sessionExpiryDate"`
	SessionExtendable bool     `json:"sessionExtendable"`
	AccountIds        []string `json:"accountIds"`
	StatusUpdated     int64    `json:"statusUpdated"`
}

type transactionResponse struct {
	NextPageToken string        `json:"nextPageToken"`
	Transactions  []transaction `json:"transactions"`
}

type transaction struct {
	ID                  string                 `json:"id"`
	AccountID           string                 `json:"accountId"`
	Amount              amount                 `json:"amount"` // Supponendo Amount sia in un file common
	BookedDateTime      time.Time              `json:"bookedDateTime"`
	TransactionDateTime time.Time              `json:"transactionDateTime"`
	ValueDateTime       time.Time              `json:"valueDateTime"`
	Dates               transactionDates       `json:"dates"`
	Descriptions        descriptions           `json:"descriptions"`
	Identifiers         transactionIdentifiers `json:"identifiers"`
	Status              string                 `json:"status"`
	Reference           string                 `json:"reference"`
	ProviderMutability  string                 `json:"providerMutability"`
	Categories          categories             `json:"categories"`
	Counterparties      counterparties         `json:"counterparties"`
	MerchantInformation merchantInformation    `json:"merchantInformation"`
	Types               transactionTypes       `json:"types"`
}

type transactionDates struct {
	Booked      string `json:"booked"`
	Transaction string `json:"transaction"`
	Value       string `json:"value"`
}

type descriptions struct {
	Display  string   `json:"display"`
	Original string   `json:"original"`
	Detailed detailed `json:"detailed"`
}

type detailed struct {
	Unstructured string `json:"unstructured"`
}

type transactionIdentifiers struct {
	ProviderTransactionID string `json:"providerTransactionId"`
}

type counterparties struct {
	Payee party `json:"payee"`
	Payer party `json:"payer"`
}

type party struct {
	Name        string           `json:"name"`
	Identifiers partyIdentifiers `json:"identifiers"`
}

type partyIdentifiers struct {
	FinancialInstitution financialInstitution `json:"financialInstitution"` // Supponendo financialInstitution sia common
}

type merchantInformation struct {
	MerchantCategoryCode string `json:"merchantCategoryCode"`
	MerchantName         string `json:"merchantName"`
}

type categories struct {
	Pfm pfmCategory `json:"pfm"`
}

type pfmCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type transactionTypes struct {
	Type                         string `json:"type"`
	FinancialInstitutionTypeCode string `json:"financialInstitutionTypeCode"`
}

// FetchTransactions retrieves all transactions for a specific user from Tink.
// It handles pagination automatically and uses a callback to save batches of transactions.
func (c *Client) FetchTransactions(ctx context.Context, externalUserID string,
	saveTransactions func(context.Context, []domain.Transaction) error) error {

	scopes := []string{"transactions:read"}
	accessToken, err := c.ExchangeUserToken(ctx, externalUserID, scopes)
	if err != nil {
		return err
	}

	var pageToken string
	for {
		u, _ := url.Parse(string(PathListTransactions))
		q := u.Query()
		q.Set("pageSize", "100")
		if pageToken != "" {
			q.Set("pageToken", pageToken)
		}
		u.RawQuery = q.Encode()

		var res transactionResponse

		props := Props{
			ctx:         ctx,
			httpClient:  c.httpClient,
			baseUrl:     c.Cfg.BaseUrl,
			path:        u.String(),
			method:      http.MethodGet,
			contentType: ContentTypeJSON,
			body:        nil,
			result:      &res,
			token:       accessToken,
		}

		if err := call(props); err != nil {
			return err
		}

		transactions := ToDomainTransactionList(res.Transactions)
		if err := saveTransactions(ctx, transactions); err != nil {
			return err
		}

		if res.NextPageToken == "" {
			break
		}
		pageToken = res.NextPageToken
	}

	return nil
}
