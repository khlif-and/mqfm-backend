package user

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	userModel "mqfm-backend/internal/models/auth/user"
	userService "mqfm-backend/internal/services/auth/user"
	"mqfm-backend/internal/utils"

)

type UserAuthController struct {
	service *userService.UserAuthService
}

func NewUserAuthController(s *userService.UserAuthService) *UserAuthController {
	return &UserAuthController{service: s}
}

func (ctrl *UserAuthController) Register(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	user := userModel.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	if err := ctrl.service.Register(&user); err != nil {
		utils.Log.Error("User registration error: " + err.Error())
		utils.ErrorResponse(c, http.StatusInternalServerError, "User registration failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", user)
}

func (ctrl *UserAuthController) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	token, user, err := ctrl.service.Login(input.Email, input.Password)
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
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Token:     token,
	}

	utils.SuccessResponse(c, http.StatusOK, "Login success", responseData)
}

func (ctrl *UserAuthController) Update(c *gin.Context) {
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

	updatedUser, err := ctrl.service.UpdateUser(uint(id), updates)
	if err != nil {
		utils.Log.Error("User update error: " + err.Error())
		utils.ErrorResponse(c, http.StatusInternalServerError, "User update failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User updated successfully", updatedUser)
}

func (ctrl *UserAuthController) Logout(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "User logged out successfully", nil)
}

func (ctrl *UserAuthController) Me(c *gin.Context) {
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

	user, err := ctrl.service.GetUserByID(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", err.Error())
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
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	utils.SuccessResponse(c, http.StatusOK, "User profile retrieved successfully", response)
}