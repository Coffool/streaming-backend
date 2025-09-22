package services

import (
	"auth-service/models"
	"auth-service/repositories"
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshResponse struct {
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type AuthServiceInterface interface {
	Login(req LoginRequest) (*LoginResponse, error)
	RefreshToken(req RefreshRequest) (*RefreshResponse, error)
	Logout(userID uint) error
}

type AuthService struct {
	userRepo         repositories.UserRepositoryInterface
	refreshTokenRepo repositories.RefreshTokenRepositoryInterface
}

func NewAuthService(userRepo repositories.UserRepositoryInterface, refreshTokenRepo repositories.RefreshTokenRepositoryInterface) AuthServiceInterface {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (s *AuthService) RefreshToken(req RefreshRequest) (*RefreshResponse, error) {
	now := time.Now()

	// Buscar refresh token en DB
	stored, err := s.refreshTokenRepo.FindByToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("no se pudo renovar el token")
	}

	// Verificar expiración
	if now.After(stored.ExpiresAt) {
		return nil, errors.New("no se pudo renovar el token")
	}

	// Buscar usuario asociado
	user, err := s.userRepo.FindByID(stored.UserID)
	if err != nil {
		return nil, errors.New("no se pudo renovar el token")
	}

	// Generar nuevo access token
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "defaultsecret"
	}

	accessTTLStr := os.Getenv("ACCESS_TOKEN_TTL")
	if accessTTLStr == "" {
		accessTTLStr = "15m" // valor por defecto
	}
	accessTTL, err := time.ParseDuration(accessTTLStr)
	if err != nil {
		accessTTL = 15 * time.Minute // fallback
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"exp":      now.Add(accessTTL).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return nil, errors.New("no se pudo generar access token")
	}

	// Rotar refresh token si está cerca de expirar (<20% de TTL)
	refreshTTLStr := os.Getenv("REFRESH_TOKEN_TTL")
	if refreshTTLStr == "" {
		refreshTTLStr = "7d" // valor por defecto
	}
	refreshTTL, err := time.ParseDuration(refreshTTLStr)
	if err != nil {
		refreshTTL = 7 * 24 * time.Hour // fallback
	}

	if stored.ExpiresAt.Sub(now) < refreshTTL/5 {
		stored.Token = uuid.New().String()
		stored.ExpiresAt = now.Add(refreshTTL)
		if err := s.refreshTokenRepo.Update(stored); err != nil {
			return nil, errors.New("error actualizando refresh token")
		}
	}

	return &RefreshResponse{
		Message:      "token renovado exitosamente",
		AccessToken:  accessTokenString,
		RefreshToken: stored.Token,
		ExpiresIn:    int(accessTTL.Seconds()),
	}, nil
}

func (s *AuthService) Logout(userID uint) error {
	return s.refreshTokenRepo.DeleteByUserID(userID)
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Message      string      `json:"message"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int         `json:"expires_in"`
	User         UserSummary `json:"user"`
}

type UserSummary struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func (s *AuthService) Login(req LoginRequest) (*LoginResponse, error) {
	// Buscar usuario por username o email
	user, err := s.userRepo.FindByUsernameOrEmail(req.Identifier, req.Identifier)
	if err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	// Comparar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("contraseña incorrecta")
	}

	// JWT_SECRET
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "defaultsecret"
	}

	// Access Token TTL
	accessTTLStr := os.Getenv("ACCESS_TOKEN_TTL")
	if accessTTLStr == "" {
		accessTTLStr = "15m"
	}
	accessTTL, err := time.ParseDuration(accessTTLStr)
	if err != nil {
		accessTTL = 15 * time.Minute
	}

	// Crear claims para el access token
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"exp":      time.Now().Add(accessTTL).Unix(),
	}

	// Generar access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return nil, errors.New("no se pudo generar access token")
	}

	// Refresh Token TTL
	refreshTTLStr := os.Getenv("REFRESH_TOKEN_TTL")
	if refreshTTLStr == "" {
		refreshTTLStr = "7d"
	}
	refreshTTL, err := time.ParseDuration(refreshTTLStr)
	if err != nil {
		refreshTTL = 7 * 24 * time.Hour
	}

	// Generar refresh token
	refreshToken := uuid.New().String()
	expiry := time.Now().Add(refreshTTL)

	// Limpiar refresh tokens anteriores del usuario
	s.refreshTokenRepo.DeleteByUserID(user.ID)

	// Crear nuevo refresh token
	newRT := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: expiry,
	}

	if err := s.refreshTokenRepo.Create(newRT); err != nil {
		return nil, errors.New("no se pudo guardar refresh token")
	}

	return &LoginResponse{
		Message:      "login exitoso",
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
		ExpiresIn:    int(accessTTL.Seconds()),
		User: UserSummary{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
	}, nil
}
