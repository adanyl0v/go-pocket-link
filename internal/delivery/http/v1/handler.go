package v1

import (
	"github.com/gin-gonic/gin"
	"go-pocket-link/internal/domain"
	"go-pocket-link/internal/service"
	"log/slog"
	"net/http"
)

const logError = "error"

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
		c.JSON(http.StatusBadRequest, gin.H{logError: "invalid input"})
		slog.Error("binding input", logError, err.Error())
		return
	}

	user := domain.User{Name: input.Name, Email: input.Email, Password: input.Password}
	if err := h.services.Users.Save(c, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{logError: "saving user"})
		slog.Error("saving user", logError, err.Error())
		return
	}

	//TODO: create JWT session

	c.JSON(http.StatusCreated, user)
	slog.Debug("signed up", "user", user)
}

func (h *Handler) signIn(c *gin.Context) {
	var input struct {
		Email    string `json:"email" form:"email" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}
	slog.Debug("sign in input", "input", input)

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{logError: "invalid input"})
		slog.Error("binding input", logError, err.Error())
		return
	}

	user, err := h.services.Users.GetByCredentials(c, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{logError: "getting user"})
		slog.Error("getting user", logError, err.Error())
		return
	}

	//TODO: refresh JWT session

	c.JSON(http.StatusOK, user)
	slog.Debug("signed in", "user", user)
}
