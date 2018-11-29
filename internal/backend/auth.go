package backend

import (
	"strings"
	"time"

	"github.com/ratchetcc/commons/pkg/util"
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
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
