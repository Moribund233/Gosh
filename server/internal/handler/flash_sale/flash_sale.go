package flash_sale

import (
	"github.com/gin-gonic/gin"
	svc "gosh/internal/service/flash_sale"
	"gosh/pkg/errcode"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

// @Summary List active flash sales
// @Description Get all active flash sales
// @Tags Flash Sales
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /flash-sales [get]
func (h *Handler) ListActive(c *gin.Context) {
	list, err := h.svc.GetActive()
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "get flash sales failed")
		return
	}
	response.Success(c, list)
}
