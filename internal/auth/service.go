package auth

import (
	"errors"
	"expense-tracker/internal/user"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	db         *gorm.DB
	jwtService *JWTService
}

func NewService(db *gorm.DB, jwtService *JWTService) *Service {
	return &Service{
		db:         db,
		jwtService: jwtService,
	}
}
func (s *Service) Register(req RegisterRequest) (*AuthResponse, error) {
	var existingUser user.User

	err := s.db.Where("email= ?", req.Email).First(&existingUser).Error

	if err == nil {
		return nil, errors.New("user already exist")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	newUser := user.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	err = s.db.Create(&newUser).Error
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwtService.GenerateAccessToken(newUser.Id, newUser.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(newUser.Id)

	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (s *Service) Login(req LoginRequest) (*AuthResponse, error) {
	var existinguser user.User

	err := s.db.Where("email= ?", req.Email).First(&existinguser).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(existinguser.Password), []byte([]byte(req.Password)))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := s.jwtService.GenerateAccessToken(existinguser.Id, existinguser.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(existinguser.Id)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshToken(req RefreshTokenRequest) (*AuthResponse, error) {
	_, claims, err := s.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID := uint(userIDFloat)

	var existingUser user.User
	if err := s.db.First(&existingUser, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	accessToken, err := s.jwtService.GenerateAccessToken(existingUser.Id, existingUser.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(existingUser.Id)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
