package service

import (
	"context"
	"errors"

	"github.com/Autodoc-Technology/interview-templates/template/golang/pkg/models"
	"github.com/Autodoc-Technology/interview-templates/template/golang/pkg/storage"
	validation "github.com/go-ozzo/ozzo-validation"
	"go.uber.org/zap"
)

type IService interface {
	CreateUser(ctx context.Context, name, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, limit, skip int) ([]models.User, error)
}

type service struct {
	logger *zap.Logger
	db     storage.IMongoClient
}

var (
	ErrBadValues = errors.New("bad values")
)

func NewService(l *zap.Logger, d storage.IMongoClient) IService {
	return &service{
		logger: l,
		db:     d,
	}
}

func (s *service) CreateUser(ctx context.Context, name, email string) (*models.User, error) {
	u := models.User{
		Email:     email,
		Name:      name,
		IsDeleted: false,
	}

	if err := validation.ValidateStruct(&u, validation.Field(&u.Email, validation.Required),
		validation.Field(&u.Name, validation.Required)); err != nil {
		return nil, err
	}

	resp, err := s.db.CreateUser(ctx, &u)
	if err != nil {
		s.logger.Error("error CreateUser", zap.Error(err))
		return nil, err
	}

	return resp, nil
}

func (s *service) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	resp, err := s.db.GetUserByID(ctx, id)
	if err != nil {
		s.logger.Error("error GetUserByID", zap.Error(err))
		return nil, err
	}

	return resp, nil
}

func (s *service) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	resp, err := s.db.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error("error GetUserByEmail", zap.Error(err))
		return nil, err
	}

	return resp, nil
}

func (s *service) DeleteUser(ctx context.Context, id string) error {
	err := s.db.DeleteUser(ctx, id)
	if err != nil {
		s.logger.Error("error DeleteUser", zap.Error(err))
		return err
	}

	return nil
}

func (s *service) ListUsers(ctx context.Context, limit, skip int) ([]models.User, error) {
	resp, err := s.db.ListUsers(ctx, limit, skip)
	if err != nil {
		s.logger.Error("error ListUsers", zap.Error(err))
		return nil, err
	}

	return resp, nil
}
