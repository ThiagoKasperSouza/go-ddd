package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	domainIdentity "go-ddd/internal/domain/identity"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(ctx context.Context, user *domainIdentity.User) error {
	query := `
		INSERT INTO users (id, name, email, password_hash, roles, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	var rolesStr []string
	for _, role := range user.Roles() {
		rolesStr = append(rolesStr, string(role))
	}

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID(),
		user.Name(),
		user.Email(),
		user.PasswordHash(),
		pq.Array(rolesStr),
		user.IsActive(),
		user.CreatedAt(),
	)

	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domainIdentity.User, error) {
	query := `SELECT id, name, email, password_hash, roles, is_active, created_at FROM users WHERE email = $1`

	var id, name, userEmail, passwordHash string
	var rolesRaw []string
	var isActive bool
	var createdAt time.Time

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&id,
		&name,
		&userEmail,
		&passwordHash,
		pq.Array(&rolesRaw),
		&isActive,
		&createdAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	var roles []domainIdentity.Role
	for _, rStr := range rolesRaw {
		roles = append(roles, domainIdentity.Role(rStr))
	}

	return domainIdentity.NewUser(id, name, userEmail, passwordHash, roles)
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domainIdentity.User, error) {
	query := `SELECT id, name, email, password_hash, roles, is_active, created_at FROM users WHERE id = $1`

	var name, userEmail, passwordHash string
	var rolesRaw []string
	var isActive bool
	var createdAt time.Time

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&id,
		&name,
		&userEmail,
		&passwordHash,
		pq.Array(&rolesRaw),
		&isActive,
		&createdAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	var roles []domainIdentity.Role
	for _, rStr := range rolesRaw {
		roles = append(roles, domainIdentity.Role(rStr))
	}

	return domainIdentity.NewUser(id, name, userEmail, passwordHash, roles)
}