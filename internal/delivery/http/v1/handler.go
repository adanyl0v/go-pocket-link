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
	PublicEndpointSignIn = "/sign-in"
	PublicEndpointSignUp = "/sign-up"

	ApiEndpointPing   = "/ping"
	ApiEndpointSignIn = PublicEndpointSignIn
	ApiEndpointSignUp = PublicEndpointSignUp
	ApiEndpointLogOut = "/log-out"
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
	routerGroup.GET(ApiEndpointPing, h.handlePing)
	routerGroup.POST(ApiEndpointSignUp, h.handleSignUp)
	routerGroup.POST(ApiEndpointSignIn, h.handleSignIn)

	protectedGroup := routerGroup.Group("/", h.useAuth)
	{
		protectedGroup.POST(ApiEndpointLogOut, h.handleLogOut)
	}
}

func (h *Handler) handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
	slog.Debug("connection is healthy")
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
	slog.Debug("bound input", "input", input)

	user := domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}
	if err := h.services.Users.Save(c, &user); err != nil {
		writeError(c, http.StatusInternalServerError, "saving user", err)
		return
	}

	tokens, err := h.services.Tokens.NewTokenPair(user.ID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "creating token pair", err)
		return
	}

	if err = h.services.Tokens.SaveRefreshTokenFromString(c, tokens.RefreshToken); err != nil {
		writeError(c, http.StatusInternalServerError, "saving refresh token", err)
		return
	}
	slog.Debug("saved refresh token", "token", tokens.RefreshToken)

	setRefreshTokenCookies(c, tokens.RefreshToken, h.services.Tokens.RefreshTokenTTL)
	slog.Debug("set refresh token cookies")

	c.JSON(http.StatusCreated, tokens.AccessToken)
	slog.Debug("signed up", "user", user, "jwt", tokens)
}

func (h *Handler) handleSignIn(c *gin.Context) {
	var input struct {
		Email    string `json:"email" form:"email" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}
	if err := bindInput(c, &input); err != nil {
		return
	}
	slog.Debug("bound input", "input", input)

	user, err := h.services.Users.GetByCredentials(c, input.Email, input.Password)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "getting user", err)
		return
	}

	tokens, err := h.services.Tokens.NewTokenPair(user.ID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "creating token pair", err)
		return
	}

	if err = h.services.Tokens.SaveRefreshTokenFromString(c, tokens.RefreshToken); err != nil {
		writeError(c, http.StatusInternalServerError, "saving refresh token", err)
		return
	}
	slog.Debug("saved refresh token", "token", tokens.RefreshToken)

	setRefreshTokenCookies(c, tokens.RefreshToken, h.services.Tokens.RefreshTokenTTL)
	slog.Debug("set refresh token cookies")

	c.JSON(http.StatusOK, tokens.AccessToken)
	slog.Debug("signed in", "user", user, "jwt", tokens)
}

func (h *Handler) handleLogOut(c *gin.Context) {
	rawUserID, exists := c.Get(contextUserID)
	if !exists {
		writeError(c, http.StatusUnauthorized, "no user id in context", nil)
		return
	}

	userID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		writeError(c, http.StatusUnauthorized, "parsing user id", err)
		return
	}
	slog.Debug("got user id from context", "id", userID)

	if err = h.services.Tokens.InvalidateUser(c, userID); err != nil {
		writeError(c, http.StatusInternalServerError, "invalidating user", err)
		return
	}
	slog.Debug("invalidated user")

	c.Redirect(http.StatusTemporaryRedirect, PublicEndpointSignIn)
}

func bindInput(c *gin.Context, input any) error {
	if err := c.ShouldBind(input); err != nil {
		writeError(c, http.StatusBadRequest, "invalid input", err)
		return err
	}
	return nil
}

func setRefreshTokenCookies(c *gin.Context, token string, ttl time.Duration) {
	c.SetCookie(cookiesRefreshToken, token, int(ttl.Seconds()), "/",
		strings.Split(c.Request.Host, ":")[0], false, true)
}
