package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"

	"events/server/response"
)

type PingResponse struct {
	*response.Basic
	Version string `json:"version"`
	Service string `json:"service"`
	Commit  string `json:"commit"`
}

type PingHandler struct {
	data []byte
}

func NewPingHandler(service string) gin.HandlerFunc {
	resp := &PingResponse{
		Basic:   response.OK,
		Version: "fixme",
		Service: service,
		Commit:  "fixme",
	}

	// Cache ping response data
	data, _ := jsoniter.Marshal(resp)

	return (&PingHandler{
		data: data,
	}).Handle
}

func (h *PingHandler) Handle(c *gin.Context) {
	c.Data(http.StatusOK, "application/json", h.data)
}
