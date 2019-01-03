package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	a "github.com/stufff-ml/stufff-ml/pkg/api"
	"github.com/stufff-ml/stufff-ml/pkg/helper"

	"github.com/stufff-ml/stufff-ml/internal/backend"
)

// InitEndpoint creates an initial set of records to get started
func InitEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	token := GetToken(ctx, c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	if token != os.Getenv("ADMIN_CLIENT_TOKEN") {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	clientID, _ := helper.ShortUUID()
	clientSecret, _ := helper.SimpleUUID()
	t, _ := helper.RandomToken()

	err := backend.CreateClientAndAuthentication(ctx, clientID, clientSecret, "admin", t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error"})
		return
	}

	resp := a.ClientResource{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Token:        t,
	}
	c.JSON(http.StatusOK, &resp)

}
