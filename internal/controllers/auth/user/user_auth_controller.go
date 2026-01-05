package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	dto "mqfm-backend/internal/dto/auth"
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
	var input dto.RegisterRequest
	if err := c.ShouldBind(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	file, _ := c.FormFile("profile_picture")

	user, err := ctrl.service.Register(input, file)
	if err != nil {
		utils.Log.Error("User registration error: " + err.Error())
		utils.ErrorResponse(c, http.StatusInternalServerError, "User registration failed", err.Error())
		return
	}

	var initials, avatarColor string
	if user.ProfilePicture == "" {
		initials = utils.GetInitials(user.Username)
		avatarColor = utils.GenerateAmbientColor(user.Username)
	}

	response := dto.UserResponse{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		Role:           user.Role,
		ProfilePicture: user.ProfilePicture,
		Initials:       initials,
		AvatarColor:    avatarColor,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", response)
}

func (ctrl *UserAuthController) Login(c *gin.Context) {
	var input dto.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	token, user, err := ctrl.service.Login(input)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	var initials, avatarColor string
	if user.ProfilePicture == "" {
		initials = utils.GetInitials(user.Username)
		avatarColor = utils.GenerateAmbientColor(user.Username)
	}

	response := dto.UserResponse{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		Role:           user.Role,
		ProfilePicture: user.ProfilePicture,
		Initials:       initials,
		AvatarColor:    avatarColor,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		Token:          token,
	}

	utils.SuccessResponse(c, http.StatusOK, "Login success", response)
}

func (ctrl *UserAuthController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", nil)
		return
	}

	var input dto.UpdateUserRequest
	if err := c.ShouldBind(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	file, _ := c.FormFile("profile_picture")

	updatedUser, err := ctrl.service.UpdateUser(uint(id), input, file)
	if err != nil {
		utils.Log.Error("User update error: " + err.Error())
		utils.ErrorResponse(c, http.StatusInternalServerError, "User update failed", err.Error())
		return
	}

	var initials, avatarColor string
	if updatedUser.ProfilePicture == "" {
		initials = utils.GetInitials(updatedUser.Username)
		avatarColor = utils.GenerateAmbientColor(updatedUser.Username)
	}

	response := dto.UserResponse{
		ID:             updatedUser.ID,
		Username:       updatedUser.Username,
		Email:          updatedUser.Email,
		Role:           updatedUser.Role,
		ProfilePicture: updatedUser.ProfilePicture,
		Initials:       initials,
		AvatarColor:    avatarColor,
		CreatedAt:      updatedUser.CreatedAt,
		UpdatedAt:      updatedUser.UpdatedAt,
	}

	utils.SuccessResponse(c, http.StatusOK, "User updated successfully", response)
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

	var initials, avatarColor string
	if user.ProfilePicture == "" {
		initials = utils.GetInitials(user.Username)
		avatarColor = utils.GenerateAmbientColor(user.Username)
	}

	response := dto.UserResponse{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		Role:           user.Role,
		ProfilePicture: user.ProfilePicture,
		Initials:       initials,
		AvatarColor:    avatarColor,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}

	utils.SuccessResponse(c, http.StatusOK, "User profile retrieved successfully", response)
}