package backend

import (
	"strings"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"

	"github.com/majordomusio/commons/pkg/gae/logger"
	"github.com/majordomusio/commons/pkg/util"
)

// ClientIDFromToken retrieves the client id based on the access token
func ClientIDFromToken(ctx context.Context, token string) (string, bool) {
	var auth = Authorization{}

	key := "token.client_id." + strings.ToLower(token)
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
			return "", false
		}
	}

	// check if the toen has been revoked or is expired
	if auth.Revoked {
		return "", false
	}

	if auth.Expires > 0 {
		if auth.Expires < util.Timestamp() {
			return "", false
		}
	}

	// all good
	return auth.ClientID, true
}

// CreateClientAndAuthentication creates a new client and its authentication
func CreateClientAndAuthentication(ctx context.Context, clientID, clientSecret, token string) error {

	// ClientResource
	cr := ClientResource{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Created:      util.Timestamp(),
	}

	key := ClientResourceKey(ctx, cr.ClientID)
	_, err := datastore.Put(ctx, key, &cr)
	if err != nil {
		logger.Error(ctx, "backend.client.create", err.Error())
		return err
	}

	// Authorization
	auth := Authorization{
		ClientID: cr.ClientID,
		Token:    token,
		Revoked:  false,
		Expires:  0,
		Created:  util.Timestamp(),
	}

	key = AuthorizationKey(ctx, auth.Token)
	_, err = datastore.Put(ctx, key, &auth)
	if err != nil {
		logger.Error(ctx, "backend.client.create", err.Error())
		return err
	}

	return nil
}
