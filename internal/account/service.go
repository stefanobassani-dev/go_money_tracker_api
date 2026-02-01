package account

import (
	"context"
	"slices"

	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type Service struct {
	accountRepo domain.AccountRepository
	tinkClient  domain.TinkClient
}

func NewService(accountRepo domain.AccountRepository, tinkClient domain.TinkClient) *Service {
	return &Service{
		accountRepo: accountRepo,
		tinkClient:  tinkClient,
	}
}

func (s *Service) SaveAccountsByCredentialID(ctx context.Context, credentialID, userID string) error {
	consentList, err := s.tinkClient.ProviderConsent(ctx, userID)
	if err != nil {
		return err
	}

	var matching domain.ProviderConsent
	for _, consent := range consentList {
		if consent.CredentialsID == credentialID {
			matching = consent
			break
		}
	}
	accountIDList := matching.AccountIDs

	accountList, err := s.tinkClient.ListAccounts(ctx, userID)
	if err != nil {
		return err
	}

	for _, account := range accountList {
		if slices.Contains(accountIDList, account.ID) {
			err := s.accountRepo.CreateCredential(ctx, account, credentialID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
