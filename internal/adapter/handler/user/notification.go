package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/internal/shared/dto/response"
	"mqfm-backend/internal/shared/helper"
	resp "mqfm-backend/internal/shared/response"
	"mqfm-backend/internal/shared/security"
)

type NotificationHandler struct {
	service port.NotificationService
}

func NewNotificationHandler(s port.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: s}
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	notifications, err := h.service.GetByUser(userID, page, limit)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgNotificationListFail, err.Error())
		return
	}

	var result []response.NotificationResponse
	for _, n := range notifications {
		result = append(result, response.NotificationResponse{
			ID:          n.ID,
			UserID:      n.UserID,
			Title:       n.Title,
			Body:        n.Body,
			Type:        n.Type,
			ReferenceID: n.ReferenceID,
			IsRead:      n.IsRead,
			CreatedAt:   n.CreatedAt,
		})
	}

	resp.Success(c, http.StatusOK, constant.MsgNotificationListOK, result)
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	if err := h.service.MarkAsRead(id, userID); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgNotificationReadFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgNotificationReadOK, nil)
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	if err := h.service.MarkAllAsRead(userID); err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgNotificationReadFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgNotificationReadAllOK, nil)
}

func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	count, err := h.service.CountUnread(userID)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgInternalError, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgNotificationCountOK, response.UnreadCountResponse{Count: count})
}

func (h *NotificationHandler) GetSetting(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	setting, err := h.service.GetSetting(userID)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgInternalError, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgNotificationSettingOK, response.NotificationSettingResponse{
		DailyReminder: setting.DailyReminder,
		NewContent:    setting.NewContent,
		EventReminder: setting.EventReminder,
	})
}

func (h *NotificationHandler) UpdateSetting(c *gin.Context) {
	var input request.UpdateNotificationSettingRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	setting, err := h.service.UpdateSetting(userID, input)
	if err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgNotificationSettingFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgNotificationSettingUpdateOK, response.NotificationSettingResponse{
		DailyReminder: setting.DailyReminder,
		NewContent:    setting.NewContent,
		EventReminder: setting.EventReminder,
	})
}
