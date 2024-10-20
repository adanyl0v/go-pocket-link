package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-pocket-link/internal/domain"
	"go-pocket-link/pkg/errb"
	"net/http"
	"time"
)

func (h *handlerImpl) GetAllSessions() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessions, err := h.services.Sessions.Repository().GetAll(c)
		if err != nil {
			writeError(c, http.StatusInternalServerError, errb.Errorf("fetching all sessions: %v", err))
			return
		}
		writeResponse(c, http.StatusOK, sessions)
	}
}

func (h *handlerImpl) GetSessionById() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := urlParamID(c)
		if err != nil {
			writeError(c, http.StatusBadRequest, err)
			return
		}

		session, err := h.services.Sessions.Repository().GetByID(c, id)
		if err != nil {
			writeError(c, http.StatusInternalServerError, errb.Errorf("fetching session by id: %v", err))
			return
		}
		writeResponse(c, http.StatusOK, session)
	}
}

type createSessionInput struct {
	UserID       uuid.UUID `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (h *handlerImpl) CreateSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input createSessionInput
		if err := c.ShouldBindJSON(&input); err != nil {
			writeError(c, http.StatusBadRequest, errb.Errorf("invalid input: %v", err))
			return
		}

		service := h.services.Sessions
		session := domain.Session{
			UserID:       input.UserID,
			RefreshToken: input.RefreshToken,
			ExpiresAt:    input.ExpiresAt,
		}
		err := service.Repository().Save(c, &session)
		if err != nil {
			writeError(c, http.StatusInternalServerError, errb.Errorf("creating session: %v", err))
			return
		}

		writeResponse(c, http.StatusCreated, session)
	}
}

type updateSessionInput struct {
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (h *handlerImpl) UpdateSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := urlParamID(c)
		if err != nil {
			writeError(c, http.StatusBadRequest, err)
			return
		}

		var input updateSessionInput
		if err = c.ShouldBindJSON(&input); err != nil {
			writeError(c, http.StatusBadRequest, errb.Errorf("invalid input: %v", err))
			return
		}

		service := h.services.Sessions
		session, err := service.Repository().GetByID(c, id)
		if err != nil {
			writeError(c, http.StatusInternalServerError, errb.Errorf("fetching session: %v", err))
			return
		}

		err = service.Repository().Update(c, &session)
		if err != nil {
			writeError(c, http.StatusInternalServerError, errb.Errorf("updating session: %v", err))
			return
		}

		writeResponse(c, http.StatusOK, session)
	}
}

func (h *handlerImpl) DeleteSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := urlParamID(c)
		if err != nil {
			writeError(c, http.StatusBadRequest, err)
			return
		}

		err = h.services.Sessions.Repository().Delete(c, id)
		if err != nil {
			writeError(c, http.StatusInternalServerError, errb.Errorf("deleting session: %v", err))
			return
		}

		writeMessage(c, http.StatusOK, gin.H{keyMessage: fmt.Sprint("deleted session with id=", id.String())})

	}
}
