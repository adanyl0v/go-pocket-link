package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-pocket-link/internal/domain"
	"go-pocket-link/pkg/errb"
	"net/http"
)

func (h *handlerImpl) GetAllLinks() gin.HandlerFunc {
	return func(c *gin.Context) {
		links, err := h.services.Links.Repository().GetAll(c)
		if err != nil {
			writeError(c, http.StatusInternalServerError, errb.Errorf("fetching all links: %v", err))
			return
		}
		writeResponse(c, http.StatusOK, links)
	}
}

func (h *handlerImpl) GetLinkById() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := urlParamID(c)
		if err != nil {
			writeError(c, http.StatusBadRequest, err)
			return
		}

		link, err := h.services.Links.Repository().GetByID(c, id)
		if err != nil {
			writeError(c, http.StatusInternalServerError, errb.Errorf("fetching link by id: %v", err))
			return
		}
		writeResponse(c, http.StatusOK, link)
	}
}

type createLinkInput struct {
	UserID uuid.UUID `json:"user_id"`
	Title  string    `json:"title,omitempty"`
	URL    string    `json:"url"`
}

func (h *handlerImpl) CreateLink() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input createLinkInput
		if err := c.ShouldBindJSON(&input); err != nil {
			writeError(c, http.StatusBadRequest, errb.Errorf("invalid input: %v", err))
			return
		}

		link := domain.Link{UserID: input.UserID, Title: input.Title, URL: input.URL}
		err := h.services.Links.Repository().Save(c, &link)
		if err != nil {
			writeError(c, http.StatusInternalServerError, errb.Errorf("creating link: %v", err))
			return
		}

		writeResponse(c, http.StatusCreated, link)
	}
}

type updateLinkInput struct {
	Title string `json:"title,omitempty"`
	URL   string `json:"url,omitempty"`
}

func (h *handlerImpl) UpdateLink() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := urlParamID(c)
		if err != nil {
			writeError(c, http.StatusBadRequest, err)
			return
		}

		var input updateLinkInput
		if err = c.ShouldBindJSON(&input); err != nil {
			writeError(c, http.StatusBadRequest, errb.Errorf("invalid input: %v", err))
			return
		}

		link, err := h.services.Links.Repository().GetByID(c, id)
		if err != nil {
			writeError(c, http.StatusInternalServerError, errb.Errorf("fetching link by id: %v", err))
			return
		}

		if input.Title == "" && input.URL == "" {
			writeError(c, http.StatusBadRequest, errb.Errorf("no updates provided"))
			return
		}

		if input.Title != "" {
			link.Title = input.Title
		}
		if input.URL != "" {
			link.URL = input.URL
		}
		err = h.services.Links.Repository().Update(c, &link)
		if err != nil {
			writeError(c, http.StatusInternalServerError, errb.Errorf("updating link: %v", err))
			return
		}

		writeResponse(c, http.StatusOK, link)
	}
}

func (h *handlerImpl) DeleteLink() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := urlParamID(c)
		if err != nil {
			writeError(c, http.StatusBadRequest, err)
			return
		}

		err = h.services.Links.Repository().Delete(c, id)
		if err != nil {
			writeError(c, http.StatusInternalServerError, errb.Errorf("deleting link: %v", err))
			return
		}

		writeMessage(c, http.StatusOK, gin.H{keyMessage: fmt.Sprint("deleted link with id=", id.String())})
	}
}
