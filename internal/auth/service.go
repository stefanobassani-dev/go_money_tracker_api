package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/stefanobassani-dev/money-tracker/internal/tink"
)

type Service struct {
	repo       *Repository
	tinkClient *tink.Client
}

func NewService(repo *Repository, tinkClient *tink.Client) *Service {
	return &Service{
		repo:       repo,
		tinkClient: tinkClient,
	}
}

func (s *Service) register(ctx context.Context) (string, error) {
	clientToken, err := s.tinkClient.GetClientAccessToken()
	if err != nil {
		return "", err
	}

	externalUserId := "2d7b9b46-94fe-435e-aa7d-95ac55fc188d"

	tinkUserId, err := s.ensureTinkUser(ctx, externalUserId, clientToken)
	if err != nil {
		return "", err
	}

	return s.tinkClient.GetAuthorizationURL(externalUserId, clientToken, tinkUserId, "IT", "it_IT")
}

func (s *Service) ensureTinkUser(ctx context.Context, extID, token string) (string, error) {
	res, err := s.tinkClient.CreateUser(extID, "IT", "it_IT", token)
	if err == nil {
		_ = s.repo.LinkTinkUser(ctx, extID, res.UserID)
		return res.UserID, nil
	}

	if tErr, ok := tink.FromError(err); ok && tErr.StatusCode == 409 {
		return s.recoverExistingUser(ctx, extID, token)
	}

	return "", err // Qualsiasi altro errore blocca tutto
}

func (s *Service) recoverExistingUser(ctx context.Context, externalUserID, clientToken string) (string, error) {
	// 1. Tentativo prioritario: Database locale
	tinkUserID, err := s.repo.GetTinkUserID(ctx, externalUserID)
	if err == nil && tinkUserID != "" {
		return tinkUserID, nil
	}

	// 2. Se non è nel DB (disallineamento), lo recuperiamo da Tink
	log.Printf("User %s exists on Tink but not in DB, starting API recovery...", externalUserID)

	// 2a. Ottieni il codice di autorizzazione per l'utente esistente
	code, err := s.tinkClient.AuthorizationGrant(externalUserID, clientToken)
	if err != nil {
		return "", fmt.Errorf("recovery: failed authorization grant: %w", err)
	}

	// 2b. Scambia il codice con un User Access Token
	tokenResponse, err := s.tinkClient.GetUserAccessToken(code)
	if err != nil {
		return "", fmt.Errorf("recovery: failed to get user access token: %w", err)
	}

	// 2c. Ottieni i dettagli per estrarre il tink_user_id (id interno di Tink)
	userResponse, err := s.tinkClient.GetUserDetails(tokenResponse.AccessToken)
	if err != nil {
		return "", fmt.Errorf("recovery: failed to fetch user details: %w", err)
	}

	// 3. Ora che l'abbiamo recuperato, "riariamo" il DB per evitare questo giro la prossima volta
	err = s.repo.LinkTinkUser(ctx, externalUserID, userResponse.ID)
	if err != nil {
		// Logghiamo l'errore ma non blocchiamo l'utente, abbiamo comunque l'ID
		log.Printf("Warning: recovered user from Tink but failed to update local DB: %v", err)
	}

	return userResponse.ID, nil
}

func (s *Service) deleteUser() {
	externalUserId := "2d7b9b46-94fe-435e-aa7d-95ac55fc188d"

	clientToken, err1 := s.tinkClient.GetClientAccessToken()
	tink.HandleError(err1)
	code, err4 := s.tinkClient.AuthorizationGrant(externalUserId, clientToken)
	tink.HandleError(err4)

	tokenResponse, err2 := s.tinkClient.GetUserAccessToken(code)
	tink.HandleError(err2)
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

	//TODO salvare il credential id sulla riga dello user
	return s.repo.SaveTinkCredential(r.Context(), userID, credID)
}
