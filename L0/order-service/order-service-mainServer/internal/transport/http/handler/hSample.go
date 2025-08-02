package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type respSample struct {
	Field1 string
	Field2 string
	Field3 int
}

func (h *Handler) hSample(c *gin.Context) {
	h.Logger.Debug().Str("handler", "hSample").Msg("Запущен hSample")

	resp := respSample{
		Field1: "field1",
		Field2: "feild2",
		Field3: 3,
	}

	c.IndentedJSON(http.StatusOK, resp)

}
