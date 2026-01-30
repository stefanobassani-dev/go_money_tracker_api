package account

import "github.com/stefanobassani-dev/money-tracker/internal/domain"

type AccountService struct {
	accountRepo domain.AccountRepository
}

func NewAccountService(accountRepo domain.AccountRepository) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
	}
}
