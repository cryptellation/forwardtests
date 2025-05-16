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

// GetStatus gets the status of the forwardtest.
func (ft Forwardtest) GetStatus(
	ctx context.Context,
) (forwardtest.Status, error) {
	res, err := ft.rawClient.GetForwardtestStatus(ctx, api.GetForwardtestStatusWorkflowParams{
		ForwardtestID: ft.ID,
	})
	if err != nil {
		return forwardtest.Status{}, err
	}

	return res.Status, nil
}
