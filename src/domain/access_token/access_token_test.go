package access_token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetNewAccessToken(t *testing.T) {
	at := GetNewAccessToken(1)

	//Go default implementation
	if at.IsExpired() {
		t.Error("brand new access token should not be expired")
	}

	if at.AccessToken != "" {
		t.Error("new access token shoud not have defined access token id")
	}

	if at.UserId != 0 {
		t.Error("new access token should not have an associated user id")
	}

	//Testify implemantation
	assert.False(t, at.IsExpired(), "brand new access token should not be expired")
	assert.Empty(t, at.AccessToken, "new access token should not have defined access token id")
	assert.Equal(t, int64(0), at.UserId, "new access token should not have an associated user id")
}

func TestAccessTokenIsExpired(t *testing.T) {

	at := AccessToken{}

	//Default go testing implementation
	if !at.IsExpired() {
		t.Error("empty access token should be expired by default")
	}

	//Testify implemantation
	at.Expires = time.Now().UTC().Add(3 * time.Hour).Unix()
	assert.False(t, at.IsExpired(), "empty access token should be expired by default")

}
