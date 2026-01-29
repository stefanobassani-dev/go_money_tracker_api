package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT user_id, email, password, tink_user_id,
       created_at, updated_at
		FROM users WHERE email = $1`

	var user domain.User

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.TinkUserId,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING user_id, created_at, updated_at"

	err := r.db.QueryRow(ctx, query, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

func (r *UserRepository) GetTinkIDByUserID(ctx context.Context, userID string) (string, error) {
	var tinkID sql.NullString

	query := `SELECT tink_user_id FROM users WHERE user_id = $1`

	err := r.db.QueryRow(ctx, query, userID).Scan(&tinkID)
	if err != nil {
		return "", err
	}

	if tinkID.Valid {
		return tinkID.String, nil
	}

	return "", nil
}

func (r *UserRepository) UpdateTinkID(ctx context.Context, userID string, newTinkID string) error {
	query := `UPDATE users SET tink_user_id = $2 WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID, newTinkID)
	return err
}
