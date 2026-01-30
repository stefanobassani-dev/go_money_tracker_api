package tink

import (
	"math"
	"strconv"
	"time"

	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

func ToDomainCredential(dto *Credential, externalUserID string) domain.Credential {
	return domain.Credential{
		CredentialID: dto.ID,
		TinkUserID:   dto.UserID,
		UserID:       externalUserID,
		ProviderName: dto.ProviderName,
		Status:       dto.Status,
		ExpiresAt:    time.Unix(0, dto.SessionExpiryDate*int64(time.Millisecond)),
		LastUpdated:  time.Unix(0, dto.Updated*int64(time.Millisecond)),
		Type:         dto.Type,
	}
}

func ToDomainAccountList(dtos []Account, userID string) []domain.Account {
	result := make([]domain.Account, len(dtos))
	for i, dto := range dtos {
		result[i] = ToDomainAccount(dto, userID)
	}
	return result
}

func ToDomainAccount(dto Account, userID string) domain.Account {
	balance := parseBalance(dto.Balances.Booked.Amount.Value)

	return domain.Account{
		ID:                     dto.ID,
		UserID:                 userID,
		CredentialID:           nil,
		Name:                   dto.Name,
		Type:                   domain.AccountType(dto.Type),
		Balance:                balance,
		Currency:               dto.Balances.Booked.Amount.CurrencyCode,
		FinancialInstitutionID: dto.FinancialInstitutionID,
		LastRefreshed:          dto.Dates.LastRefreshed,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}
}

func parseBalance(val Value) float64 {
	unscaled, err := strconv.ParseInt(val.UnscaledValue, 10, 64)
	if err != nil {
		return 0
	}

	scale, err := strconv.Atoi(val.Scale)
	if err != nil {
		return 0
	}

	return float64(unscaled) * math.Pow(10, -float64(scale))
}

func ToDomainProviderConsentList(dtos []ProviderConsent) []domain.ProviderConsent {
	result := make([]domain.ProviderConsent, len(dtos))
	for i, dto := range dtos {
		result[i] = ToDomainProviderConsent(dto)
	}
	return result
}

func ToDomainProviderConsent(dto ProviderConsent) domain.ProviderConsent {
	return domain.ProviderConsent{
		CredentialsID:     dto.CredentialsId,
		ProviderName:      dto.ProviderName,
		Status:            dto.Status,
		SessionExpiryDate: time.Unix(0, dto.SessionExpiryDate*int64(time.Millisecond)),
		SessionExtendable: dto.SessionExtendable,
		AccountIDs:        dto.AccountIds,
		StatusUpdated:     time.Unix(0, dto.StatusUpdated*int64(time.Millisecond)),
	}
}
