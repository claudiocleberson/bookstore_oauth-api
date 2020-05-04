package rest_apis

import (
	"encoding/json"
	"time"

	"github.com/claudiocleberson/bookstore_oauth-api/src/domain/users"
	"github.com/claudiocleberson/bookstore_utils-shared/utils/rest_err"
	"github.com/mercadolibre/golang-restclient/rest"
)

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, *rest_err.RestErr)
}

func NewRepository() RestUsersRepository {
	return &restUserRepository{}
}

var (
	usersRestClient = rest.RequestBuilder{
		BaseURL: "http://localhost:8080",
		Timeout: 10000 * time.Millisecond,
	}
)

type restUserRepository struct {
}

func (r *restUserRepository) LoginUser(email string, password string) (*users.User, *rest_err.RestErr) {

	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}

	response := usersRestClient.Post("/users/login", request)
	if response == nil || response.Response == nil {
		return nil, rest_err.NewInternalServerError("invalid restClient response when trying to login user")
	}

	if response.StatusCode > 299 {
		var restErr rest_err.RestErr
		err := json.Unmarshal(response.Bytes(), &restErr)
		if err != nil {
			return nil, rest_err.NewNotFoundError("invalid error interface when trying to login user")
		}
		return nil, &restErr
	}

	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, rest_err.NewInternalServerError("error when trying to unmarshal users response")
	}

	return &user, nil
}
