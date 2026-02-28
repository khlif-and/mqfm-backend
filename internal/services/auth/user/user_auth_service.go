package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"

	dto "mqfm-backend/internal/dto/auth"
	userModel "mqfm-backend/internal/models/auth/user"
	userRepo "mqfm-backend/internal/repositories/auth/user"
	"mqfm-backend/internal/utils"
)

type UserAuthService struct {
	repo userRepo.UserAuthRepository
}

func NewUserAuthService(repo userRepo.UserAuthRepository) *UserAuthService {
	return &UserAuthService{repo: repo}
}

func (s *UserAuthService) Register(req dto.RegisterRequest, file *multipart.FileHeader) (*userModel.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Log.Error("Failed to hash user password")
		return nil, err
	}

	var profilePicturePath string
	if file != nil {
		filename := utils.GenerateUniqueFilename(file.Filename)
		path := "uploads/profiles/" + filename
		if err := utils.SaveUploadedFile(file, path); err != nil {
			utils.Log.Error("Failed to save profile picture: " + err.Error())
		} else {
			profilePicturePath = path
		}
	}

	user := userModel.User{
		Username:       req.Username,
		Email:          req.Email,
		Password:       string(hashedPassword),
		ProfilePicture: profilePicturePath,
		Role:           "user",
		Provider:       "local",
	}

	if err := s.repo.Create(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserAuthService) Login(req dto.LoginRequest) (string, *userModel.User, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		utils.Log.Warn("User login attempt failed: email not found")
		return "", nil, errors.New("invalid user credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Log.Warn("User login attempt failed: incorrect password")
		return "", nil, errors.New("invalid user credentials")
	}

	token, err := utils.GenerateToken(user.ID, "user")
	if err != nil {
		utils.Log.Error("Failed to generate user JWT token: " + err.Error())
		return "", nil, err
	}

	return token, user, nil
}

type googleTokenInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Aud           string `json:"aud"`
}

func (s *UserAuthService) GoogleLogin(req dto.GoogleLoginRequest) (string, *userModel.User, error) {
	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + req.IDToken)
	if err != nil {
		return "", nil, errors.New("failed to verify google token")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, errors.New("invalid google token")
	}

	var tokenInfo googleTokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return "", nil, errors.New("failed to parse google token info")
	}

	clientIDs := strings.Split(os.Getenv("GOOGLE_CLIENT_ID"), ",")
	validAud := false
	for _, id := range clientIDs {
		if strings.TrimSpace(id) == tokenInfo.Aud {
			validAud = true
			break
		}
	}
	if !validAud {
		return "", nil, errors.New("invalid google token audience")
	}

	if tokenInfo.EmailVerified != "true" {
		return "", nil, errors.New("google email not verified")
	}

	user, err := s.repo.FindByProviderID("google", tokenInfo.Sub)
	if err != nil {
		username := tokenInfo.Name
		if username == "" {
			username = strings.Split(tokenInfo.Email, "@")[0]
		}

		existingUser, emailErr := s.repo.FindByEmail(tokenInfo.Email)
		if emailErr == nil {
			existingUser.Provider = "google"
			existingUser.ProviderID = tokenInfo.Sub
			if tokenInfo.Picture != "" && existingUser.ProfilePicture == "" {
				existingUser.ProfilePicture = tokenInfo.Picture
			}
			updates := map[string]interface{}{
				"provider":    "google",
				"provider_id": tokenInfo.Sub,
			}
			if tokenInfo.Picture != "" && existingUser.ProfilePicture == "" {
				updates["profile_picture"] = tokenInfo.Picture
			}
			s.repo.Update(existingUser.ID, updates)
			user = existingUser
		} else {
			newUser := userModel.User{
				Username:       username,
				Email:          tokenInfo.Email,
				ProfilePicture: tokenInfo.Picture,
				Role:           "user",
				Provider:       "google",
				ProviderID:     tokenInfo.Sub,
			}
			if err := s.repo.Create(&newUser); err != nil {
				return "", nil, fmt.Errorf("failed to create google user: %w", err)
			}
			user = &newUser
		}
	}

	token, err := utils.GenerateToken(user.ID, "user")
	if err != nil {
		utils.Log.Error("Failed to generate JWT token for google user: " + err.Error())
		return "", nil, err
	}

	return token, user, nil
}

func (s *UserAuthService) UpdateUser(id uint, req dto.UpdateUserRequest, file *multipart.FileHeader) (*userModel.User, error) {
	updates := make(map[string]interface{})
	if req.Username != "" {
		updates["username"] = req.Username
	}

	if file != nil {
		filename := utils.GenerateUniqueFilename(file.Filename)
		path := "uploads/profiles/" + filename
		if err := utils.SaveUploadedFile(file, path); err != nil {
			utils.Log.Error("Failed to save profile picture: " + err.Error())
		} else {
			updates["profile_picture"] = path
		}
	}

	if len(updates) == 0 {
		return nil, errors.New("no updates provided")
	}

	if err := s.repo.Update(id, updates); err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

func (s *UserAuthService) GetUserByID(id uint) (*userModel.User, error) {
	return s.repo.FindByID(id)
}