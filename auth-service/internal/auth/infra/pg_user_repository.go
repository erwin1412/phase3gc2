package infra

import (
	"auth-service/internal/auth/domain"
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
)

type pgUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *pgUserRepository {
	if db == nil {
		log.Fatal("Postgres DB is nil")
	}
	return &pgUserRepository{db}
}

func (r *pgUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1 LIMIT 1`,
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *pgUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	id := uuid.NewString()
	if id == "" {
		return nil, errors.New("failed to generate UUID")
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)`,
		id, user.Name, user.Email, user.Password,
	)
	if err != nil {
		return nil, err
	}

	user.ID = id
	return user, nil
}
