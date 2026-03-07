package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/internal/shared/helper"
	"mqfm-backend/internal/shared/logger"
	"mqfm-backend/internal/shared/security"
)

type userAuthService struct {
	repo port.UserRepository
}

func NewUserAuthService(repo port.UserRepository) port.UserAuthService {
	return &userAuthService{repo: repo}
}

func (s *userAuthService) Register(req request.UserRegisterRequest, file *multipart.FileHeader) (*entity.User, error) {
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		logger.Error("failed to hash user password")
		return nil, err
	}

	var profilePath string
	if file != nil {
		filename := helper.GenerateUniqueFilename(file.Filename)
		path := "uploads/profiles/" + filename
		if err := helper.SaveUploadedFile(file, path); err != nil {
			logger.Error("failed to save profile picture")
		} else {
			profilePath = path
		}
	}

	user := entity.User{
		Username:       req.Username,
		Email:          req.Email,
		Password:       hashedPassword,
		ProfilePicture: profilePath,
		Role:           constant.RoleUser,
		Provider:       "local",
	}

	if err := s.repo.Create(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userAuthService) Login(req request.UserLoginRequest) (string, *entity.User, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		logger.Warn("user login: email not found")
		return "", nil, errors.New("invalid user credentials")
	}

	if !security.CheckPassword(req.Password, user.Password) {
		logger.Warn("user login: incorrect password")
		return "", nil, errors.New("invalid user credentials")
	}

	token, err := security.GenerateToken(user.ID, constant.RoleUser)
	if err != nil {
		logger.Error("failed to generate user token")
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

func (s *userAuthService) GoogleLogin(req request.GoogleLoginRequest) (string, *entity.User, error) {
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

	var localProfilePic string
	if tokenInfo.Picture != "" {
		downloadedPath, downloadErr := helper.DownloadImage(tokenInfo.Picture, "uploads/profiles")
		if downloadErr == nil && downloadedPath != "" {
			localProfilePic = downloadedPath
		}
	}

	username := tokenInfo.Name
	if username == "" {
		username = strings.Split(tokenInfo.Email, "@")[0]
	}

	user, err := s.repo.FindByProviderID("google", tokenInfo.Sub)
	if err != nil {
		existingUser, emailErr := s.repo.FindByEmail(tokenInfo.Email)
		if emailErr == nil {
			updates := map[string]interface{}{
				"provider":    "google",
				"provider_id": tokenInfo.Sub,
				"username":    username,
			}
			if localProfilePic != "" {
				updates["profile_picture"] = localProfilePic
			}
			_ = s.repo.Update(existingUser.ID, updates)
			existingUser.Provider = "google"
			existingUser.ProviderID = tokenInfo.Sub
			existingUser.Username = username
			if localProfilePic != "" {
				existingUser.ProfilePicture = localProfilePic
			}
			user = existingUser
		} else {
			newUser := entity.User{
				Username:       username,
				Email:          tokenInfo.Email,
				ProfilePicture: localProfilePic,
				Role:           constant.RoleUser,
				Provider:       "google",
				ProviderID:     tokenInfo.Sub,
			}
			if err := s.repo.Create(&newUser); err != nil {
				return "", nil, fmt.Errorf("failed to create google user: %w", err)
			}
			user = &newUser
		}
	} else {
		updates := map[string]interface{}{"username": username}
		if localProfilePic != "" {
			updates["profile_picture"] = localProfilePic
			user.ProfilePicture = localProfilePic
		}
		user.Username = username
		_ = s.repo.Update(user.ID, updates)
	}

	token, err := security.GenerateToken(user.ID, constant.RoleUser)
	if err != nil {
		logger.Error("failed to generate google user token")
		return "", nil, err
	}

	return token, user, nil
}

func (s *userAuthService) UpdateUser(id uint, req request.UpdateUserRequest, file *multipart.FileHeader) (*entity.User, error) {
	updates := make(map[string]interface{})
	if req.Username != "" {
		updates["username"] = req.Username
	}

	if file != nil {
		filename := helper.GenerateUniqueFilename(file.Filename)
		path := "uploads/profiles/" + filename
		if err := helper.SaveUploadedFile(file, path); err != nil {
			logger.Error("failed to save profile picture")
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

func (s *userAuthService) GetUserByID(id uint) (*entity.User, error) {
	return s.repo.FindByID(id)
}
