//go:build e2e
// +build e2e

package test

import (
	"context"

	"github.com/cryptellation/forwardtests/api"
	"github.com/cryptellation/forwardtests/pkg/forwardtest"
	"github.com/cryptellation/runtime"
	"github.com/cryptellation/runtime/account"
	"github.com/cryptellation/runtime/order"
)

// createTestCallbacks creates test callbacks for testing
func createTestCallbacks() runtime.Callbacks {
	return runtime.Callbacks{
		OnInitCallback: runtime.CallbackWorkflow{
			Name:          "test-init-workflow",
			TaskQueueName: "test-queue",
		},
		OnNewPricesCallback: runtime.CallbackWorkflow{
			Name:          "test-prices-workflow",
			TaskQueueName: "test-queue",
		},
		OnExitCallback: runtime.CallbackWorkflow{
			Name:          "test-exit-workflow",
			TaskQueueName: "test-queue",
		},
	}
}

func (suite *EndToEndSuite) TestGetForwardtestBalance() {
	// GIVEN a forwardtest

	params := api.CreateForwardtestWorkflowParams{
		Accounts: map[string]account.Account{
			"binance": {
				Balances: map[string]float64{
					"USDT": 1000,
				},
			},
		},
		Callbacks: createTestCallbacks(),
	}
	ft, err := suite.client.NewForwardtest(context.Background(), params)
	suite.Require().NoError(err)

	// WHEN getting the forwardtest balance

	balance, err := ft.GetBalance(context.Background())
	suite.Require().NoError(err)

	// THEN the balance is correct

	suite.Require().Equal(1000.0, balance)
}

func (suite *EndToEndSuite) TestListForwardtestStatus() {
	// GIVEN 3 forwardtests

	params := api.CreateForwardtestWorkflowParams{
		Accounts: map[string]account.Account{
			"binance": {
				Balances: map[string]float64{
					"USDT": 1000,
				},
			},
		},
		Callbacks: createTestCallbacks(),
	}
	ft1, err := suite.client.NewForwardtest(context.Background(), params)
	suite.Require().NoError(err)
	ft2, err := suite.client.NewForwardtest(context.Background(), params)
	suite.Require().NoError(err)
	ft3, err := suite.client.NewForwardtest(context.Background(), params)
	suite.Require().NoError(err)

	// WHEN listing the forwardtests

	list, err := suite.client.ListForwardtests(context.Background(), api.ListForwardtestsWorkflowParams{})
	suite.Require().NoError(err)

	// THEN the list contains the forwardtests

	suite.Require().Contains(list, ft1)
	suite.Require().Contains(list, ft2)
	suite.Require().Contains(list, ft3)
}

func (suite *EndToEndSuite) TestCreateOrder() {
	// GIVEN a forwardtest

	params := api.CreateForwardtestWorkflowParams{
		Accounts: map[string]account.Account{
			"binance": {
				Balances: map[string]float64{
					"USDT": 1000000,
				},
			},
		},
		Callbacks: createTestCallbacks(),
	}
	ft, err := suite.client.NewForwardtest(context.Background(), params)
	suite.Require().NoError(err)

	// WHEN creating an order

	_, err = ft.CreateOrder(context.Background(), order.Order{
		Type:     order.TypeIsMarket,
		Side:     order.SideIsBuy,
		Exchange: "binance",
		Pair:     "BTC-USDT",
		Quantity: 1,
	})
	suite.Require().NoError(err)

	// THEN the balances are in order

	accounts, err := ft.ListAccounts(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(1.0, accounts["binance"].Balances["BTC"])
	suite.Require().NotEqual(1000000.0, accounts["binance"].Balances["USDT"])
}

func (suite *EndToEndSuite) TestListForwardtestAccounts() {
	// GIVEN a forwardtest with multiple accounts

	params := api.CreateForwardtestWorkflowParams{
		Accounts: map[string]account.Account{
			"binance": {
				Balances: map[string]float64{
					"USDT": 1000,
					"BTC":  0.5,
				},
			},
			"kucoin": {
				Balances: map[string]float64{
					"USDT": 2000,
					"ETH":  2.0,
				},
			},
		},
		Callbacks: createTestCallbacks(),
	}
	ft, err := suite.client.NewForwardtest(context.Background(), params)
	suite.Require().NoError(err)

	// WHEN listing the accounts

	accounts, err := ft.ListAccounts(context.Background())
	suite.Require().NoError(err)

	// THEN all accounts and balances are returned correctly

	suite.Require().Len(accounts, 2)

	// Check Binance account
	suite.Require().Equal(1000.0, accounts["binance"].Balances["USDT"])
	suite.Require().Equal(0.5, accounts["binance"].Balances["BTC"])

	// Check KuCoin account
	suite.Require().Equal(2000.0, accounts["kucoin"].Balances["USDT"])
	suite.Require().Equal(2.0, accounts["kucoin"].Balances["ETH"])
}

func (suite *EndToEndSuite) TestGetForwardtest() {
	// GIVEN a forwardtest with accounts and callbacks

	params := api.CreateForwardtestWorkflowParams{
		Accounts: map[string]account.Account{
			"binance": {
				Balances: map[string]float64{
					"USDT": 1000,
					"BTC":  0.5,
				},
			},
		},
		Callbacks: createTestCallbacks(),
	}
	ft, err := suite.client.NewForwardtest(context.Background(), params)
	suite.Require().NoError(err)

	// WHEN getting the forwardtest data

	retrievedFt, err := ft.Get(context.Background())
	suite.Require().NoError(err)

	// THEN the forwardtest data is retrieved correctly

	suite.Require().Equal(ft.ID, retrievedFt.ID)
	suite.Require().Len(retrievedFt.Accounts, 1)
	suite.Require().Equal(1000.0, retrievedFt.Accounts["binance"].Balances["USDT"])
	suite.Require().Equal(0.5, retrievedFt.Accounts["binance"].Balances["BTC"])
	suite.Require().Equal(createTestCallbacks(), retrievedFt.Callbacks)
	suite.Require().Equal(forwardtest.StatusReady, retrievedFt.Status)
}
