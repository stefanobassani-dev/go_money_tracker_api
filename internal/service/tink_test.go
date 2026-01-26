package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stefanobassani-dev/money-tracker/internal/domain"
	"github.com/stefanobassani-dev/money-tracker/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*TinkService, *mocks.MockUserRepository, *mocks.MockCredentialRepository,
	*mocks.MockTinkClient) {
	userRepo := mocks.NewMockUserRepository(t)
	credentialRepo := mocks.NewMockCredentialRepository(t)
	tinkClient := mocks.NewMockTinkClient(t)

	tinkService := NewTinkService(tinkClient, userRepo, credentialRepo)

	return tinkService, userRepo, credentialRepo, tinkClient
}

func Test_UserHasTinkID(t *testing.T) {
	tinkService, userRepo, _, _ := setup(t)

	userID := "test_user_id"
	tinkID := "test_tink_id"
	ctx := context.Background()
	userRepo.EXPECT().GetTinkIDByUserID(ctx, userID).Return(tinkID, nil)

	resultID, err := tinkService.GetOrCreateTinkUser(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, tinkID, resultID)
}

func Test_UserHasNoTinkID_Success(t *testing.T) {
	tinkService, userRepo, _, tinkClient := setup(t)

	userID := "test_user_id"
	tinkID := "test_tink_id"
	ctx := context.Background()

	userRepo.EXPECT().GetTinkIDByUserID(ctx, userID).Return("", nil)
	tinkClient.EXPECT().CreateUser(ctx, userID).Return(tinkID, nil)
	userRepo.EXPECT().UpdateTinkID(ctx, userID, tinkID).Return(nil)

	resTinkID, err := tinkService.GetOrCreateTinkUser(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, tinkID, resTinkID)
}

func Test_UserHasNoTinkID_DBError(t *testing.T) {
	tinkService, userRepo, _, _ := setup(t)

	userID := "test_user_id"
	ctx := context.Background()

	genericErr := errors.New("Generic error")

	userRepo.EXPECT().GetTinkIDByUserID(ctx, userID).Return("", genericErr)

	_, err := tinkService.GetOrCreateTinkUser(ctx, userID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, genericErr)
}

func Test_UserHasNoTinkID_UserConflict(t *testing.T) {
	tinkService, userRepo, _, tinkClient := setup(t)

	userID := "test_user_id"
	tinkID := "test_tink_id"
	ctx := context.Background()

	userRepo.EXPECT().GetTinkIDByUserID(ctx, userID).Return("", nil)
	tinkClient.EXPECT().CreateUser(ctx, userID).Return("", domain.ErrUserAlreadyExists)
	tinkClient.EXPECT().GetUserByExternalID(ctx, userID).Return(tinkID, nil)
	userRepo.EXPECT().UpdateTinkID(ctx, userID, tinkID).Return(nil)

	resTinkID, err := tinkService.GetOrCreateTinkUser(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, tinkID, resTinkID)
}
