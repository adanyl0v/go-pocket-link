package v1

import (
	"github.com/gin-gonic/gin"
	"go-pocket-link/internal/domain"
	"go-pocket-link/internal/service"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const (
	logError = "error"

	cookieRefreshToken = "refresh_token"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitEndpoints(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/ping", h.ping)
	routerGroup.POST("/sign-up", h.signUp)
	routerGroup.POST("/sign-in", h.signIn)
}

func (h *Handler) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
	slog.Debug("pinged connection")
}

func (h *Handler) signUp(c *gin.Context) {
	var input struct {
		Name     string `json:"name" form:"name" binding:"required"`
		Email    string `json:"email" form:"email" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}
	slog.Debug("sign up input", "input", input)

	if err := c.ShouldBind(&input); err != nil {
		writeError(c, http.StatusBadRequest, "invalid input", err)
		return
	}

	user := domain.User{Name: input.Name, Email: input.Email, Password: input.Password}
	if err := h.services.Users.Save(c, &user); err != nil {
		writeError(c, http.StatusInternalServerError, "saving user", err)
		return
	}

	refreshToken, err := h.services.Auth.NewRefreshToken()
	if err != nil {
		writeError(c, http.StatusInternalServerError, "creating refresh token", err)
		return
	}
	slog.Debug("created refresh token", "token", refreshToken)

	session := domain.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
	}
	if err = h.services.Sessions.Save(c, &session); err != nil {
		writeError(c, http.StatusInternalServerError, "saving session", err)
		return
	}
	slog.Debug("saved session", "session", session)

	accessToken, err := h.services.Auth.NewAccessToken(session.ID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "creating access token", err)
		return
	}
	slog.Debug("created access token", "token", accessToken)

	setRefreshTokenCookie(c, refreshToken, h.services.Auth.RefreshTokenTTL())

	c.JSON(http.StatusCreated, accessToken)
	slog.Debug("signed up", "user", user, "jwt", struct {
		AccessToken, RefreshToken string
	}{AccessToken: accessToken, RefreshToken: refreshToken})
}

func (h *Handler) signIn(c *gin.Context) {
	var input struct {
		Email    string `json:"email" form:"email" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}
	slog.Debug("sign in input", "input", input)

	if err := c.ShouldBind(&input); err != nil {
		writeError(c, http.StatusBadRequest, "invalid input", err)
		return
	}

	user, err := h.services.Users.GetByCredentials(c, input.Email, input.Password)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "fetching user", err)
		return
	}

	refreshToken, err := h.services.Auth.NewRefreshToken()
	if err != nil {
		writeError(c, http.StatusInternalServerError, "creating refresh token", err)
		return
	}
	slog.Debug("created refresh token", "token", refreshToken)

	session := domain.Session{
		RefreshToken: refreshToken,
	}
	if err = h.services.Sessions.Update(c, &session); err != nil {
		writeError(c, http.StatusInternalServerError, "updating session", err)
		return
	}
	slog.Debug("updated session")

	accessToken, err := h.services.Auth.NewAccessToken(session.ID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "creating access token", err)
		return
	}
	slog.Debug("created access token", "token", accessToken)

	setRefreshTokenCookie(c, refreshToken, h.services.Auth.RefreshTokenTTL())

	c.JSON(http.StatusOK, accessToken)
	slog.Debug("signed in", "user", user, "jwt", struct {
		AccessToken, RefreshToken string
	}{AccessToken: accessToken, RefreshToken: refreshToken})
}

func writeError(c *gin.Context, status int, message string, err error) {
	c.JSON(status, gin.H{logError: message})
	slog.Error(message, logError, err)
}

func setRefreshTokenCookie(c *gin.Context, token string, ttl time.Duration) {
	hostDomain := strings.Split(c.Request.Host, ":")[0]
	c.SetCookie(cookieRefreshToken, token, int(ttl.Seconds()), "/", hostDomain, false, true)
}
