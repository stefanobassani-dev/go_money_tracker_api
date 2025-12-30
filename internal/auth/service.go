package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"

	tinkapi2 "github.com/stefanobassani-dev/money-tracker/internal/integration/tinkapi"
)

type Service struct {
	repo         *Repository
	tinkClient   *tinkapi2.Client
	tokenManager *tinkapi2.TokenManager
}

func NewService(repo *Repository, tinkClient *tinkapi2.Client, tokenManager *tinkapi2.TokenManager) *Service {
	return &Service{
		repo:         repo,
		tinkClient:   tinkClient,
		tokenManager: tokenManager,
	}
}

func (s *Service) register(ctx context.Context) (string, error) {
	clientToken, err := s.tokenManager.GetToken()
	if err != nil {
		return "", err
	}

	externalUserId := "2d7b9b46-94fe-435e-aa7d-95ac55fc188d"

	err = s.ensureTinkUser(ctx, externalUserId, clientToken)
	if err != nil {
		return "", err
	}

	return s.tinkClient.GetAuthorizationURL(externalUserId, clientToken, "IT", "it_IT")
}

func (s *Service) ensureTinkUser(ctx context.Context, extID, token string) error {
	res, err := s.tinkClient.CreateUser(extID, "IT", "it_IT", token)
	if err == nil {
		_ = s.repo.LinkTinkUser(ctx, extID, res.UserID)
		return nil
	}

	if tErr, ok := tinkapi2.FromError(err); ok && tErr.StatusCode == 409 {
		return s.recoverExistingUser(ctx, extID, token)
	}

	return err
}

func (s *Service) recoverExistingUser(ctx context.Context, externalUserID, clientToken string) error {
	tinkUserID, err := s.repo.GetTinkUserID(ctx, externalUserID)
	if err == nil && tinkUserID != "" {
		return nil
	}

	log.Printf("User %s exists on Tink but not in DB, starting API recovery...", externalUserID)

	code, err := s.tinkClient.AuthorizationGrant(externalUserID, clientToken)
	if err != nil {
		return fmt.Errorf("recovery: failed authorization grant: %w", err)
	}

	tokenResponse, err := s.tinkClient.GetUserAccessToken(code)
	if err != nil {
		return fmt.Errorf("recovery: failed to get user access token: %w", err)
	}

	userResponse, err := s.tinkClient.GetUserDetails(tokenResponse.AccessToken)
	if err != nil {
		return fmt.Errorf("recovery: failed to fetch user details: %w", err)
	}

	err = s.repo.LinkTinkUser(ctx, externalUserID, userResponse.ID)
	if err != nil {
		// Logghiamo l'errore ma non blocchiamo l'utente, abbiamo comunque l'ID
		log.Printf("Warning: recovered user from Tink but failed to update local DB: %v", err)
	}

	return nil
}

func (s *Service) deleteUser() {
	externalUserId := "2d7b9b46-94fe-435e-aa7d-95ac55fc188d"

	clientToken, _, err1 := s.tinkClient.GetClientAccessToken()
	tinkapi2.HandleError(err1)
	code, err4 := s.tinkClient.AuthorizationGrant(externalUserId, clientToken)
	tinkapi2.HandleError(err4)

	tokenResponse, err2 := s.tinkClient.GetUserAccessToken(code)
	tinkapi2.HandleError(err2)
	log.Println(s.tinkClient.GetUserDetails(tokenResponse.AccessToken))
	err := s.tinkClient.DeleteUser(tokenResponse.AccessToken)
	if err != nil {
		log.Println(err)
	}
}

func (s *Service) ProcessCallback(r *http.Request) error {
	credID := r.URL.Query().Get("credentials_id")
	if credID == "" {
		credID = r.URL.Query().Get("credentialsId")
	}

	if credID == "" {
		return nil
	}
	//TODO controllare hmac
	userID := r.URL.Query().Get("state")

	return s.repo.SaveTinkCredential(r.Context(), userID, credID)
}
