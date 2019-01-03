package app

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/stufff-ml/stufff-ml/pkg/api"
	"github.com/stufff-ml/stufff-ml/pkg/helper"

	"github.com/stufff-ml/stufff-ml/internal/backend"
	"github.com/stufff-ml/stufff-ml/internal/cloud"
)

// InitEndpoint creates an initial set of records to get started
func InitEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	token := backend.GetToken(ctx, c)
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

	err := cloud.CreateClientAndAuthentication(ctx, clientID, clientSecret, "admin", t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error"})
		return
	}

	resp := api.ClientResource{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Token:        t,
	}
	c.JSON(http.StatusOK, &resp)

}
