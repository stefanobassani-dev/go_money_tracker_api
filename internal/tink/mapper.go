package tink

import (
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
