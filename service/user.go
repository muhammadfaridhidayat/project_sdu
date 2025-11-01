package service

import (
	"errors"
	"project_sdu/model"
	"project_sdu/repository"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(user model.User) (token *string, dbUser model.User, err error)
	Register(user model.User) error
	GetUserByID(id int) (model.User, error)

	// CheckPassLength(pass string) bool
	// CheckPassAlphabet(pass string) bool
	
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository}
}

func (s *userService) Login(user model.User) (token *string, dbUser model.User, err error) {
	dbUser, err = s.userRepository.CheckAvail(user)
	if err != nil {
		return nil, model.User{}, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		return nil, model.User{}, errors.New("wrong email or password")
	}

	// if dbUser.Email != user.Email || dbUser.Password != user.Password {
	// 	return nil, 0, errors.New("wrong email or password")
	// }

	// expirationTime := time.Now().Add(10 * time.Minute)
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := model.Claims{
		UserID: dbUser.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	tokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenJwtString, err := tokenJwt.SignedString(model.JwtKey)
	if err != nil {
		return nil, model.User{}, err
	}

	return &tokenJwtString, dbUser, nil
}

func (s *userService) Register(user model.User) error {
	dbUser, _ := s.userRepository.CheckAvail(user)
	// if err != nil {
	// 	return err
	// }

	if dbUser.Email != "" || dbUser.ID != 0 {
		return errors.New("email already exists")
	}

	err := s.userRepository.Add(user)
	if err != nil {
		return err
	}

	return nil
}

// func (s *userService) GetUser() (model.User, error)

// func (s *userService) CheckPassLength(pass string) bool {
// 	if len(pass) <= 5 {
// 		return true
// 	}

// 	return false
// }

// func (s *userService) CheckPassAlphabet(pass string) bool {
// 	for _, charVariable := range pass {
// 		if (charVariable < 'a' || charVariable > 'z') && (charVariable < 'A' || charVariable > 'Z') {
// 			return false
// 		}
// 	}
// 	return true
// }

func (s *userService) GetUserByID(id int) (model.User, error) {
	user, err := s.userRepository.GetUserByID(id)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
