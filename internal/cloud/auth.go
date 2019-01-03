package cloud

import (
	"strings"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"

	"github.com/majordomusio/commons/pkg/gae/logger"
	"github.com/majordomusio/commons/pkg/util"

	"github.com/stufff-ml/stufff-ml/internal/types"
)

// GetAuthorization returns the authorization for the given token
func GetAuthorization(ctx context.Context, token string) (*types.AuthorizationDS, error) {
	var auth = types.AuthorizationDS{}

	key := "auth.token." + strings.ToLower(token)
	_, err := memcache.Gob.Get(ctx, key, &auth)

	if err != nil {
		err = datastore.Get(ctx, AuthorizationKey(ctx, strings.ToLower(token)), &auth)
		if err == nil {
			cache := memcache.Item{}
			cache.Key = key
			cache.Object = auth
			cache.Expiration, _ = time.ParseDuration(types.DefaultCacheDuration)
			memcache.Gob.Set(ctx, &cache)
		} else {
			logger.Error(ctx, "backend.auth.get", err.Error())
			return nil, err
		}
	}

	return &auth, err
}

// CreateClientAndAuthentication creates a new client and its authentication
func CreateClientAndAuthentication(ctx context.Context, clientID, clientSecret, scope, token string) error {

	// ClientResource
	cr := types.ClientResourceDS{
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
	auth := types.AuthorizationDS{
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
