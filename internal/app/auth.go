package app

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthorizeRequest looks for a token and returns the matching app id
func AuthorizeRequest(c *gin.Context) (string, bool) {

	auth := c.Request.Header["Authorization"]
	if len(auth) == 0 {
		return "", false
	}

	parts := strings.Split(auth[0], " ")
	if len(parts) != 2 {
		return "", false
	}

	return strings.ToLower(parts[1]), true
}
