package app

import (
	"github.com/claudiocleberson/bookstore_oauth-api/src/clients/cassandra"
	"github.com/claudiocleberson/bookstore_oauth-api/src/domain/access_token"
	"github.com/claudiocleberson/bookstore_oauth-api/src/http"
	"github.com/claudiocleberson/bookstore_oauth-api/src/repository/db"
	"github.com/claudiocleberson/bookstore_oauth-api/src/repository/rest_apis"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {

	//Initialize Database
	session, dbErr := cassandra.GetCassandraSession()
	if dbErr != nil {
		panic(dbErr)
	}
	//defer session.Close()

	//Initialize dependencies
	dbRepository := db.New(session)
	restRepository := rest_apis.NewRepository()
	atService := access_token.NewService(dbRepository, restRepository)

	atHandler := http.NewHandler(atService)
	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("oauth/access_token/", atHandler.Create)
	//Initialize router

	//Initialize application
	router.Run(":8081")

}
