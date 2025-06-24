//go:build e2e
// +build e2e

package test

import (
	"context"
	"time"

	"github.com/cryptellation/forwardtests/api"
	"github.com/cryptellation/forwardtests/pkg/clients"
	"github.com/cryptellation/forwardtests/pkg/forwardtest"
	"github.com/cryptellation/runtime"
	"github.com/cryptellation/runtime/account"
	"github.com/google/uuid"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type testRunner struct {
	Suite *EndToEndSuite

	ForwardtestID uuid.UUID
	Accounts      map[string]account.Account

	WfClient clients.WfClient

	OnInitCalls      int
	OnNewPricesCalls int
	OnExitCalls      int
}

func (r *testRunner) Name() string {
	return "ForwardtestE2eRunner"
}

func (r *testRunner) OnInit(ctx workflow.Context, params runtime.OnInitCallbackWorkflowParams) error {
	checkForwardtestRunContext(r.Suite, params.Context, r.ForwardtestID)

	// Check that forwardtest status is "running" during OnInit
	r.checkForwardtestStatus(ctx, forwardtest.StatusRunning)

	// Subscribe to price
	_, err := r.WfClient.SubscribeToPrice(ctx, api.SubscribeToPriceWorkflowParams{
		ForwardtestID: r.ForwardtestID,
		Exchange:      "binance",
		Pair:          "BTC-USDT",
	})
	r.Suite.Require().NoError(err)

	r.OnInitCalls++
	return err
}

func (r *testRunner) OnNewPrices(ctx workflow.Context, params runtime.OnNewPricesCallbackWorkflowParams) error {
	checkForwardtestRunContext(r.Suite, params.Context, r.ForwardtestID)

	// Check that forwardtest status is "running" during OnNewPrices
	r.checkForwardtestStatus(ctx, forwardtest.StatusRunning)

	// TODO(#6): test order passing in OnNewPrices

	r.OnNewPricesCalls++
	return nil
}

func (r *testRunner) OnExit(ctx workflow.Context, params runtime.OnExitCallbackWorkflowParams) error {
	checkForwardtestRunContext(r.Suite, params.Context, r.ForwardtestID)

	// Check that forwardtest status is "finished" during OnExit
	r.checkForwardtestStatus(ctx, forwardtest.StatusFinished)

	r.OnExitCalls++
	return nil
}

// checkForwardtestStatus verifies that the forwardtest has the expected status.
func (r *testRunner) checkForwardtestStatus(ctx workflow.Context, expectedStatus forwardtest.Status) {
	result, err := r.WfClient.GetForwardtest(ctx, api.GetForwardtestWorkflowParams{
		ForwardtestID: r.ForwardtestID,
	})
	r.Suite.Require().NoError(err)
	r.Suite.Require().Equal(expectedStatus, result.Forwardtest.Status)
}

func (suite *EndToEndSuite) TestForwardtestRun() {
	// GIVEN a running worker
	tq := "ForwardtestE2eRunner-TaskQueue"
	w := worker.New(suite.temporalclient, tq, worker.Options{})
	go func() {
		if err := w.Run(nil); err != nil {
			suite.Require().NoError(err)
		}
	}()
	defer w.Stop()

	// AND a runner

	accounts := map[string]account.Account{
		"binance": {
			Balances: map[string]float64{
				"BTC": 1,
			},
		},
	}
	r := &testRunner{
		Accounts: accounts,
		Suite:    suite,
		WfClient: clients.NewWfClient(),
	}

	// WHEN creating a new forwardtest

	params := api.CreateForwardtestWorkflowParams{
		Accounts:  accounts,
		Callbacks: runtime.RegisterRunnable(w, tq, r),
	}
	forwardtest, err := suite.client.NewForwardtest(context.Background(), params)

	// THEN no error is returned

	suite.Require().NoError(err)

	// WHEN running the forwardtest with a runner

	r.ForwardtestID = forwardtest.ID // Add forwardtest ID to runner for checking forwardtest run context
	err = forwardtest.Run(context.Background())

	// THEN no error is returned

	suite.Require().NoError(err)

	// AND we wait for the OnInit callback to be called
	suite.Require().Eventually(func() bool {
		return r.OnInitCalls >= 1
	}, time.Second*10, time.Millisecond*100)

	// AND we wait for the OnNewPrices callback to be called at least twice

	suite.Require().Eventually(func() bool {
		return r.OnNewPricesCalls >= 2
	}, 10*time.Minute, time.Millisecond*100)

	// WHEN stopping the forwardtest
	err = forwardtest.Stop(context.Background())

	// THEN no error is returned

	suite.Require().NoError(err)

	// AND the we wait for the OnExit callback to be called
	suite.Require().Eventually(func() bool {
		return r.OnExitCalls >= 1
	}, time.Second*10, time.Millisecond*100)
}

func checkForwardtestRunContext(suite *EndToEndSuite, ctx runtime.Context, forwardtestID uuid.UUID) {
	suite.Require().Equal(forwardtestID, ctx.ID)
	suite.Require().Equal(runtime.ModeForwardtest, ctx.Mode)
	suite.Require().NotEmpty(ctx.ParentTaskQueue)
}
