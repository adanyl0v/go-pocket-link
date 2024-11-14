package v1

import (
	"github.com/adanyl0v/go-pocket-link/internal/service"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func (h *Handler) handleGetUser(c *gin.Context) {
	rawUserID, exists := getRawUserIDFromContext(c)
	if !exists {
		return
	}

	userID, err := parseRawUserID(c, rawUserID)
	if err != nil {
		return
	}

	user, err := h.services.Users.Get(c, userID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to get user", err)
		return
	}

	c.JSON(http.StatusOK, user)
	slog.Debug("got user", "user", user)
}

func (h *Handler) handleUpdateUser(c *gin.Context) {
	var input struct {
		Name            string `json:"name" form:"name"`
		Email           string `json:"email" form:"email"`
		Password        string `json:"password" form:"password"`
		CurrentPassword string `json:"current_password" form:"current_password"`
	}
	if err := bindInput(c, &input); err != nil {
		return
	}

	rawUserID, exists := getRawUserIDFromContext(c)
	if !exists {
		return
	}

	userID, err := parseRawUserID(c, rawUserID)
	if err != nil {
		return
	}

	user, err := h.services.Users.Get(c, userID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to get user", err)
		return
	}

	//TODO (DONE): add credentials validation
	if err = validateCredentials(h.services.Users, input.Name, input.Email, input.Password); err != nil {
		writeError(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	user.Name = input.Name
	user.Email = input.Email

	if !h.services.Users.ComparePasswordAndHash(input.CurrentPassword, user.Password) {
		writeError(c, http.StatusBadRequest, "incorrect current password", nil)
		return
	}

	user.Password = input.Password

	if err = h.services.Users.Update(c, &user); err != nil {
		writeError(c, http.StatusInternalServerError, "failed to update user", err)
		return
	}
	slog.Debug("updated user", "id", user.ID)
}

func validateCredentials(s *service.UsersService, name, email, password string) error {
	if err := s.ValidateName(name); err != nil {
		return err
	} else if err = validateEmailAndPassword(s, email, password); err != nil {
		return err
	}
	return nil
}

func validateEmailAndPassword(s *service.UsersService, email, password string) error {
	if err := s.ValidateEmail(email); err != nil {
		return err
	} else if err = s.ValidatePassword(password); err != nil {
		return err
	}
	return nil
}
