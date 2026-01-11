package service

import "github.com/stefanobassani-dev/money-tracker/internal/domain"

type TinkService struct {
	t domain.TinkClient
}

func NewTinkService(t domain.TinkClient) *TinkService {
	return &TinkService{
		t: t,
	}
}
