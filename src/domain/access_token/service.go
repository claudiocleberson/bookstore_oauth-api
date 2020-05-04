package access_token

import (
	"strings"

	"github.com/claudiocleberson/bookstore_oauth-api/src/repository/rest_apis"
	"github.com/claudiocleberson/bookstore_utils-shared/utils/rest_err"
)

type Service interface {
	GetById(string) (*AccessToken, *rest_err.RestErr)
	Create(AccessTokenRequest) (*AccessToken, *rest_err.RestErr)
	UpdateExpirationTime(AccessToken) *rest_err.RestErr
}

type DbRepository interface {
	GetById(string) (*AccessToken, *rest_err.RestErr)
	Create(AccessToken) (*AccessToken, *rest_err.RestErr)
	UpdateExpirationTime(AccessToken) *rest_err.RestErr
}

type service struct {
	dbRepository   DbRepository
	restRepository rest_apis.RestUsersRepository
}

func NewService(dbRepo DbRepository, restRepo rest_apis.RestUsersRepository) Service {
	return &service{
		dbRepository:   dbRepo,
		restRepository: restRepo,
	}
}

func (s *service) GetById(tokenId string) (*AccessToken, *rest_err.RestErr) {

	accessTokenId := strings.TrimSpace(tokenId)
	if len(accessTokenId) == 0 {
		return nil, rest_err.NewBadRequestError("invalid access token")
	}
	accessToken, err := s.dbRepository.GetById(tokenId)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (s *service) Create(request AccessTokenRequest) (*AccessToken, *rest_err.RestErr) {

	if err := request.Validate(); err != nil {
		return nil, err
	}

	//Todo - Support both client_credentials and password

	//Authenticate the user against the Users API:
	user, err := s.restRepository.LoginUser(request.Username, request.Password)

	if err != nil {
		return nil, err
	}

	//Genereate a new access token
	at := GetNewAccessToken(user.Id)
	at.Generate()

	if err := at.Validate(); err != nil {
		return nil, err
	}
	return s.dbRepository.Create(at)

}

func (s *service) UpdateExpirationTime(token AccessToken) *rest_err.RestErr {
	if err := token.Validate(); err != nil {
		return err
	}
	return s.dbRepository.UpdateExpirationTime(token)
}
