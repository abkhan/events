package adminevents

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var errParameterNotFound = errors.New("Parameter not found")

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) BindEndpoint(rg *gin.RouterGroup) {
	rg.GET("/:location/:entity/:status", h.getLastEvent)
	rg.POST("/:location/:entity/:status", h.postEvent)
}
