package useraccess

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"core/internal/domain"
)

type UserAccess struct {
	pool *pgxpool.Pool
}

func NewUserAccess(pool *pgxpool.Pool) *UserAccess {
	return &UserAccess{pool: pool}
}

func (a *UserAccess) Create(ctx context.Context, input domain.CreateProfileInput) (domain.Profile, error) {
	const query = `
		INSERT INTO users (email, display_name, bio)
		VALUES ($1, $2, $3)
		RETURNING id, email, display_name, bio, created_at, updated_at
	`

	profile, err := scanProfile(a.pool.QueryRow(ctx, query, input.Email, input.DisplayName, input.Bio))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.Profile{}, domain.ErrEmailAlreadyTaken
		}
		return domain.Profile{}, fmt.Errorf("tworzenie użytkownika: %w", err)
	}

	return profile, nil
}

func (a *UserAccess) GetByID(ctx context.Context, id int64) (domain.Profile, error) {
	const query = `
		SELECT id, email, display_name, bio, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	profile, err := scanProfile(a.pool.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Profile{}, domain.ErrNotFound
		}
		return domain.Profile{}, fmt.Errorf("pobieranie użytkownika: %w", err)
	}

	return profile, nil
}

func (a *UserAccess) Update(ctx context.Context, id int64, input domain.UpdateProfileInput) (domain.Profile, error) {
	const query = `
		UPDATE users
		SET display_name = CASE WHEN $2 THEN $3 ELSE display_name END,
		    bio = CASE WHEN $4 THEN $5 ELSE bio END,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, email, display_name, bio, created_at, updated_at
	`

	var displayName string
	if input.DisplayName != nil {
		displayName = *input.DisplayName
	}

	var bio string
	if input.Bio != nil {
		bio = *input.Bio
	}

	profile, err := scanProfile(a.pool.QueryRow(
		ctx,
		query,
		id,
		input.DisplayName != nil,
		displayName,
		input.Bio != nil,
		bio,
	))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Profile{}, domain.ErrNotFound
		}
		return domain.Profile{}, fmt.Errorf("aktualizacja użytkownika: %w", err)
	}

	return profile, nil
}

func scanProfile(row pgx.Row) (domain.Profile, error) {
	var profile domain.Profile
	err := row.Scan(
		&profile.ID,
		&profile.Email,
		&profile.DisplayName,
		&profile.Bio,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	return profile, err
}
