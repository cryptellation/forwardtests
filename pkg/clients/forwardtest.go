package clients

import (
	"context"

	"github.com/cryptellation/forwardtests/api"
	"github.com/cryptellation/forwardtests/pkg/forwardtest"
	"github.com/cryptellation/runtime/account"
	"github.com/cryptellation/runtime/order"
	"github.com/google/uuid"
)

// Forwardtest is a local representation of a forwardtest running on the Cryptellation API.
type Forwardtest struct {
	ID        uuid.UUID
	rawClient RawClient
}

// Run runs the forwardtest with the given bot.
func (ft *Forwardtest) Run(ctx context.Context) error {
	// Run forwardtest
	_, err := ft.rawClient.RunForwardtest(ctx, api.RunForwardtestWorkflowParams{
		ForwardtestID: ft.ID,
	})

	return err
}

// CreateOrder creates an order on the forwardtest.
func (ft Forwardtest) CreateOrder(
	ctx context.Context,
	order order.Order,
) (api.CreateForwardtestOrderWorkflowResults, error) {
	return ft.rawClient.CreateForwardtestOrder(ctx, api.CreateForwardtestOrderWorkflowParams{
		ForwardtestID: ft.ID,
		Order:         order,
	})
}

// ListAccounts lists the accounts of the forwardtest.
func (ft Forwardtest) ListAccounts(
	ctx context.Context,
) (map[string]account.Account, error) {
	res, err := ft.rawClient.ListForwardtestAccounts(ctx, api.ListForwardtestAccountsWorkflowParams{
		ForwardtestID: ft.ID,
	})
	if err != nil {
		return nil, err
	}

	return res.Accounts, nil
}

// Get retrieves the forwardtest data from the database.
func (ft Forwardtest) Get(ctx context.Context) (forwardtest.Forwardtest, error) {
	res, err := ft.rawClient.GetForwardtest(ctx, api.GetForwardtestWorkflowParams{
		ForwardtestID: ft.ID,
	})
	if err != nil {
		return forwardtest.Forwardtest{}, err
	}

	return res.Forwardtest, nil
}

// GetBalance gets the balance of the forwardtest.
func (ft Forwardtest) GetBalance(
	ctx context.Context,
) (float64, error) {
	res, err := ft.rawClient.GetForwardtestBalance(ctx, api.GetForwardtestBalanceWorkflowParams{
		ForwardtestID: ft.ID,
	})
	if err != nil {
		return 0, err
	}

	return res.Balance, nil
}

// Stop stops the forwardtest by executing the exit callback.
func (ft Forwardtest) Stop(ctx context.Context) error {
	_, err := ft.rawClient.StopForwardtest(ctx, api.StopForwardtestWorkflowParams{
		ForwardtestID: ft.ID,
	})

	return err
}
