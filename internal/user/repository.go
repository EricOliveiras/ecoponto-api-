package user

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// Create insere um novo usuário no banco
func (r *Repository) Create(ctx context.Context, u User) error {
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	return r.db.QueryRowxContext(ctx, query, u.Email, u.PasswordHash).Scan(&u.ID, &u.CreatedAt)
}

// FindByEmail busca um usuário pelo email
func (r *Repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	query := "SELECT * FROM users WHERE email = $1"

	err := r.db.GetContext(ctx, &u, query, email)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
