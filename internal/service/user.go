package service

import (
	"errors"
	"go-api/internal/entity"
	"go-api/internal/repo"
	"go-api/internal/utils"
)

type UserServiceInterface interface {
	GetAll() []entity.UserDTO
	Register(request entity.RegisterRequest) error
	Login(request entity.LoginRequest) (*entity.LoginResult, error)
	Logout(session string) error
}

type userService struct {
	repo repo.UserRepoInterface
}

func NewUserService(repo repo.UserRepoInterface) UserServiceInterface {
	return &userService{repo: repo}
}

func (s *userService) GetAll() []entity.UserDTO {
	data, _ := s.repo.GetAll()
	return data
}

func (s *userService) Register(request entity.RegisterRequest) error {
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return err
	}
	hashedUser := entity.RegisterRequest{
		Email:    request.Email,
		Password: hashedPassword,
	}
	return s.repo.AddUser(hashedUser)
}

func (s *userService) Login(request entity.LoginRequest) (*entity.LoginResult, error) {
	// check if account exist
	storedUser, err := s.repo.GetUserByEmail(request.Email)
	if err != nil {
		return nil, errors.New("user with this email doesn't exist")
	}

	//check if password match
	if !utils.CheckPasswordHash(request.Password, storedUser.Password) {
		return nil, errors.New("invalid Password")
	}

	//generate and store session
	session, _ := utils.GenerateSessionID()
	err = s.repo.StoreSession(session, storedUser.Id)
	if err != nil {
		return nil, errors.New("could not store Session")
	}

	//create result object
	result := entity.LoginResult{
		Session: session,
		User: entity.UserDTO{
			Id:    storedUser.Id,
			Email: storedUser.Email,
		},
	}
	return &result, nil
}

func (s *userService) Logout(session string) error {
	return s.repo.DeleteSession(session)
}
