package api

import (
	"strings"

	"golang.org/x/net/context"

	"github.com/gin-gonic/gin"
	"github.com/majordomusio/commons/pkg/errors"
	"github.com/majordomusio/commons/pkg/gae/logger"
	"github.com/majordomusio/commons/pkg/util"

	"github.com/stufff-ml/stufff-ml/internal/backend"
	"github.com/stufff-ml/stufff-ml/internal/types"
)

// AuthenticateAndAuthorize authenicates and authorizes a client based on its token
func AuthenticateAndAuthorize(ctx context.Context, scope, token string) (string, error) {

	auth, err := backend.GetAuthorization(ctx, token)
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
	if strings.Contains(auth.Scope, types.ScopeAdmin) {
		return auth.ClientID, nil
	}

	if strings.Contains(auth.Scope, scope) {
		return auth.ClientID, nil
	}

	return "", errors.New("Not authorized")
}

// ClientIDFromToken retrieves the client id based on the access token
func ClientIDFromToken(ctx context.Context, token string) string {

	auth, err := backend.GetAuthorization(ctx, token)
	if err != nil {
		logger.Error(ctx, "backend.auth.client_id", err.Error())
		return ""
	}

	// all good
	return auth.ClientID
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
