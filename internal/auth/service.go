package auth

import (
	"log"

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

func (s *Service) register() (string, error) {
	//get client access token
	clientToken, err := s.tinkClient.GetClientAccessToken()
	if err != nil {
		return "", err
	}

	//create user
	externalUserId := "test_id"
	userRes, err := s.tinkClient.CreateUser(externalUserId, "IT", "it_IT", clientToken)
	if err != nil {
		if tErr, ok := tink.FromError(err); ok && tErr.StatusCode == 409 {
			log.Printf("User with external_user_id %v already exists", externalUserId)
		} else {
			return "", err
		}
	}

	return s.tinkClient.GetAuthorizationURL(externalUserId, clientToken, userRes.UserID, "IT", "it_IT")
}

func (s *Service) deleteUser() {
	externalUserId := "test_id"

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
