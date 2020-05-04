package db

import (
	"github.com/claudiocleberson/bookstore_oauth-api/src/domain/access_token"
	"github.com/claudiocleberson/bookstore_users-api/utils/rest_err"
	"github.com/gocql/gocql"
)

const (
	queryGetAccessToken       = "SELECT access_token, user_id,client_id,expires FROM access_token where access_token=?;"
	queryInsertAccessToken    = "INSERT INTO access_token(access_token,user_id,client_id,expires) VALUES (?,?,?,?);"
	queryUpdateExpirationTime = "UPDATE access_token SET expires=? WHERE access_token=?;"
)

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, *rest_err.RestErr)
	Create(access_token.AccessToken) (*access_token.AccessToken, *rest_err.RestErr)
	UpdateExpirationTime(access_token.AccessToken) *rest_err.RestErr
}

func New(dbsession *gocql.Session) DbRepository {
	return &dbRepository{
		dbSession: dbsession,
	}
}

type dbRepository struct {
	dbSession *gocql.Session
}

func (r *dbRepository) Create(token access_token.AccessToken) (*access_token.AccessToken, *rest_err.RestErr) {

	if err := r.dbSession.Query(queryInsertAccessToken,
		token.AccessToken,
		token.UserId,
		token.ClientId, token.Expires).Exec(); err != nil {
		return nil, rest_err.NewInternalServerError(err.Error())
	}

	return &token, nil
}

func (r *dbRepository) UpdateExpirationTime(token access_token.AccessToken) *rest_err.RestErr {

	if err := r.dbSession.Query(queryUpdateExpirationTime,
		token.Expires,
		token.AccessToken,
	).Exec(); err != nil {
		return rest_err.NewInternalServerError(err.Error())
	}

	return nil
}

func (r *dbRepository) GetById(tokenId string) (*access_token.AccessToken, *rest_err.RestErr) {

	var resultAt access_token.AccessToken
	err := r.dbSession.Query(queryGetAccessToken, tokenId).Scan(
		&resultAt.AccessToken,
		&resultAt.UserId,
		&resultAt.ClientId,
		&resultAt.Expires)

	if err != nil {
		if err == gocql.ErrNotFound {
			return nil, rest_err.NewNotFoundError("no token for the given id")
		}
		return nil, rest_err.NewInternalServerError(err.Error())
	}

	return &resultAt, nil
}
