package v1

import (
	"github.com/adanyl0v/go-pocket-link/internal/domain"
	"github.com/adanyl0v/go-pocket-link/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const (
	PublicSignIn = "/sign-in"
	PublicSignUp = "/sign-up"

	ApiPing   = "/ping"
	ApiSignIn = "/sign-in"
	ApiSignUp = "/sign-up"
	ApiLogOut = "/log-out"

	GroupUser = "/user"
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
	routerGroup.GET(ApiPing, h.handlePing)
	routerGroup.POST(ApiSignUp, h.handleSignUp)
	routerGroup.POST(ApiSignIn, h.handleSignIn)

	protectedGroup := routerGroup.Group("/", h.useAuth)
	{
		protectedGroup.POST(ApiLogOut, h.handleLogOut)

		usersGroup := protectedGroup.Group(GroupUser)
		{
			// operations with user id retrieved from jwt access token
			usersGroup.GET("/", h.handleGetUser)
			usersGroup.PUT("/", h.handleUpdateUser)

			//TODO: implement me
			//usersGroup.DELETE("/") // [?](add email verification)
		}
	}
}

func (h *Handler) handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
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

	user := domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	//TODO (DONE): add credentials validation
	if err := validateCredentials(h.services.Users, input.Name, input.Email, input.Password); err != nil {
		writeError(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	if err := h.services.Users.Save(c, &user); err != nil {
		writeError(c, http.StatusInternalServerError, "failed to save user", err)
		return
	}

	tokens, err := h.services.Tokens.NewTokenPair(user.ID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to create token pair", err)
		return
	}

	if err = h.services.Tokens.SaveRefreshTokenFromString(c, tokens.RefreshToken); err != nil {
		writeError(c, http.StatusInternalServerError, "failed to save refresh token", err)
		return
	}

	setRefreshTokenCookies(c, tokens.RefreshToken, h.services.Tokens.RefreshTokenTTL)

	c.JSON(http.StatusCreated, tokens.AccessToken)
	slog.Debug("signed up", "id", user.ID, "jwt", tokens)

	//TODO: notify user by email
}

func (h *Handler) handleSignIn(c *gin.Context) {
	var input struct {
		Email    string `json:"email" form:"email" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}
	if err := bindInput(c, &input); err != nil {
		return
	}

	//TODO (DONE): add credentials validation
	if err := validateEmailAndPassword(h.services.Users, input.Email, input.Password); err != nil {
		writeError(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	user, err := h.services.Users.GetByCredentials(c, input.Email, input.Password)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to get user", err)
		return
	}

	tokens, err := h.services.Tokens.NewTokenPair(user.ID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to create token pair", err)
		return
	}

	if err = h.services.Tokens.SaveRefreshTokenFromString(c, tokens.RefreshToken); err != nil {
		writeError(c, http.StatusInternalServerError, "failed to save refresh token", err)
		return
	}

	setRefreshTokenCookies(c, tokens.RefreshToken, h.services.Tokens.RefreshTokenTTL)

	c.JSON(http.StatusOK, tokens.AccessToken)
	slog.Debug("signed in", "id", user.ID, "jwt", tokens)

	//TODO: notify user by email, rollback changes if user hasn't logged in
}

func (h *Handler) handleLogOut(c *gin.Context) {
	rawUserID, exists := getRawUserIDFromContext(c)
	if !exists {
		return
	}

	userID, err := parseRawUserID(c, rawUserID)
	if err != nil {
		return
	}

	if err = h.services.Tokens.InvalidateUser(c, userID); err != nil {
		writeError(c, http.StatusInternalServerError, "failed to invalidate user", err)
		return
	}
	slog.Debug("invalidated user")

	//TODO: don't forget to remove me if the frontend will make the redirections
	//c.Redirect(http.StatusTemporaryRedirect, PublicSignIn)
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

func getRawUserIDFromContext(c *gin.Context) (string, bool) {
	raw, exists := c.Get(contextUserID)
	if !exists {
		writeError(c, http.StatusUnauthorized, "no user id in context", nil)
		return "", false
	}
	return raw.(string), true
}

func parseRawUserID(c *gin.Context, raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		writeError(c, http.StatusUnauthorized, "failed to parse user id", err)
		return uuid.Nil, err
	}
	return id, nil
}
