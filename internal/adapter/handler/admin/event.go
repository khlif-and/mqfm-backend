package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/internal/shared/dto/response"
	"mqfm-backend/internal/shared/helper"
	resp "mqfm-backend/internal/shared/response"
)

type EventHandler struct {
	service port.EventService
}

func NewEventHandler(s port.EventService) *EventHandler {
	return &EventHandler{service: s}
}

func (h *EventHandler) Create(c *gin.Context) {
	var input request.CreateEventRequest
	if err := c.ShouldBind(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	file, _ := c.FormFile("image")

	event, err := h.service.Create(input, file)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgEventCreateFail, err.Error())
		return
	}

	resp.Success(c, http.StatusCreated, constant.MsgEventCreateOK, toEventResponse(event))
}

func (h *EventHandler) FindAll(c *gin.Context) {
	events, err := h.service.FindAll()
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgEventListFail, err.Error())
		return
	}

	var result []response.EventResponse
	for _, e := range events {
		result = append(result, toEventResponseVal(e))
	}

	resp.Success(c, http.StatusOK, constant.MsgEventListOK, result)
}

func (h *EventHandler) FindByID(c *gin.Context) {
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

	resp.Success(c, http.StatusOK, constant.MsgEventGetOK, toEventResponse(event))
}

func (h *EventHandler) Update(c *gin.Context) {
	id, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	var input request.UpdateEventRequest
	if err := c.ShouldBind(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	file, _ := c.FormFile("image")

	event, err := h.service.Update(id, input, file)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgEventUpdateFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgEventUpdateOK, toEventResponse(event))
}

func (h *EventHandler) Delete(c *gin.Context) {
	id, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	if err := h.service.Delete(id); err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgEventDeleteFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgEventDeleteOK, nil)
}
