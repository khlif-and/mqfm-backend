package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/internal/shared/dto/response"
	"mqfm-backend/internal/shared/helper"
	"mqfm-backend/internal/shared/logger"
	resp "mqfm-backend/internal/shared/response"
	"mqfm-backend/internal/shared/security"
)

type AuthHandler struct {
	service port.UserAuthService
}

func NewAuthHandler(s port.UserAuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input request.UserRegisterRequest
	if err := c.ShouldBind(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	file, _ := c.FormFile("profile_picture")

	u, err := h.service.Register(input, file)
	if err != nil {
		logger.Error("user register: " + err.Error())
		resp.Error(c, http.StatusInternalServerError, constant.MsgUserRegisterFail, err.Error())
		return
	}

	resp.Success(c, http.StatusCreated, constant.MsgUserRegisterOK, toUserResponse(u, ""))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input request.UserLoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	token, u, err := h.service.Login(input)
	if err != nil {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUserLoginFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgUserLoginOK, toUserResponse(u, token))
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var input request.GoogleLoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	token, u, err := h.service.GoogleLogin(input)
	if err != nil {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUserGoogleLoginFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgUserGoogleLoginOK, toUserResponse(u, token))
}

func (h *AuthHandler) Update(c *gin.Context) {
	id, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	var input request.UpdateUserRequest
	if err := c.ShouldBind(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	file, _ := c.FormFile("profile_picture")

	u, err := h.service.UpdateUser(id, input, file)
	if err != nil {
		logger.Error("user update: " + err.Error())
		resp.Error(c, http.StatusInternalServerError, constant.MsgUserUpdateFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgUserUpdateOK, toUserResponse(u, ""))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	resp.Success(c, http.StatusOK, constant.MsgUserLogoutOK, nil)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	u, err := h.service.GetUserByID(userID)
	if err != nil {
		resp.Error(c, http.StatusNotFound, constant.MsgUserNotFound, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgUserProfileOK, toUserResponse(u, ""))
}

func toUserResponse(u *entity.User, token string) response.UserResponse {
	var initials, avatarColor string
	if u.ProfilePicture == "" {
		initials = helper.GetInitials(u.Username)
		avatarColor = helper.GenerateAmbientColor(u.Username)
	}

	return response.UserResponse{
		ID:             u.ID,
		Username:       u.Username,
		Email:          u.Email,
		Role:           u.Role,
		ProfilePicture: u.ProfilePicture,
		Initials:       initials,
		AvatarColor:    avatarColor,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		Token:          token,
	}
}
