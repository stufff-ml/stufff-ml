package backend

import (
	"strings"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"

	"github.com/gin-gonic/gin"
	"github.com/majordomusio/commons/pkg/errors"
	"github.com/majordomusio/commons/pkg/gae/logger"
	"github.com/majordomusio/commons/pkg/util"
)

// AuthenticateAndAuthorize authenicates and authorizes a client based on its token
func AuthenticateAndAuthorize(ctx context.Context, scope, token string) (string, error) {

	auth, err := GetAuthorization(ctx, token)
	if err != nil {
		logger.Error(ctx, "backend.auth.authenticate", err.Error())
		return "", errors.New("Invalid Token")
	}

	// check if the token has been revoked or is expired
	if auth.Revoked {
		return "", errors.New("Token has been revoked")
	}

	if auth.Expires > 0 {
		if auth.Expires < util.Timestamp() {
			return "", errors.New("Token has expired")
		}
	}

	// check the authorization
	if strings.Contains(auth.Scope, ScopeAdmin) {
		return auth.ClientID, nil
	}

	if strings.Contains(auth.Scope, scope) {
		return auth.ClientID, nil
	}

	return "", errors.New("Not authorized")
}

// ClientIDFromToken retrieves the client id based on the access token
func ClientIDFromToken(ctx context.Context, token string) string {

	auth, err := GetAuthorization(ctx, token)
	if err != nil {
		logger.Error(ctx, "backend.auth.client_id", err.Error())
		return ""
	}

	// all good
	return auth.ClientID
}

// GetAuthorization returns the authorization for the given token
func GetAuthorization(ctx context.Context, token string) (*AuthorizationDS, error) {
	var auth = AuthorizationDS{}

	key := "auth.token." + strings.ToLower(token)
	_, err := memcache.Gob.Get(ctx, key, &auth)

	if err != nil {
		err = datastore.Get(ctx, AuthorizationKey(ctx, strings.ToLower(token)), &auth)
		if err == nil {
			cache := memcache.Item{}
			cache.Key = key
			cache.Object = auth
			cache.Expiration, _ = time.ParseDuration(DefaultCacheDuration)
			memcache.Gob.Set(ctx, &cache)
		} else {
			logger.Error(ctx, "backend.auth.get", err.Error())
			return nil, err
		}
	}

	return &auth, err
}

// GetToken extracts the bearer token
func GetToken(ctx context.Context, c *gin.Context) string {

	auth := c.Request.Header["Authorization"]
	if len(auth) == 0 {
		return ""
	}

	parts := strings.Split(auth[0], " ")
	if len(parts) != 2 {
		return ""
	}

	return strings.ToLower(parts[1])
}

// CreateClientAndAuthentication creates a new client and its authentication
func CreateClientAndAuthentication(ctx context.Context, clientID, clientSecret, scope, token string) error {

	// ClientResource
	cr := ClientResourceDS{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Created:      util.Timestamp(),
	}

	key := ClientResourceKey(ctx, cr.ClientID)
	_, err := datastore.Put(ctx, key, &cr)
	if err != nil {
		logger.Error(ctx, "backend.auth.create", err.Error())
		return err
	}

	// Authorization
	auth := AuthorizationDS{
		ClientID: cr.ClientID,
		Scope:    scope,
		Token:    token,
		Revoked:  false,
		Expires:  0,
		Created:  util.Timestamp(),
	}

	key = AuthorizationKey(ctx, auth.Token)
	_, err = datastore.Put(ctx, key, &auth)
	if err != nil {
		logger.Error(ctx, "backend.auth.create", err.Error())
		return err
	}

	return nil
}
