package tink

import (
	"context"
	"net/http"

	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

func (c *Client) FetchTransactions(ctx context.Context, externalUserID string,
	saveTransactions func(context.Context, []domain.Transaction) error) error {
	scopes := []string{"transactions:read"}
	accessToken, err := c.ExchangeUserToken(ctx, externalUserID, scopes)
	if err != nil {
		return err
	}

	var pageToken string
	for {
		path := string(PathListTransactions) + "?pageSize=100"
		if pageToken != "" {
			path += "&pageToken=" + pageToken
		}

		var res = struct {
			Transactions  []Transaction
			NextPageToken string
		}{}
		props := Props{
			ctx:         ctx,
			httpClient:  c.httpClient,
			baseUrl:     c.Cfg.BaseUrl,
			path:        path,
			method:      http.MethodGet,
			contentType: ContentTypeJSON,
			body:        nil,
			result:      &res,
			token:       accessToken,
		}

		err = call(props)
		if err != nil {
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
