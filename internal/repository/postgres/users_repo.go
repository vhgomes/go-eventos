package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"vhgomes-eventos/internal/domain"
	erro "vhgomes-eventos/internal/pkg/errors"
	"vhgomes-eventos/internal/repository"
)

var _ repository.UserRepository = (*userRepo)(nil)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) repository.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, user.Email, user.Password, user.Name).Scan(&user.Id)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id int) (*domain.User, error) {
	query := `SELECT id, email, name, password FROM users WHERE id = $1`
	var user domain.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.Email, &user.Name, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, erro.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, name, password FROM users WHERE email = $1`
	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.Id, &user.Email, &user.Name, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, erro.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

func (r *userRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return exists, nil
}
