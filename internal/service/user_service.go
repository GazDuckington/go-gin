package service

import (
	"github.com/GazDuckington/go-gin/internal/models/dto"
	"github.com/GazDuckington/go-gin/internal/models/entity"
	"github.com/GazDuckington/go-gin/internal/repository"
)

type UserService interface {
	GetAll() ([]dto.UserResponse, error)
	GetByID(id uint) (*dto.UserResponse, error)
	Create(req dto.CreateUserRequest) (*dto.UserResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) GetAll() ([]dto.UserResponse, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	out := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		out = append(out, dto.UserResponse{ID: u.ID, Name: u.Name, Email: u.Email})
	}
	return out, nil
}

func (s *userService) GetByID(id uint) (*dto.UserResponse, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, nil
	}
	resp := &dto.UserResponse{ID: u.ID, Name: u.Name, Email: u.Email}
	return resp, nil
}

func (s *userService) Create(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	user := &entity.User{Name: req.Name, Email: req.Email}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return &dto.UserResponse{ID: user.ID, Name: user.Name, Email: user.Email}, nil
}
