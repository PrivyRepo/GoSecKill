package services

import (
	"golang.org/x/crypto/bcrypt"
	"homework/models/datamodels"
	"homework/models/repositories"
)

type IUserService interface {
	IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool)
	AddUser(user *datamodels.User) (userId int64, err error)
}

type UserService struct {
	UserRepository repositories.IUserRepository
}

func NewUserService(repository repositories.IUserRepository) IUserService {
	return &UserService{repository}
}

func (u *UserService) IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOK bool) {
	user, err := u.UserRepository.Select(userName)
	if err != nil {
		return
	}
	if !ValidatePassword(pwd, user.HashPassword) {
		return &datamodels.User{}, false
	}
	return user, true
}

func (u *UserService) AddUser(user *datamodels.User) (userId int64, err error) {
	pwdByte, errPwd := GeneratePassword(user.HashPassword)
	if errPwd != nil {
		return userId, errPwd
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.Insert(user)
}

func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

func ValidatePassword(userPassword string, hashed string) (isOK bool) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err == nil {
		return true
	}
	return
}
