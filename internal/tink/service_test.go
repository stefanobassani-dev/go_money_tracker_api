package tink

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

func Test_GetOrCreateTinkUser_UserExistsInDB_ReturnsLocalID(t *testing.T) {
	tinkService, userRepo, _, _ := setup(t)

	userID := "test_user_id"
	tinkID := "test_tink_id"
	ctx := context.Background()
	userRepo.EXPECT().GetTinkIDByUserID(ctx, userID).Return(tinkID, nil)

	resultID, err := tinkService.GetOrCreateTinkUser(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, tinkID, resultID)
}

func Test_GetOrCreateTinkUser_DBLookupFails_ReturnsError(t *testing.T) {
	tinkService, userRepo, _, _ := setup(t)

	userID := "test_user_id"
	ctx := context.Background()

	genericErr := errors.New("generic error")

	userRepo.EXPECT().GetTinkIDByUserID(ctx, userID).Return("", genericErr)

	_, err := tinkService.GetOrCreateTinkUser(ctx, userID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, genericErr)
}

func Test_GetOrCreateTinkUser_UserNotFoundInDB_CreateSuccess_ReturnsNewID(t *testing.T) {
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

func Test_GetOrCreateTinkUser_UserNotFoundInDB_CreateFailsConflict_ReturnsNewID(t *testing.T) {
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

func Test_GetOrCreateTinkUser_UserNotFoundInDB_CreateFailsConflict_RecoverFails_ReturnsError(t *testing.T) {
	tinkService, userRepo, _, tinkClient := setup(t)

	userID := "test_user_id"
	ctx := context.Background()

	genericErr := errors.New("generic error")

	userRepo.EXPECT().GetTinkIDByUserID(ctx, userID).Return("", nil)
	tinkClient.EXPECT().CreateUser(ctx, userID).Return("", domain.ErrUserAlreadyExists)
	tinkClient.EXPECT().GetUserByExternalID(ctx, userID).Return("", genericErr)

	_, err := tinkService.GetOrCreateTinkUser(ctx, userID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, genericErr)
}

func Test_GetOrCreateTinkUser_UserNotFoundInDB_CreateFailsGeneric_ReturnsError(t *testing.T) {
	tinkService, userRepo, _, tinkClient := setup(t)

	userID := "test_user_id"
	ctx := context.Background()

	genericErr := errors.New("generic error")

	userRepo.EXPECT().GetTinkIDByUserID(ctx, userID).Return("", nil)
	tinkClient.EXPECT().CreateUser(ctx, userID).Return("", genericErr)

	_, err := tinkService.GetOrCreateTinkUser(ctx, userID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, genericErr)
}

func Test_GetOrCreateTinkUser_UserNotFoundInDB_CreateSuccess_UpdateTinkIDFails_ReturnsNewID(t *testing.T) {
	tinkService, userRepo, _, tinkClient := setup(t)

	userID := "test_user_id"
	tinkID := "test_tink_id"
	ctx := context.Background()

	genericErr := errors.New("generic error")

	userRepo.EXPECT().GetTinkIDByUserID(ctx, userID).Return("", nil)
	tinkClient.EXPECT().CreateUser(ctx, userID).Return(tinkID, nil)
	userRepo.EXPECT().UpdateTinkID(ctx, userID, tinkID).Return(genericErr)

	newTinkID, err := tinkService.GetOrCreateTinkUser(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, newTinkID, tinkID)
}

//func Test_SaveCredential_(t *testing.T) {
//	tinkService, userRepo, credRepo, tinkClient := setup(t)
//
//	ctx := context.Background()
//	credID := "cred_id"
//	userID := "user_id"
//
//	tinkClient.EXPECT().GetUserCredential(ctx, credID, userID).Return()
//}
