package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/dto/response"
	"mqfm-backend/internal/shared/helper"
	resp "mqfm-backend/internal/shared/response"
	"mqfm-backend/internal/shared/security"
)

type EventHandler struct {
	service port.EventService
}

func NewEventHandler(s port.EventService) *EventHandler {
	return &EventHandler{service: s}
}

func (h *EventHandler) GetUpcoming(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	events, err := h.service.GetUpcoming(limit)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgEventListFail, err.Error())
		return
	}

	userID := security.GetUserID(c)
	var result []response.EventResponse
	for _, e := range events {
		rsvpCount, _ := h.service.GetRSVPCount(e.ID)
		hasRSVP := false
		if userID > 0 {
			rsvps, _ := h.service.GetUserRSVPs(userID)
			for _, r := range rsvps {
				if r.EventID == e.ID {
					hasRSVP = true
					break
				}
			}
		}
		result = append(result, response.EventResponse{
			ID:          e.ID,
			Title:       e.Title,
			Description: e.Description,
			EventDate:   e.EventDate,
			Location:    e.Location,
			Image:       e.Image,
			RSVPCount:   rsvpCount,
			HasRSVP:     hasRSVP,
			CreatedAt:   e.CreatedAt,
			UpdatedAt:   e.UpdatedAt,
		})
	}

	resp.Success(c, http.StatusOK, constant.MsgEventListOK, result)
}

func (h *EventHandler) GetByID(c *gin.Context) {
	id, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	event, err := h.service.FindByID(id)
	if err != nil {
		resp.Error(c, http.StatusNotFound, constant.MsgEventNotFound, nil)
		return
	}

	rsvpCount, _ := h.service.GetRSVPCount(event.ID)
	userID := security.GetUserID(c)
	hasRSVP := false
	if userID > 0 {
		rsvps, _ := h.service.GetUserRSVPs(userID)
		for _, r := range rsvps {
			if r.EventID == event.ID {
				hasRSVP = true
				break
			}
		}
	}

	resp.Success(c, http.StatusOK, constant.MsgEventGetOK, response.EventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		EventDate:   event.EventDate,
		Location:    event.Location,
		Image:       event.Image,
		RSVPCount:   rsvpCount,
		HasRSVP:     hasRSVP,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	})
}

func (h *EventHandler) RSVP(c *gin.Context) {
	eventID, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	if err := h.service.RSVP(userID, eventID); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgEventRSVPFail, err.Error())
		return
	}

	resp.Success(c, http.StatusCreated, constant.MsgEventRSVPOK, nil)
}

func (h *EventHandler) CancelRSVP(c *gin.Context) {
	eventID, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	if err := h.service.CancelRSVP(userID, eventID); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgEventRSVPCancelFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgEventRSVPCancelOK, nil)
}

func (h *EventHandler) GetMyRSVPs(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	rsvps, err := h.service.GetUserRSVPs(userID)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgEventRSVPListFail, err.Error())
		return
	}

	var result []response.EventRSVPResponse
	for _, r := range rsvps {
		item := response.EventRSVPResponse{
			ID:        r.ID,
			UserID:    r.UserID,
			EventID:   r.EventID,
			CreatedAt: r.CreatedAt,
		}
		if r.Event != nil {
			item.Event = &response.EventResponse{
				ID:        r.Event.ID,
				Title:     r.Event.Title,
				EventDate: r.Event.EventDate,
				Location:  r.Event.Location,
				Image:     r.Event.Image,
			}
		}
		result = append(result, item)
	}

	resp.Success(c, http.StatusOK, constant.MsgEventRSVPListOK, result)
}
