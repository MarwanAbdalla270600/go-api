package repo

import (
	"context"
	"go-api/internal/entity"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type UserRepoInterface interface {
	GetAll() ([]entity.UserDTO, error)
	AddUser(request entity.RegisterRequest) error
	GetUserById(id int) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	StoreSession(session string, userId int) error
	DeleteSession(session string) error
	CheckSession(session string) bool
}

type userRepo struct {
	db    *sqlx.DB
	redis *redis.Client
	ctx   context.Context
}

func NewUserRepo(db *sqlx.DB, redis *redis.Client, ctx context.Context) UserRepoInterface {
	return &userRepo{db: db, redis: redis, ctx: ctx}
}

func (r *userRepo) AddUser(request entity.RegisterRequest) error {
	_, err := r.db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", request.Email, request.Password)
	return err
}

func (r *userRepo) GetAll() ([]entity.UserDTO, error) {
	var users []entity.UserDTO
	err := r.db.Select(&users, "SELECT id, email FROM users")
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepo) GetUserById(id int) (*entity.User, error) {
	var user entity.User
	err := r.db.Get(&user, "SELECT id, email, password FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Get(&user, "SELECT id, email, password FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) StoreSession(session string, userId int) error {
	_, err := r.redis.Set(r.ctx, session, userId, 60*time.Minute).Result()
	return err
}

func (r *userRepo) DeleteSession(session string) error {
	err := r.redis.Del(r.ctx, session).Err()
	return err
}

func (r *userRepo) CheckSession(session string) bool {
	_, err := r.redis.Get(r.ctx, session).Result()
	return err == nil
}
