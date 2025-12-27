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

func (s *Service) register() {
	//get client access token
	clientToken, err1 := s.tinkClient.GetClientAccessToken()
	tink.HandleError(err1)

	//create user
	externalUserId := "test_id"
	userRes, err2 := s.tinkClient.CreateUser(externalUserId, "IT", "it_IT", clientToken)
	tink.HandleError(err2)
	log.Println(userRes)

	code, err3 := s.tinkClient.AuthorizationGrantDelegate(externalUserId, clientToken)
	tink.HandleError(err3)
	log.Println(code)

	log.Println(s.tinkClient.BuildUrl(
		s.tinkClient.ClientId, "", "http://localhost:8080/auth/callback", code, "IT", "it_IT"),
	)

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
