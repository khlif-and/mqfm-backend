package admin

import (
	"net/http"
	"strconv"
	"time"

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
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	admin := adminModel.Admin{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
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

	token, admin, err := ctrl.service.Login(input.Email, input.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	type LoginResponse struct {
		ID        uint      `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Role      string    `json:"role"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Token     string    `json:"token"`
	}

	responseData := LoginResponse{
		ID:        admin.ID,
		Username:  admin.Username,
		Email:     admin.Email,
		Role:      admin.Role,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
		Token:     token,
	}

	utils.SuccessResponse(c, http.StatusOK, "Login success", responseData)
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

	updatedAdmin, err := ctrl.service.UpdateAdmin(uint(id), updates)
	if err != nil {
		utils.Log.Error("Admin update error: " + err.Error())
		utils.ErrorResponse(c, http.StatusInternalServerError, "Admin update failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Admin updated successfully", updatedAdmin)
}

func (ctrl *AdminAuthController) Logout(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Admin logged out successfully", nil)
}

func (ctrl *AdminAuthController) Me(c *gin.Context) {
	userIDClaim, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var userID uint
	if idFloat, ok := userIDClaim.(float64); ok {
		userID = uint(idFloat)
	} else if idUint, ok := userIDClaim.(uint); ok {
		userID = idUint
	} else {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Invalid token claim", nil)
		return
	}

	admin, err := ctrl.service.GetAdminByID(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Admin not found", err.Error())
		return
	}

	type MeResponse struct {
		ID        uint      `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Role      string    `json:"role"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	response := MeResponse{
		ID:        admin.ID,
		Username:  admin.Username,
		Email:     admin.Email,
		Role:      admin.Role,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
	}

	utils.SuccessResponse(c, http.StatusOK, "User profile retrieved successfully", response)
}