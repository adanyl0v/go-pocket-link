package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-pocket-link/internal/domain"
	"go-pocket-link/pkg/errb"
	"net/http"
)

func (h *handlerImpl) GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := h.services.Users.Repository().GetAll(c)
		if err != nil {
			writeError(c, http.StatusInternalServerError, errb.Errorf("fetching all users: %v", err))
			return
		}
		writeResponse(c, http.StatusOK, users)
	}
}

func (h *handlerImpl) GetUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := urlParamID(c)
		if err != nil {
			writeError(c, http.StatusBadRequest, err)
			return
		}

		user, err := h.services.Users.Repository().GetByID(c, id)
		if err != nil {
			writeError(c, http.StatusInternalServerError, errb.Errorf("fetching user by id: %v", err))
			return
		}
		writeResponse(c, http.StatusOK, user)
	}
}

type createUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *handlerImpl) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input createUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			writeError(c, http.StatusBadRequest, errb.Errorf("invalid input: %v", err))
			return
		}

		service := h.services.Users
		user := domain.User{
			Name:     input.Name,
			Email:    input.Email,
			Password: service.HashPassword(input.Password),
		}
		err := service.Repository().Save(c, &user)
		if err != nil {
			writeError(c, http.StatusInternalServerError, errb.Errorf("creating user: %v", err))
			return
		}

		writeResponse(c, http.StatusCreated, user)
	}
}

type updateUserInput struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (h *handlerImpl) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := urlParamID(c)
		if err != nil {
			writeError(c, http.StatusBadRequest, err)
			return
		}

		var input updateUserInput
		if err = c.ShouldBindJSON(&input); err != nil {
			writeError(c, http.StatusBadRequest, errb.Errorf("invalid input: %v", err))
			return
		}

		service := h.services.Users
		user, err := service.Repository().GetByID(c, id)
		if err != nil {
			writeError(c, http.StatusInternalServerError, errb.Errorf("fetching user: %v", err))
			return
		}

		if input.Name == "" && input.Email == "" && input.Password == "" {
			writeError(c, http.StatusBadRequest, errb.Errorf("no updates provided"))
			return
		}

		if input.Name != "" {
			user.Name = input.Name
		}
		if input.Email != "" {
			user.Email = input.Email
		}
		if input.Password != "" {
			user.Password = service.HashPassword(input.Password)
		}
		err = service.Repository().Update(c, &user)
		if err != nil {
			writeError(c, http.StatusInternalServerError, errb.Errorf("updating user: %v", err))
			return
		}

		writeResponse(c, http.StatusOK, user)
	}
}

func (h *handlerImpl) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := urlParamID(c)
		if err != nil {
			writeError(c, http.StatusBadRequest, err)
			return
		}

		err = h.services.Users.Repository().Delete(c, id)
		if err != nil {
			writeError(c, http.StatusInternalServerError, errb.Errorf("deleting user: %v", err))
			return
		}

		writeMessage(c, http.StatusOK, gin.H{keyMessage: fmt.Sprint("deleted user with id=", id.String())})
	}
}
