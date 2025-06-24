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
	ID     uuid.UUID
	client RawClient
}

// Run runs the forwardtest with the given bot.
func (ft *Forwardtest) Run(ctx context.Context) error {
	// Start forwardtest
	_, err := ft.client.StartForwardtest(ctx, api.StartForwardtestWorkflowParams{
		ForwardtestID: ft.ID,
	})

	return err
}

// CreateOrder creates an order on the forwardtest.
func (ft Forwardtest) CreateOrder(
	ctx context.Context,
	order order.Order,
) (api.CreateForwardtestOrderWorkflowResults, error) {
	return ft.client.CreateForwardtestOrder(ctx, api.CreateForwardtestOrderWorkflowParams{
		ForwardtestID: ft.ID,
		Order:         order,
	})
}

// ListAccounts lists the accounts of the forwardtest.
func (ft Forwardtest) ListAccounts(
	ctx context.Context,
) (map[string]account.Account, error) {
	res, err := ft.client.ListForwardtestAccounts(ctx, api.ListForwardtestAccountsWorkflowParams{
		ForwardtestID: ft.ID,
	})
	if err != nil {
		return nil, err
	}

	return res.Accounts, nil
}

// GetStatus gets the status of the forwardtest.
func (ft Forwardtest) GetStatus(
	ctx context.Context,
) (forwardtest.Status, error) {
	res, err := ft.client.GetForwardtestStatus(ctx, api.GetForwardtestStatusWorkflowParams{
		ForwardtestID: ft.ID,
	})
	if err != nil {
		return forwardtest.Status{}, err
	}

	return res.Status, nil
}

// Stop stops the forwardtest by executing the exit callback.
func (ft Forwardtest) Stop(ctx context.Context) error {
	_, err := ft.client.StopForwardtest(ctx, api.StopForwardtestWorkflowParams{
		ForwardtestID: ft.ID,
	})

	return err
}
