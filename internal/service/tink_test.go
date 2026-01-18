package service

import (
	"testing"

	"github.com/stefanobassani-dev/money-tracker/internal/mocks"
)

func Test_GetConnectURL(t *testing.T) {
	mockTink := mocks.NewMockTinkClient(t)

	mockTink.EXPECT().GetUserCredential()
}
