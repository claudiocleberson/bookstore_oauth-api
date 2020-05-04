package http

import (
	"net/http"
	"strings"

	"github.com/claudiocleberson/bookstore_oauth-api/src/domain/access_token"
	"github.com/claudiocleberson/bookstore_utils-shared/utils/rest_err"
	"github.com/gin-gonic/gin"
)

type AccessTokenHandler interface {
	GetById(c *gin.Context)
	Create(c *gin.Context)
}

type accessTokenHandler struct {
	service access_token.Service
}

func NewHandler(service access_token.Service) AccessTokenHandler {
	return &accessTokenHandler{
		service: service,
	}
}

func (handler *accessTokenHandler) GetById(c *gin.Context) {

	accessTokeId := strings.TrimSpace(c.Param("access_token_id"))

	accessToken, err := handler.service.GetById(accessTokeId)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}
	c.JSON(http.StatusOK, accessToken)
}

func (handler *accessTokenHandler) Create(c *gin.Context) {

	var request access_token.AccessTokenRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := rest_err.NewBadRequestError("invalid json body")
		c.JSON(restErr.Code, restErr)
		return
	}

	token, err := handler.service.Create(request)

	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusCreated, token)

}
