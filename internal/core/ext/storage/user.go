package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/compiai/engine/internal/core/domain/user"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"strings"
)

// PostgresStorage implements Storage using a PostgreSQL database.
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage creates a new PostgresStorage.
func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) Save(ctx context.Context, usr user.User) error {
	query := `
	INSERT INTO users (id, username, solana_wallet_public_key, password_hash, games)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (id) DO UPDATE
	SET username = EXCLUDED.username,
	    solana_wallet_public_key = EXCLUDED.solana_wallet_public_key,
	    password_hash = EXCLUDED.password_hash,
	    games = EXCLUDED.games
	`
	_, err := s.db.ExecContext(ctx, query,
		usr.ID, usr.Username, usr.SolanaWalletPublicKey, usr.PasswordHash, pq.Array(usr.Games),
	)
	return err
}

func (s *PostgresStorage) FindOne(ctx context.Context, filter user.SingleFilter) (user.User, error) {
	// delegate to specific methods
	if filter.ID != nil {
		return s.FindOneByID(ctx, *filter.ID)
	}
	if filter.Username != nil {
		return s.FindOneByUsername(ctx, *filter.Username)
	}
	if filter.SolanaWallet != nil {
		return s.FindOneBySolanaWallet(ctx, *filter.SolanaWallet)
	}
	return user.User{}, sql.ErrNoRows
}

func (s *PostgresStorage) FindOneByUsername(ctx context.Context, username string) (user.User, error) {
	var u user.User
	query := `
	SELECT id, username, solana_wallet_public_key, password_hash, games
	FROM users WHERE username = $1
	`
	err := s.db.QueryRowContext(ctx, query, username).Scan(
		&u.ID, &u.Username, &u.SolanaWalletPublicKey, &u.PasswordHash, pq.Array(&u.Games),
	)
	return u, err
}

func (s *PostgresStorage) FindOneByID(ctx context.Context, id uuid.UUID) (user.User, error) {
	var u user.User
	query := `
	SELECT id, username, solana_wallet_public_key, password_hash, games
	FROM users WHERE id = $1
	`
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Username, &u.SolanaWalletPublicKey, &u.PasswordHash, pq.Array(&u.Games),
	)
	return u, err
}

func (s *PostgresStorage) FindOneBySolanaWallet(ctx context.Context, solanaWallet string) (user.User, error) {
	var u user.User
	query := `
	SELECT id, username, solana_wallet_public_key, password_hash, games
	FROM users WHERE solana_wallet_public_key = $1
	`
	err := s.db.QueryRowContext(ctx, query, solanaWallet).Scan(
		&u.ID, &u.Username, &u.SolanaWalletPublicKey, &u.PasswordHash, pq.Array(&u.Games),
	)
	return u, err
}

func (s *PostgresStorage) Find(ctx context.Context, filter user.Filter) ([]user.User, error) {
	// build dynamic IN clauses
	clauses := []string{}
	args := []interface{}{}
	idx := 1
	if len(filter.IDs) > 0 {
		clauses = append(clauses, fmt.Sprintf("id = ANY($%d)", idx))
		args = append(args, pq.Array(filter.IDs))
		idx++
	}
	if len(filter.Usernames) > 0 {
		clauses = append(clauses, fmt.Sprintf("username = ANY($%d)", idx))
		args = append(args, pq.Array(filter.Usernames))
		idx++
	}
	if len(filter.SolanaWallets) > 0 {
		clauses = append(clauses, fmt.Sprintf("solana_wallet_public_key = ANY($%d)", idx))
		args = append(args, pq.Array(filter.SolanaWallets))
		idx++
	}
	if len(clauses) == 0 {
		return nil, nil
	}
	query := "SELECT id, username, solana_wallet_public_key, password_hash, games FROM users WHERE " + strings.Join(clauses, " AND ")
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []user.User{}
	for rows.Next() {
		var u user.User
		err := rows.Scan(&u.ID, &u.Username, &u.SolanaWalletPublicKey, &u.PasswordHash, pq.Array(&u.Games))
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}
