package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-pocket-link/internal/domain"
	"go-pocket-link/internal/service"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const (
	logError = "error"

	cookiesRefreshToken = "refresh_token"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitEndpoints(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/ping", h.handlePing)
	routerGroup.POST("/sign-up", h.handleSignUp)
	routerGroup.POST("/sign-in", h.handleSignIn)

	protectedGroup := routerGroup.Group("/", h.useAuth)
	{
		protectedGroup.POST("/log-out", h.handleLogOut)
	}
}

func (h *Handler) handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
	slog.Debug("pinged connection")
}

func (h *Handler) handleSignUp(c *gin.Context) {
	var input struct {
		Name     string `json:"name" form:"name" binding:"required"`
		Email    string `json:"email" form:"email" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}
	if err := bindInput(c, &input); err != nil {
		return
	}

	user := domain.User{Name: input.Name, Email: input.Email, Password: input.Password}
	if err := h.services.Users.Save(c, &user); err != nil {
		writeError(c, http.StatusInternalServerError, "saving user", err)
		return
	}

	refreshToken, err := newRefreshToken(c, h.services.Auth)
	if err != nil {
		return
	}

	session := domain.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
	}
	if err = h.services.Sessions.Save(c, &session); err != nil {
		writeError(c, http.StatusInternalServerError, "saving session", err)
		return
	}
	slog.Debug("saved session", "session", session)

	accessToken, err := newAccessToken(c, h.services.Auth, user.ID)
	if err != nil {
		return
	}

	setRefreshTokenCookies(c, refreshToken, h.services.Auth.RefreshTokenTTL())

	c.JSON(http.StatusCreated, accessToken)
	slog.Debug("signed up", "user", user, "jwt", struct {
		AccessToken, RefreshToken string
	}{AccessToken: accessToken, RefreshToken: refreshToken})
}

func (h *Handler) handleSignIn(c *gin.Context) {
	var input struct {
		Email    string `json:"email" form:"email" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}
	if err := bindInput(c, &input); err != nil {
		return
	}

	user, err := h.services.Users.GetByCredentials(c, input.Email, input.Password)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "getting user", err)
		return
	}

	refreshToken, err := newRefreshToken(c, h.services.Auth)
	if err != nil {
		return
	}

	session := domain.Session{
		RefreshToken: refreshToken,
	}
	if err = h.services.Sessions.Update(c, &session); err != nil {
		writeError(c, http.StatusInternalServerError, "updating session", err)
		return
	}
	slog.Debug("updated session")

	accessToken, err := newAccessToken(c, h.services.Auth, user.ID)
	if err != nil {
		return
	}

	setRefreshTokenCookies(c, refreshToken, h.services.Auth.RefreshTokenTTL())

	c.JSON(http.StatusOK, accessToken)
	slog.Debug("signed in", "user", user, "jwt", struct {
		AccessToken, RefreshToken string
	}{AccessToken: accessToken, RefreshToken: refreshToken})
}

func (h *Handler) handleLogOut(c *gin.Context) {
	_, sessionID, err := userAndSessionIDsFromContext(c)
	if err != nil {
		return
	}
	slog.Debug("got session", "id", sessionID)

	if err = h.services.Sessions.Invalidate(c, sessionID); err != nil {
		writeError(c, http.StatusInternalServerError, "invalidating session", err)
		return
	}
	slog.Debug("invalidated session", "id", sessionID)

	c.Redirect(http.StatusTemporaryRedirect, "/sign-in")
}

func bindInput(c *gin.Context, input any) error {
	if err := c.ShouldBind(input); err != nil {
		writeError(c, http.StatusBadRequest, "invalid input", err)
		return err
	}
	slog.Debug("bound input", "input", input)
	return nil
}

func newAccessToken(c *gin.Context, auth *service.AuthService, userID uuid.UUID) (string, error) {
	token, err := auth.NewAccessToken(userID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "creating access token", err)
		return "", err
	}
	slog.Debug("created access token", "token", token)
	return token, nil
}

func newRefreshToken(c *gin.Context, auth *service.AuthService) (string, error) {
	token, err := auth.NewRefreshToken()
	if err != nil {
		writeError(c, http.StatusInternalServerError, "creating refresh token", err)
		return "", err
	}
	slog.Debug("created refresh token", "token", token)
	return token, nil
}

func setRefreshTokenCookies(c *gin.Context, token string, ttl time.Duration) {
	hostDomain := strings.Split(c.Request.Host, ":")[0]
	c.SetCookie(cookiesRefreshToken, token, int(ttl.Seconds()), "/", hostDomain, false, true)
}

func userAndSessionIDsFromContext(c *gin.Context) (uuid.UUID, uuid.UUID, error) {
	userID, err := idFromContext(c, contextUserID)
	if err != nil {
		writeError(c, http.StatusBadRequest, "getting user id from context", err)
		return uuid.Nil, uuid.Nil, err
	}

	sessionID, err := idFromContext(c, contextSessionID)
	if err != nil {
		writeError(c, http.StatusBadRequest, "getting session id from context", err)
		return uuid.Nil, uuid.Nil, err
	}

	return userID, sessionID, nil
}
