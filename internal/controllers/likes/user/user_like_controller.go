package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	likeService "mqfm-backend/internal/services/likes/user"
	"mqfm-backend/internal/utils"

)

type UserLikeController struct {
	service *likeService.UserLikeService
}

func NewUserLikeController(s *likeService.UserLikeService) *UserLikeController {
	return &UserLikeController{service: s}
}

func (ctrl *UserLikeController) Like(c *gin.Context) {
	var input struct {
		AudioID uint `json:"audio_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	userID := utils.GetUserID(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	if err := ctrl.service.LikeAudio(userID, input.AudioID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to like audio", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Audio liked successfully", nil)
}

func (ctrl *UserLikeController) Unlike(c *gin.Context) {
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

	if err := ctrl.service.UnlikeAudio(userID, uint(audioID)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to unlike audio", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Audio unliked successfully", nil)
}

func (ctrl *UserLikeController) GetLikes(c *gin.Context) {
	userID := utils.GetUserID(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	likes, err := ctrl.service.GetLikedAudios(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch liked audios", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Liked audios retrieved", likes)
}