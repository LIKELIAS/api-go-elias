package repository

import (
	"database/sql"
	"fmt"

	"github.com/elias/api-go-elias/internal/domain"
)

type postgresUserRepo struct {
	db *sql.DB
}

// NewPostgresUserRepository crea el adaptador de repositorio PostgreSQL
func NewPostgresUserRepository(db *sql.DB) domain.UserRepository {
	return &postgresUserRepo{db: db}
}

func (r *postgresUserRepo) Create(user *domain.User) error {
	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, user.Name, user.Email, user.Password).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *postgresUserRepo) FindByID(id uint) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = $1`

	err := r.db.QueryRow(query, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuario con id %d no encontrado", id)
		}
		return nil, err
	}
	return user, nil
}

func (r *postgresUserRepo) FindByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1`

	err := r.db.QueryRow(query, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *postgresUserRepo) FindAll() ([]*domain.User, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM users ORDER BY id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if users == nil {
		users = []*domain.User{}
	}
	return users, nil
}

func (r *postgresUserRepo) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at`

	return r.db.QueryRow(query, user.Name, user.Email, user.ID).Scan(&user.UpdatedAt)
}

func (r *postgresUserRepo) Delete(id uint) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("usuario con id %d no encontrado", id)
	}
	return nil
}
