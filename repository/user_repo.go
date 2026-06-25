package repository

import(
	"context"
	"database/sql"
)

type User struct {
	ID 		int64
	Name 	string
	Email 	string
	Role 	string
}

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	
}

type postgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
	query := "SELECT id, name, email, role FROM users WHERE id = $1"

	var user User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Role)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}