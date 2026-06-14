package flash_sale

import (
	"github.com/gin-gonic/gin"
	svc "gosh/internal/service/flash_sale"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

func (h *Handler) ListActive(c *gin.Context) {
	list, err := h.svc.GetActive()
	if err != nil {
		response.InternalError(c, "get flash sales failed")
		return
	}
	response.Success(c, list)
}
