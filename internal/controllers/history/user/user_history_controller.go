package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	historyDto "mqfm-backend/internal/dto/history"
	historyService "mqfm-backend/internal/services/history/user"
	"mqfm-backend/internal/utils"
)

type UserHistoryController struct {
	service *historyService.UserHistoryService
}

func NewUserHistoryController(s *historyService.UserHistoryService) *UserHistoryController {
	return &UserHistoryController{service: s}
}

func (ctrl *UserHistoryController) GetHistory(c *gin.Context) {
	userID := utils.GetUserID(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	histories, err := ctrl.service.GetHistory(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch play history", err.Error())
		return
	}

	var response []historyDto.HistoryResponse
	for _, h := range histories {
		response = append(response, historyDto.HistoryResponse{
			ID:        h.ID,
			UserID:    h.UserID,
			AudioID:   h.AudioID,
			PlayCount: h.PlayCount,
			PlayedAt:  h.PlayedAt,
			CreatedAt: h.CreatedAt,
		})
	}

	utils.SuccessResponse(c, http.StatusOK, "Play history retrieved", response)
}

func (ctrl *UserHistoryController) DeleteHistory(c *gin.Context) {
	audioIDParam := c.Param("audio_id")
	audioID, err := strconv.Atoi(audioIDParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid Audio ID", nil)
		return
	}

	userID := utils.GetUserID(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	if err := ctrl.service.DeleteHistory(userID, uint(audioID)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to delete history", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "History deleted successfully", nil)
}

func (ctrl *UserHistoryController) ClearHistory(c *gin.Context) {
	userID := utils.GetUserID(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	if err := ctrl.service.ClearHistory(userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to clear history", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "All history cleared successfully", nil)
}
