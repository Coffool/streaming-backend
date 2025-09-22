// services/user_service.go
package services

import (
	"auth-service/models"
	"auth-service/repositories"
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=20"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	Birthdate string `json:"birthdate" binding:"required"`
}

type UpdateRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Birthdate string `json:"birthdate,omitempty"`
}

type UserDetailResponse struct {
	ID                 uint        `json:"id"`
	Name               string      `json:"name"`
	Username           string      `json:"username"`
	Email              string      `json:"email"`
	Role               string      `json:"role"`
	Birthdate          string      `json:"birthdate"`
	RegisterDate       string      `json:"register_date"`
	LastUsernameChange interface{} `json:"last_username_change"`
	LastEmailChange    interface{} `json:"last_email_change"`
	LastPasswordChange interface{} `json:"last_password_change"`
}

type UserServiceInterface interface {
	RegisterUser(req RegisterRequest) (*UserResponse, error)
	UpdateUser(userID uint, req UpdateRequest) (*UserResponse, error)
	GetUserDetails(userID uint) (*UserDetailResponse, error)
}

type UserService struct {
	userRepo repositories.UserRepositoryInterface
}

func NewUserService(userRepo repositories.UserRepositoryInterface) UserServiceInterface {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) RegisterUser(req RegisterRequest) (*UserResponse, error) {
	// Parse birthdate
	birthdate, err := time.Parse("2006-01-02", req.Birthdate)
	if err != nil {
		return nil, errors.New("fecha inválida, formato esperado YYYY-MM-DD")
	}

	// Verificar si username o email ya existen
	_, err = s.userRepo.FindByUsernameOrEmail(req.Username, req.Email)
	if err == nil {
		return nil, errors.New("usuario o email ya registrados")
	}

	// Encriptar contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("no se pudo encriptar la contraseña")
	}

	// Crear usuario
	user := &models.User{
		Name:         "user",
		Role:         "user",
		Username:     req.Username,
		Email:        req.Email,
		Password:     string(hashedPassword),
		Birthdate:    birthdate,
		Registerdate: time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Birthdate: user.Birthdate.Format("2006-01-02"),
	}, nil
}

func (s *UserService) UpdateUser(userID uint, req UpdateRequest) (*UserResponse, error) {
	// Obtener usuario existente
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	// Actualizar campos
	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Username != "" && req.Username != user.Username {
		exists, err := s.userRepo.ExistsByUsername(req.Username)
		if err != nil {
			return nil, errors.New("error al verificar username")
		}
		if exists {
			return nil, errors.New("username ya está en uso")
		}
		user.Username = req.Username
	}

	if req.Email != "" && req.Email != user.Email {
		// Validar formato de email
		matched, _ := regexp.MatchString(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`, req.Email)
		if !matched {
			return nil, errors.New("formato de email inválido")
		}

		exists, err := s.userRepo.ExistsByEmail(req.Email)
		if err != nil {
			return nil, errors.New("error al verificar email")
		}
		if exists {
			return nil, errors.New("email ya está en uso")
		}
		user.Email = req.Email
	}

	if req.Password != "" {
		if len(req.Password) < 6 {
			return nil, errors.New("la contraseña debe tener al menos 6 caracteres")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("no se pudo encriptar la contraseña")
		}
		user.Password = string(hashedPassword)
	}

	// Guardar cambios
	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("no se pudo actualizar el usuario")
	}

	return &UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

func (s *UserService) GetUserDetails(userID uint) (*UserDetailResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	return &UserDetailResponse{
		ID:                 user.ID,
		Name:               user.Name,
		Username:           user.Username,
		Email:              user.Email,
		Role:               user.Role,
		Birthdate:          user.Birthdate.Format("2006-01-02"),
		RegisterDate:       user.Registerdate.Format("2006-01-02 15:04:05"),
		LastUsernameChange: formatTimePointer(user.LastUsernameChange),
		LastEmailChange:    formatTimePointer(user.LastEmailChange),
		LastPasswordChange: formatTimePointer(user.LastPasswordChange),
	}, nil
}

// formatTimePointer formatea un puntero a time.Time para JSON
func formatTimePointer(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return t.Format("2006-01-02 15:04:05")
}
