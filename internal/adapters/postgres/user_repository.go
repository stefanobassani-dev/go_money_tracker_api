package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

var (
	ErrPGXUniqueViolationCode = "23505"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

type entity struct {
	ID         string    `db:"user_id"`
	Email      string    `db:"email"`
	Password   string    `db:"password"`
	TinkUserId *string   `db:"tink_user_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, string, error) {
	query := `SELECT user_id, email, password, tink_user_id,
       created_at, updated_at
		FROM users WHERE email = $1`

	var user entity

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", domain.ErrUserNotFound
		}
		slog.Error("database error during user retrieval", "error", err, "email", email)
		return nil, "", domain.ErrInternal
	}
	return entityToDomain(user), user.Password, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, email, password string) error {
	//TODO return user with info for frontend
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING user_id, created_at, updated_at"

	_, err := r.db.ExecContext(ctx, query, email, password)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == ErrPGXUniqueViolationCode {
			return domain.ErrUserAlreadyExists
		}

		slog.Error("database error during user creation", "error", err, "email", email)
		return domain.ErrInternal
	}
	return nil
}

func (r *UserRepository) GetTinkIDByUserID(ctx context.Context, userID string) (string, error) {
	var tinkID sql.NullString

	query := `SELECT tink_user_id FROM users WHERE user_id = $1`

	err := r.db.GetContext(ctx, &tinkID, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrUserNotFound
		}
		return "", domain.ErrInternal
	}

	if tinkID.Valid {
		return tinkID.String, nil
	}

	return "", nil
}

func (r *UserRepository) UpdateTinkID(ctx context.Context, userID string, newTinkID string) error {
	query := `UPDATE users SET tink_user_id = $2 WHERE user_id = $1`
	res, err := r.db.ExecContext(ctx, query, userID, newTinkID)
	if err != nil {
		slog.Error("database error during user update", "error", err, "user_id", userID)
		return domain.ErrInternal
	}

	count, err := res.RowsAffected()
	if err != nil {
		slog.Error("error checking rows affected during user update", "error", err)
		return domain.ErrInternal
	}

	if count == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func entityToDomain(user entity) *domain.User {
	return &domain.User{
		ID:             user.ID,
		Email:          user.Email,
		ProviderUserID: toString(user.TinkUserId),
	}
}

func toString(ptr *string) string {
	if ptr != nil {
		return *ptr
	}

	return ""
}
