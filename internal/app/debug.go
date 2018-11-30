package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/majordomusio/commons/pkg/gae/logger"
	"github.com/majordomusio/commons/pkg/util"

	"github.com/stufff-ml/stufff-ml/internal/backend"
)

// DebugEndpoint is for testing only
func DebugEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// MigrateEndpoint is for testing only
func MigrateEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// SeedEndpoint creates an initial set of records to get started
func SeedEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	// ClientResource
	cr := backend.ClientResource{
		ClientID:     "aaaa",
		ClientSecret: "aaaa",
		Created:      util.Timestamp(),
	}

	key := backend.ClientResourceKey(ctx, cr.ClientID)
	_, err := datastore.Put(ctx, key, &cr)
	if err != nil {
		logger.Error(ctx, "api.seed", err.Error())
	}

	// Authorization
	auth := backend.Authorization{
		ClientID: cr.ClientID,
		Token:    "xoxo-ffffffff",
		Revoked:  false,
		Expires:  0,
		Created:  util.Timestamp(),
	}

	key = backend.AuthorizationKey(ctx, auth.Token)
	_, err = datastore.Put(ctx, key, &auth)
	if err != nil {
		logger.Error(ctx, "api.seed", err.Error())
	}

	m := backend.Model{
		ClientID: cr.ClientID,
		Domain:   "buy",
		Revision: 1,
		Event:    "buy",
		Created:  util.Timestamp(),
	}
	key = backend.ModelKey(ctx, m.ClientID, m.Domain)
	_, err = datastore.Put(ctx, key, &m)
	if err != nil {
		logger.Error(ctx, "api.seed", err.Error())
	}

	// all good
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
