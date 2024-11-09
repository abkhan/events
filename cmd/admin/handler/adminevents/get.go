package adminevents

import (
	"events/server/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getLastEvent(c *gin.Context) {
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

	// TDODO: fix the response data
	response.SuccessData(c, "Status: Good")
}
