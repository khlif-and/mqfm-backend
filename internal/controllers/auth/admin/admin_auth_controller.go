package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	adminModel "mqfm-backend/internal/models/auth/admin"
	adminService "mqfm-backend/internal/services/auth/admin"
	"mqfm-backend/internal/utils"

)

type AdminAuthController struct {
	service *adminService.AdminAuthService
}

func NewAdminAuthController(s *adminService.AdminAuthService) *AdminAuthController {
	return &AdminAuthController{service: s}
}

func (ctrl *AdminAuthController) Register(c *gin.Context) {
	var admin adminModel.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	if err := ctrl.service.Register(&admin); err != nil {
		utils.Log.Error("Admin registration error: " + err.Error())
		utils.ErrorResponse(c, http.StatusInternalServerError, "Admin registration failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Admin registered successfully", admin)
}

func (ctrl *AdminAuthController) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	token, err := ctrl.service.Login(input.Email, input.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login success", gin.H{
		"token": token,
		"role":  "admin",
	})
}

func (ctrl *AdminAuthController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", nil)
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid update data", err.Error())
		return
	}

	if err := ctrl.service.UpdateAdmin(uint(id), updates); err != nil {
		utils.Log.Error("Admin update error: " + err.Error())
		utils.ErrorResponse(c, http.StatusInternalServerError, "Admin update failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Admin updated successfully", nil)
}

func (ctrl *AdminAuthController) Logout(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Admin logged out successfully", nil)
}