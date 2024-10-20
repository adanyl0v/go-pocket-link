package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go-pocket-link/internal/service"
	"go-pocket-link/pkg/errb"
	"net/http"
)

type Handler interface {
	InitEndpoints(parent *gin.RouterGroup)

	Ping() gin.HandlerFunc

	GetAllUsers() gin.HandlerFunc
	GetUserById() gin.HandlerFunc
	CreateUser() gin.HandlerFunc
	UpdateUser() gin.HandlerFunc
	DeleteUser() gin.HandlerFunc

	GetAllLinks() gin.HandlerFunc
	GetLinkById() gin.HandlerFunc
	CreateLink() gin.HandlerFunc
	UpdateLink() gin.HandlerFunc
	DeleteLink() gin.HandlerFunc

	GetAllSessions() gin.HandlerFunc
	GetSessionById() gin.HandlerFunc
	CreateSession() gin.HandlerFunc
	UpdateSession() gin.HandlerFunc
	DeleteSession() gin.HandlerFunc

	AccessJWT() gin.HandlerFunc
	RefreshJWT() gin.HandlerFunc
	InvalidateJWT() gin.HandlerFunc

	EmailSendSuccessfulSignIn() gin.HandlerFunc
	EmailSendSuccessfulSignUp() gin.HandlerFunc
}

type handlerImpl struct {
	services *service.Services
}

func NewHandler(services *service.Services) Handler {
	return &handlerImpl{services: services}
}

func (h *handlerImpl) InitEndpoints(parent *gin.RouterGroup) {
	parent.GET("/ping", h.Ping())

	usersGroup := parent.Group("/users")
	{
		usersGroup.GET("/", h.GetAllUsers())
		usersGroup.GET("/:id", h.GetUserById())
		usersGroup.POST("/", h.CreateUser())
		usersGroup.PUT("/:id", h.UpdateUser())
		usersGroup.DELETE("/:id", h.DeleteUser())
	}

	linksGroup := parent.Group("/links")
	{
		linksGroup.GET("/", h.GetAllLinks())
		linksGroup.GET("/:id", h.GetLinkById())
		linksGroup.POST("/", h.CreateLink())
		linksGroup.PUT("/:id", h.UpdateLink())
		linksGroup.DELETE("/:id", h.DeleteLink())
	}

	sessionsGroup := parent.Group("/sessions")
	{
		sessionsGroup.GET("/", h.GetAllSessions())
		sessionsGroup.GET("/:id", h.GetSessionById())
		sessionsGroup.POST("/", h.CreateSession())
		sessionsGroup.PUT("/:id", h.UpdateSession())
		sessionsGroup.DELETE("/:id", h.DeleteSession())
	}

	authGroup := parent.Group("/auth")
	{
		jwtGroup := authGroup.Group("/jwt")
		{
			jwtGroup.GET("/:token", h.AccessJWT())
			jwtGroup.PUT("/:token", h.RefreshJWT())
			jwtGroup.DELETE("/:token", h.InvalidateJWT())
		}
	}

	emailGroup := parent.Group("/email")
	{
		emailGroup.POST("/send/successful-sign-in", h.EmailSendSuccessfulSignIn())
		emailGroup.POST("/send/successful-sign-up", h.EmailSendSuccessfulSignUp())
	}
}

func (h *handlerImpl) Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		writeMessage(c, http.StatusOK, gin.H{keyMessage: "pong"})
	}
}

func urlParamID(c *gin.Context) (uuid.UUID, error) {
	idStr := c.Param("id")
	if idStr == "" {
		return uuid.Nil, errb.Errorf("no id provided")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, errb.Errorf("invalid id provided")
	}
	return id, nil
}

func writeResponse(c *gin.Context, status int, response any) {
	c.JSON(status, response)
	log.Infoln(response)
}

const keyMessage = "message"

func writeMessage(c *gin.Context, status int, message gin.H) {
	c.JSON(status, message)
	m, exist := message[keyMessage]
	if exist {
		log.Infoln(m)
	}
}

func writeError(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"error": err.Error()})
	log.Errorln(err)
}
