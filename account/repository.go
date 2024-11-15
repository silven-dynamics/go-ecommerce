package account

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type (
	AccountRepository interface {
		Close()
		PutAccount(ctx context.Context, a Account) error
		GetAccountById(ctx context.Context, id string) (*Account, error)
		ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
	}

	postgresRepository struct {
		db *sql.DB
	}
)

func NewPostgresRepository(url string) (AccountRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) Ping() error {
	return r.db.Ping()
}

func (r *postgresRepository) PutAccount(
	ctx context.Context,
	a Account,
) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO accounts(id, name) VALUES($1, $2)", a.ID, a.Name)
	return err
}

func (r *postgresRepository) GetAccountById(
	ctx context.Context,
	id string,
) (*Account, error) {
	a := &Account{}
	row := r.db.QueryRowContext(ctx, "SELECT id, name FROM accounts WHERE id = $1", id)
	if err := row.Scan(&a.ID, &a.Name); err != nil {
		return nil, err
	}

	return a, nil
}

func (r *postgresRepository) ListAccounts(
	ctx context.Context,
	skip uint64,
	take uint64,
) ([]Account, error) {
	accounts := []Account{}

	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, name FROM accounts ORDER BY id DESC OFFSET $1 LIMIT $2",
		skip,
		take,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		a := &Account{}
		if err = rows.Scan(&a.ID, &a.Name); err == nil {
			accounts = append(accounts, *a)
		}
	}

	if rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}
