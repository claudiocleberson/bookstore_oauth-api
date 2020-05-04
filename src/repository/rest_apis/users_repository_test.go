package rest_apis

import (
	"net/http"
	"os"
	"testing"

	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestLoginUserTimeOutFromApi(t *testing.T) {

	//Flush any mockup done previously
	rest.FlushMockups()

	//Create the mock
	rest.AddMockups(&rest.Mock{
		URL:          "",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@hotmail.com","password":"password"}`,
		RespBody:     "",
		RespHTTPCode: -1,
	})

	//Initialize test
	repo := restUserRepository{}
	user, err := repo.LoginUser("luiz@hotmail.com", "klsdjflkjl")

	//Validate
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Code)
	assert.EqualValues(t, "invalid restClient response when trying to login user", err.Message)
}

func TestLoginUserInvaliderrorInterface(t *testing.T) {
	//Flush any mockup done previously
	rest.FlushMockups()

	//Create the mock
	rest.AddMockups(&rest.Mock{
		URL:          "http://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@hotmail.com","password":"password"}`,
		RespBody:     `{"message":"invalid login credentials","code":"404","error":"not_found"}`,
		RespHTTPCode: http.StatusNotFound,
	})

	//Initialize test
	repo := restUserRepository{}
	user, err := repo.LoginUser("luiz@hotmail.com", "klsdjflkjl")

	//Validate
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Code)
	assert.EqualValues(t, "invalid error interface when trying to login user", err.Message)
}

func TestLoginUserInvalidLoginCredencials(t *testing.T) {
	//Flush any mockup done previously
	rest.FlushMockups()

	//Create the mock
	rest.AddMockups(&rest.Mock{
		URL:          "http://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@hotmail.com","password":"password"}`,
		RespBody:     `{"message":"invalid login credentials","code":404,"error":"not_found"}`,
		RespHTTPCode: http.StatusNotFound,
	})

	//Initialize test
	repo := restUserRepository{}
	user, err := repo.LoginUser("luiz@hotmail.com", "klsdjflkjl")

	//Validate
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Code)
	assert.EqualValues(t, "invalid login credentials", err.Message)
}

func TestLoginUserInvalidUserJsonResponse(t *testing.T) {
	//Flush any mockup done previously
	rest.FlushMockups()

	//Create the mock
	rest.AddMockups(&rest.Mock{
		URL:          "http://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@hotmail.com","password":"password"}`,
		RespBody:     `{"id":"asdkj","first_name":"claudio","last_name":"pereira","email":"claudio@hotmail.com","status":"active"}`,
		RespHTTPCode: http.StatusOK,
	})

	//Initialize test
	repo := restUserRepository{}
	user, err := repo.LoginUser("luiz@hotmail.com", "klsdjflkjl")

	//Validate
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Code)
	assert.EqualValues(t, "error when trying to unmarshal users response", err.Message)
}

func TestLoginUserNoError(t *testing.T) {
	//Flush any mockup done previously
	rest.FlushMockups()

	//Create the mock
	rest.AddMockups(&rest.Mock{
		URL:          "http://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@hotmail.com","password":"password"}`,
		RespBody:     `{"id":1,"first_name":"claudio","last_name":"pereira","email":"claudio@hotmail.com","status":"active"}`,
		RespHTTPCode: http.StatusOK,
	})

	//Initialize test
	repo := restUserRepository{}
	user, err := repo.LoginUser("luiz@hotmail.com", "klsdjflkjl")

	//Validate
	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.EqualValues(t, "claudio", user.FirstName)
	assert.EqualValues(t, "claudio@hotmail.com", user.Email)
}
