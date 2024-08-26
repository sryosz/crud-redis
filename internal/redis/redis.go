package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"log/slog"
	"microservice-redis/internal/models"
)

type Redis struct {
	Client *redis.Client
	log    *slog.Logger
}

func New(log *slog.Logger) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	err := client.Ping().Err()
	if err != nil {
		log.Error("Failed to connect to redis", "error", err)
	}

	return &Redis{
		Client: redis.NewClient(&redis.Options{}),
		log:    log,
	}
}

func userIDKey(id int64) string {
	return fmt.Sprintf("user:%d", id)
}

func (r *Redis) Insert(user models.User) error {

	log := r.log.With("op", "redis.Insert")

	data, err := json.Marshal(user)
	if err != nil {
		log.Error("Failed to encode", "error", err)
		return err
	}

	key := userIDKey(user.ID)

	txn := r.Client.TxPipeline()

	res := txn.SetNX(key, string(data), 0)
	if res.Err() != nil {
		txn.Discard()
		log.Error("Failed to set", "error", err)
		return err
	}

	if err = txn.SAdd("orders", key).Err(); err != nil {
		txn.Discard()
		log.Error("Failed to add to set", "error", err)
		return err
	}

	if _, err = txn.Exec(); err != nil {
		log.Error("Failed to exec", "error", err)
		return err
	}

	return nil
}

var ErrNoExist error = errors.New("user does not exist")

func (r *Redis) FindById(id int64) (models.User, error) {
	log := r.log.With("op", "redis.FindById")

	key := userIDKey(id)

	value, err := r.Client.Get(key).Result()
	if errors.Is(err, redis.Nil) {
		log.Error("User does not exist", "error", ErrNoExist)
		return models.User{}, ErrNoExist
	} else if err != nil {
		log.Error("Failed to get user", "error", err)
		return models.User{}, err
	}

	var user models.User
	err = json.Unmarshal([]byte(value), &user)
	if err != nil {
		log.Error("Failed to decode", "error", err)
		return models.User{}, err
	}

	return user, nil
}
