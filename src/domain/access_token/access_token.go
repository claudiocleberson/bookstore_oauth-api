package access_token

import (
	"strings"
	"time"

	"github.com/claudiocleberson/bookstore_users-api/utils/crypto_utils"
	"github.com/claudiocleberson/bookstore_users-api/utils/rest_err"
)

const (
	expirationTime              = 24
	grantTypePassword           = "password"
	grantTypeClientCrendentials = "client_credentials"
)

type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`
	//Use when for grant type password
	Username string `json:"username"`
	Password string `json:"password"`
	//Used when for grantType client_secret
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`

	Scope string `json:"scope"`
}

func (r *AccessTokenRequest) Validate() *rest_err.RestErr {

	switch r.GrantType {
	case grantTypeClientCrendentials:
		if r.ClientId == "" || r.ClientSecret == "" {
			return rest_err.NewBadRequestError("ClientId and ClienteSecret must be informed for grant_type=client_credentials")
		}
	case grantTypePassword:
		if r.Username == "" || r.Password == "" {
			return rest_err.NewBadRequestError("Username and password must be informed for grant_type=password")
		}
	default:
		return rest_err.NewBadRequestError("that grant_type is not supported")

	}

	return nil
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
	ClientId    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

func (token *AccessToken) Validate() *rest_err.RestErr {

	token.AccessToken = strings.TrimSpace(token.AccessToken)
	if token.AccessToken == "" {
		return rest_err.NewBadRequestError("invalid access token")
	}

	if token.UserId <= 0 {
		return rest_err.NewBadRequestError("invalid token user_id")
	}

	// if token.ClientId <= 0 {
	// 	return rest_err.NewBadRequestError("invalid token client_id")
	// }
	if token.Expires <= 0 {
		return rest_err.NewBadRequestError("invalid token expiration time")
	}

	return nil
}

func GetNewAccessToken(userId int64) AccessToken {
	return AccessToken{
		UserId:  userId,
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at *AccessToken) IsExpired() bool {
	now := time.Now().UTC()
	expirationTime := time.Unix(at.Expires, 0)
	return expirationTime.Before(now)
}

func (at *AccessToken) Generate() {
	userToken := string(at.UserId) + string(at.Expires)
	at.AccessToken = crypto_utils.GetMd5(userToken)
}
