package user

import (
	"errors"
	"mime/multipart"

	"golang.org/x/crypto/bcrypt"

	"mqfm-backend/internal/dto/auth"
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
		// Logic save file sebenarnya butuh context gin, 
		// tapi idealnya service tidak bergantung Gin.
		// Untuk simplicity MVC standard di Go (tanpa clean architecture strict banget),
		// Service bisa terima multipart.FileHeader tapi butuh helper save.
		// Kita akan buat utils.SaveFile di step selanjutnya utk decouple dari Gin Context jika perlu,
		// TAPI karena Gin Context punya SaveUploadedFile yang mudah, 
		// biasnaya di-pass filenya atau logic save tetap di controller wrapper.
		// REQ User: "Logic simpan file idealnya digeser ke Service".
		// Maka kita butuh cara save file manual di sini.
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