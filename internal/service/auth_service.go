package service

import (
	"errors"
	"medical-center/internal/models/user"
	"medical-center/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(email, password, name string, role user.Role) (*user.User, error)
	Login(email, password string) (string, error) // returns JWT token
	ValidateToken(tokenString string) (*user.User, error)
}

type authService struct {
	userRepo repository.UserRepository
	jwtKey   []byte
}

type Claims struct {
	UserID uint       `json:"user_id"`
	Role   user.Role  `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo, 
		jwtKey:   []byte("your-secret-key-here"), // In production, use environment variables
	}
}

func (s *authService) Register(email, password, name string, role user.Role) (*user.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}
	
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	
	newUser := &user.User{
		Email:     email,
		Password:  string(hashedPassword),
		Name:      name,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	err = s.userRepo.Create(newUser)
	if err != nil {
		return nil, err
	}
	
	return newUser, nil
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}
	
	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}
	
	// Generate token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

func (s *authService) ValidateToken(tokenString string) (*user.User, error) {
	claims := &Claims{}
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtKey, nil
	})
	
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	
	return user, nil
} 