package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/internal/shared/dto/response"
	"mqfm-backend/internal/shared/logger"
	resp "mqfm-backend/internal/shared/response"
	"mqfm-backend/internal/shared/security"
)

type AuthHandler struct {
	service port.AdminAuthService
}

func NewAuthHandler(s port.AdminAuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input request.AdminRegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	admin, err := h.service.Register(input)
	if err != nil {
		logger.Error("admin register: " + err.Error())
		resp.Error(c, http.StatusInternalServerError, constant.MsgAdminRegisterFail, err.Error())
		return
	}

	resp.Success(c, http.StatusCreated, constant.MsgAdminRegisterOK, response.AdminResponse{
		ID:        admin.ID,
		Username:  admin.Username,
		Email:     admin.Email,
		Role:      admin.Role,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input request.AdminLoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	token, admin, err := h.service.Login(input)
	if err != nil {
		resp.Error(c, http.StatusUnauthorized, constant.MsgAdminLoginFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgAdminLoginOK, response.AdminResponse{
		ID:        admin.ID,
		Username:  admin.Username,
		Email:     admin.Email,
		Role:      admin.Role,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
		Token:     token,
	})
}

func (h *AuthHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, ok := parseUint(idParam)
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	admin, err := h.service.UpdateAdmin(id, updates)
	if err != nil {
		logger.Error("admin update: " + err.Error())
		resp.Error(c, http.StatusInternalServerError, constant.MsgAdminUpdateFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgAdminUpdateOK, response.AdminResponse{
		ID:        admin.ID,
		Username:  admin.Username,
		Email:     admin.Email,
		Role:      admin.Role,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	resp.Success(c, http.StatusOK, constant.MsgAdminLogoutOK, nil)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := security.GetUserID(c)
	if userID == 0 {
		resp.Error(c, http.StatusUnauthorized, constant.MsgUnauthorized, nil)
		return
	}

	admin, err := h.service.GetAdminByID(userID)
	if err != nil {
		resp.Error(c, http.StatusNotFound, constant.MsgAdminNotFound, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgAdminProfileOK, response.AdminResponse{
		ID:        admin.ID,
		Username:  admin.Username,
		Email:     admin.Email,
		Role:      admin.Role,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
	})
}

func parseUint(s string) (uint, bool) {
	var id uint
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, false
		}
		id = id*10 + uint(c-'0')
	}
	if len(s) == 0 {
		return 0, false
	}
	return id, true
}
