package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/dto/request"
	resp "mqfm-backend/internal/shared/response"
	"mqfm-backend/internal/shared/security"
)

type OTPHandler struct {
	otpService  port.OTPService
	userService port.UserAuthService
}

func NewOTPHandler(otpSvc port.OTPService, userSvc port.UserAuthService) *OTPHandler {
	return &OTPHandler{otpService: otpSvc, userService: userSvc}
}

func (h *OTPHandler) SendOTP(c *gin.Context) {
	var input request.SendOTPRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	if err := h.otpService.SendOTP(input.Email); err != nil {
		resp.Error(c, http.StatusTooManyRequests, constant.MsgOTPSendFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgOTPSendOK, nil)
}

func (h *OTPHandler) VerifyOTP(c *gin.Context) {
	var input request.VerifyOTPRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	if err := h.otpService.VerifyOTP(input.Email, input.Code); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgOTPVerifyFail, err.Error())
		return
	}

	user, err := h.userService.GetUserByID(security.GetUserID(c))
	if err == nil && user.Email == input.Email {
		resp.Success(c, http.StatusOK, constant.MsgOTPVerifyOK, toUserResponse(user, ""))
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgOTPVerifyOK, nil)
}
