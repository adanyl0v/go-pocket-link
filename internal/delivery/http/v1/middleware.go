package v1

import (
	"fmt"
	"github.com/adanyl0v/go-pocket-link/pkg/auth/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	contextUserID = "user_id"

	headerAuthorization = "Authorization"
)

func (h *Handler) useAuth(c *gin.Context) {
	headerValue, err := parseAuthHeader(c)
	if err != nil {
		writeAbort(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	accessTokenClaims, err := h.services.Tokens.ValidateAccessToken(headerValue)
	if err != nil {
		writeAbort(c, http.StatusUnauthorized, "invalid access token claims", err)
		return
	}

	rawUserID, ok := accessTokenClaims[jwt.ClaimsSubject].(string)
	if !ok {
		writeAbort(c, http.StatusUnauthorized, "missed user id", nil)
		return
	}

	c.Set(contextUserID, rawUserID)
}

func parseAuthHeader(c *gin.Context) (string, error) {
	header := c.GetHeader(headerAuthorization)
	if header == "" {
		return "", fmt.Errorf("no %s header", headerAuthorization)
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid %s header", headerAuthorization)
	}

	return parts[1], nil
}
