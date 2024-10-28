package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-pocket-link/internal/service"
	"net/http"
	"strings"
)

const (
	contextUserID    = "user_id"
	contextSessionID = "session_id"

	headerAuthorization = "Authorization"
)

func (h *Handler) useAuth(c *gin.Context) {
	accessToken, err := parseAuthorizationHeader(c)
	if err != nil {
		writeAbort(c, http.StatusUnauthorized, fmt.Sprintf("parsing %s header", headerAuthorization), err)
		return
	}

	rawUserID, err := parseAccessToken(h.services.Auth, accessToken)
	if err != nil {
		writeAbort(c, http.StatusUnauthorized, "parsing claims", err)
		return
	}

	c.Set(contextUserID, rawUserID)

	rawSessionID, err := rawSessionIDByUserID(c, h.services.Sessions, rawUserID)
	if err != nil {
		writeAbort(c, http.StatusUnauthorized, "getting session from user id", err)
		return
	}

	c.Set(contextSessionID, rawSessionID)
}

func parseAuthorizationHeader(c *gin.Context) (string, error) {
	header := c.GetHeader(headerAuthorization)
	if header == "" {
		return "", fmt.Errorf("empty %s header", headerAuthorization)
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", fmt.Errorf("invalid %s header", headerAuthorization)
	}

	return headerParts[1], nil
}

func parseAccessToken(auth *service.AuthService, accessToken string) (string, error) {
	claims, err := auth.ParseAccessToken(accessToken)
	if err != nil {
		return "", err
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return "", err
	}

	return sub, nil
}

func rawSessionIDByUserID(c *gin.Context, sessions *service.SessionsService, rawUserID string) (string, error) {
	userID, err := uuid.Parse(rawUserID)
	if err != nil {
		return "", err
	}

	session, err := sessions.GetByUserID(c, userID)
	if err != nil {
		return "", err
	}

	return session.ID.String(), nil
}

func idFromContext(c *gin.Context, key string) (uuid.UUID, error) {
	id, err := uuid.Parse(c.GetString(key))
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
