package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Gwinkamp/grpcauth-sso/internal/domain/models"
	"github.com/Gwinkamp/grpcauth-sso/internal/storage"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// New инициализирует новое хранилище данных
func New(storagePath string) (*Storage, error) {
	const operation = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &Storage{db: db}, nil
}

// CreateUser создает нового пользователя в хранилище
func (s *Storage) CreateUser(ctx context.Context, email string, passHash []byte) (uuid.UUID, error) {
	const operation = "storage.sqlite.CreateUser"

	stmt, err := s.db.Prepare("INSERT INTO users(id, email, pass_hash) VALUES(?, ?, ?)")
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", operation, err)
	}

	id := uuid.New()

	_, err = stmt.ExecContext(ctx, id.String(), email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return uuid.Nil, fmt.Errorf("%s: %w", operation, storage.ErrUserAlreadyExists)
		}

		return uuid.Nil, fmt.Errorf("%s: %w", operation, err)
	}

	return id, nil
}

// GetUser возвращает данные пользователя по его email
func (s *Storage) GetUser(ctx context.Context, email string) (models.User, error) {
	const operation = "strorage.sqlite.GetUser"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", operation, err)
	}

	row := stmt.QueryRowContext(ctx, email)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", operation, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", operation, err)
	}

	return user, nil
}

// IsAdmin определяет, является ли пользователь администром
func (s *Storage) IsAdmin(ctx context.Context, userId uuid.UUID) (bool, error) {
	const operation = "storage.sqlite.IsAdmin"

	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", operation, err)
	}

	row := stmt.QueryRowContext(ctx, userId.String())

	var isAdmin bool

	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", operation, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", operation, err)
	}

	return isAdmin, nil
}

// GetService получает информацию о сервисе по его ID
func (s *Storage) GetService(ctx context.Context, serviceId uuid.UUID) (models.Service, error) {
	const operation = "storage.sqlite.GetService"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM services WHERE id = ?")
	if err != nil {
		return models.Service{}, fmt.Errorf("%s: %w", operation, err)
	}

	row := stmt.QueryRowContext(ctx, serviceId.String())

	var service models.Service

	err = row.Scan(&service.ID, &service.Name, &service.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Service{}, fmt.Errorf("%s: %w", operation, storage.ErrServiceNotFound)
		}

		return models.Service{}, fmt.Errorf("%s: %w", operation, err)
	}

	return service, nil
}

// DeleteUser удаляет пользователя из хранилища
func (s *Storage) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	const operation = "storage.sqlite.DeleteUser"

	stmt, err := s.db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	_, err = stmt.ExecContext(ctx, userId.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", operation, storage.ErrUserNotFound)
		}

		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}
