package account

import (
	"context"

	"github.com/segmentio/ksuid"
)

type (
	AccountService interface {
		PostAccount(ctx context.Context, name string) (*Account, error)
		GetAccount(ctx context.Context, id string) (*Account, error)
		GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
	}

	Account struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	accountService struct {
		repository AccountRepository
	}
)

func NewAccountService(r AccountRepository) AccountService {
	return &accountService{r}
}

func (s *accountService) PostAccount(
	ctx context.Context,
	name string,
) (*Account, error) {
	account := &Account{
		ID:   ksuid.New().String(),
		Name: name,
	}

	err := s.repository.PutAccount(ctx, *account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *accountService) GetAccount(
	ctx context.Context,
	id string,
) (*Account, error) {
	return s.repository.GetAccountById(ctx, id)
}

func (s *accountService) GetAccounts(
	ctx context.Context,
	skip uint64,
	take uint64,
) ([]Account, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.ListAccounts(ctx, skip, take)
}
