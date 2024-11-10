package psql

import (
	"context"

	"github.com/abrahammegantoro/to-do-list-be/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	Conn *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		Conn: conn,
	}
}

func (u *UserRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.User, err error) {
	rows, err := u.Conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		user := domain.User{}
		err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.Password,
			&user.Name,
			&user.UpdatedAt,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, user)
	}

	return
}

func (u *UserRepository) Login(ctx context.Context, username string, password string) (res domain.User, err error) {
	query := `SELECT id, username, password, name, updated_at, created_at FROM users WHERE username=$1 AND password=$2`

	list, err := u.fetch(ctx, query, username, password)
	if err != nil {
		return domain.User{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return domain.User{}, domain.ErrNotFound
	}

	return
}

func (u *UserRepository) GetByUsername(ctx context.Context, email string) (res domain.User, err error) {
	query := `SELECT id, username, password, name, updated_at, created_at FROM users WHERE username=$1`

	list, err := u.fetch(ctx, query, email)
	if err != nil {
		return domain.User{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return domain.User{}, domain.ErrCredential
	}

	return
}

func (u *UserRepository) Register(ctx context.Context, user domain.User) (err error) {
	query := `INSERT INTO users (username, password, name, updated_at, created_at) VALUES ($1, $2, $3, $4, $5)`

	_, err = u.Conn.Exec(ctx, query, user.Username, user.Password, user.Name, user.UpdatedAt, user.CreatedAt)
	if err != nil {
		return
	}

	return
}

func (u *UserRepository) GetByID(ctx context.Context, id int64) (res domain.User, err error) {
	query := `SELECT id, username, password, name, updated_at, created_at FROM users WHERE id=$1`

	list, err := u.fetch(ctx, query, id)
	if err != nil {
		return domain.User{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return domain.User{}, domain.ErrNotFound
	}

	return
}
