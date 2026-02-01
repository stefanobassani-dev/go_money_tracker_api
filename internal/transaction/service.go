package transaction

import (
	"context"

	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type Service struct {
	transactionRepo domain.TransactionRepository
	tinkClient      domain.TinkClient
}

func NewService(transactionRepo domain.TransactionRepository, tinkClient domain.TinkClient) *Service {
	return &Service{
		transactionRepo: transactionRepo,
		tinkClient:      tinkClient,
	}
}

func (s *Service) SaveTransactions(ctx context.Context, userID string) error {
	err := s.tinkClient.FetchTransactions(ctx, userID, s.transactionRepo.UpsertBatch)
	if err != nil {
		return err
	}

	return nil
}
