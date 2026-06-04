package repository

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"student-cert-service/internal/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(dbURL string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	repo := &PostgresRepository{db: db}
	if err := repo.migrate(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return repo, nil
}

func (r *PostgresRepository) migrate() error {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL
	);`

	requestsTable := `
	CREATE TABLE IF NOT EXISTS requests (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id),
		type VARCHAR(255) NOT NULL,
		status VARCHAR(50) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`

	if _, err := r.db.Exec(usersTable); err != nil {
		return err
	}
	if _, err := r.db.Exec(requestsTable); err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) CreateUser(username, password string) (int, error) {
	var id int
	err := r.db.QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", username, password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PostgresRepository) GetUserByUsername(username string) (*domain.User, error) {
	u := &domain.User{}
	err := r.db.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).
		Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func (r *PostgresRepository) CreateRequest(userID int, reqType string) (*domain.Request, error) {
	req := &domain.Request{
		UserID:    userID,
		Type:      reqType,
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	err := r.db.QueryRow(
		"INSERT INTO requests (user_id, type, status, created_at) VALUES ($1, $2, $3, $4) RETURNING id",
		req.UserID, req.Type, req.Status, req.CreatedAt,
	).Scan(&req.ID)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *PostgresRepository) GetRequestsByUserID(userID int) ([]domain.Request, error) {
	rows, err := r.db.Query("SELECT id, user_id, type, status, created_at FROM requests WHERE user_id = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []domain.Request
	for rows.Next() {
		var req domain.Request
		if err := rows.Scan(&req.ID, &req.UserID, &req.Type, &req.Status, &req.CreatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func (r *PostgresRepository) UpdateRequestStatus(requestID int, status string) error {
	_, err := r.db.Exec("UPDATE requests SET status = $1 WHERE id = $2", status, requestID)
	return err
}
