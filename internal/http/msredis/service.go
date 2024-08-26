package msredis

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"microservice-redis/internal/models"
	"microservice-redis/internal/redis"
	"microservice-redis/internal/storage"
	"microservice-redis/internal/storage/postgres"
)

type Service struct {
	log     *slog.Logger
	Storage *storage.Storage
	Redis   *redis.Redis
}

func NewService(log *slog.Logger) *Service {
	const op = "lib.service.NewService"

	return &Service{
		log:     log,
		Storage: storage.New(log),
		Redis:   redis.New(log),
	}
}

func (s *Service) GetUserById(id int) (*postgres.User, error) {
	const op = "lib.service.GetCustomers"
	log := s.log.With("op", op)

	user, err := s.Redis.FindById(int64(id))
	if errors.Is(err, redis.ErrNoExist) {
		log.Error("User was not found at redis", "error", err)
		u, err := s.Storage.GetUserByID(context.Background(), int32(id))
		if err != nil {
			log.Error("Failed to get user", "error", err)
			return nil, err
		}
		log.Info("Successfully got user from db")
		return &u, nil
	} else if err != nil {
		log.Error("Failed to get user from redis", "error", err)
		return nil, err
	}

	log.Info("Successfully got user from redis")

	return &postgres.User{
		ID:       int32(user.ID),
		Email:    pgtype.Text{String: user.Email, Valid: true},
		Password: user.Password,
	}, nil
}

func (s *Service) CreateUser(email, password string) (*postgres.User, error) {
	const op = "lib.service.GetCustomers"
	log := s.log.With("op", op)

	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to hash password", "error", err)
		return nil, err
	}

	err = s.Storage.CreateUser(context.Background(), postgres.CreateUserParams{
		Email:    pgtype.Text{String: email, Valid: true},
		Password: pass,
	})
	if err != nil {
		log.Error("Failed to save user", "error", err)
		return nil, err
	}

	user, err := s.Storage.GetUserByEmail(context.Background(), pgtype.Text{String: email, Valid: true})
	if err != nil {
		log.Error("Failed to get user by email", "error", err)
		return nil, err
	}

	// return user from db when creating and set fields here
	err = s.Redis.Insert(models.User{
		ID:       int64(user.ID),
		Email:    user.Email.String,
		Password: user.Password,
	})
	if err != nil {
		log.Error("Failed to insert user with redis", "error", err)
		return nil, err
	}

	log.Info("Successfully created user")

	return &user, nil
}
