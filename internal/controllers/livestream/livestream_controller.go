package livestream

import (
	"net/http"

	"github.com/gin-gonic/gin"

	lsService "mqfm-backend/internal/services/livestream"
	"mqfm-backend/internal/utils"

)

type LiveStreamController struct {
	service *lsService.LiveStreamService
}

func NewLiveStreamController(s *lsService.LiveStreamService) *LiveStreamController {
	return &LiveStreamController{service: s}
}

func (ctrl *LiveStreamController) GetStatus(c *gin.Context) {
	status, err := ctrl.service.GetStatus()
	if err != nil {
		utils.SuccessResponse(c, http.StatusOK, "No data yet", gin.H{"is_live": false})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Live stream status", status)
}