package service

import (
	"api-payment/entity"
	"api-payment/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	VeryfiyCredential(email string, password string) interface{}
	IsDuplicateEmail(email string) bool
	CreateUser(user entity.User) entity.User
}

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepository: userRepo,
	}
}

func (service *authService) VeryfiyCredential(email string, password string) interface{} {
	res := service.userRepository.VerifyCredential(email, password)
	if v, ok := res.(entity.User); ok {
		comparedPassword := comparePassword(v.Password, password)
		if v.Email == email && comparedPassword {
			return res
		}
		return false
	}
	return false
}

func (service *authService) IsDuplicateEmail(email string) bool {
	res := service.userRepository.IsDuplicateEmail(email)
	if res.Error != nil {
		return false
	}
	return true
}

func (service *authService) CreateUser(user entity.User) entity.User {
	userToCreate := entity.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	} 
	return service.userRepository.InsertUser(userToCreate)
}

// fungsi compare password
func comparePassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}