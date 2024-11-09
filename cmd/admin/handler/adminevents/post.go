package adminevents

import (
	"events/server/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func (h *Handler) postEvent(c *gin.Context) {
	// Validate Parameters
	locationParam := c.Param("location")
	if locationParam == "" {
		response.AbortInvalidParams(c, errParameterNotFound)
		return
	}
	if err := validate.Var(locationParam, "required,lt=33"); err != nil {
		response.AbortInvalidParams(c, err)
		return
	}

	// TODO: validate if API-Key is valid

	// TODO: fix the return
	response.SuccessNoContent(c)
}
