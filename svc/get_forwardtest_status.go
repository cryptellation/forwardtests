package svc

import (
	"fmt"
	"time"

	candlesticksapi "github.com/cryptellation/candlesticks/api"
	"github.com/cryptellation/candlesticks/pkg/period"
	"github.com/cryptellation/forwardtests/api"
	"github.com/cryptellation/forwardtests/pkg/forwardtest"
	"go.temporal.io/sdk/workflow"
)

var (
	// ErrNoActualPrice is the error when there is no actual price when requesting status.
	ErrNoActualPrice = fmt.Errorf("no actual price")
)

const (
	// DefaultBalanceSymbol is the default symbol used to have the total balance.
	DefaultBalanceSymbol = "USDT"
)

// GetForwardtestStatusWorkflow is the workflow to get the forwardtest status.
func (wf *workflows) GetForwardtestStatusWorkflow(
	ctx workflow.Context,
	params api.GetForwardtestStatusWorkflowParams,
) (api.GetForwardtestStatusWorkflowResults, error) {
	// Read forwardtest from database
	ft, err := wf.readForwardtestFromDB(ctx, params.ForwardtestID)
	if err != nil {
		return api.GetForwardtestStatusWorkflowResults{},
			fmt.Errorf("could not read forwardtest from db: %w", err)
	}

	// Get value for each symbol in accounts
	total := 0.0
	for exchange, account := range ft.Accounts {
		for symbol, balance := range account.Balances {
			if symbol == DefaultBalanceSymbol {
				total += balance
				continue
			}

			// Get price
			start := time.Now().Add(-time.Minute * 10)
			end := time.Now()
			p := symbol + "-" + DefaultBalanceSymbol
			csRes, err := wf.candlesticks.ListCandlesticks(ctx, candlesticksapi.ListCandlesticksWorkflowParams{
				Exchange: exchange,
				Pair:     p,
				Period:   period.M1,
				Start:    &start,
				End:      &end,
				Limit:    1,
			}, &workflow.ChildWorkflowOptions{
				TaskQueue: candlesticksapi.WorkerTaskQueueName,
			})
			if err != nil {
				return api.GetForwardtestStatusWorkflowResults{},
					fmt.Errorf("could not get candlesticks from service: %w", err)
			}

			c := csRes.List[len(csRes.List)-1]

			// Calculate value
			total += balance * c.Close
		}
	}

	return api.GetForwardtestStatusWorkflowResults{
		Status: forwardtest.Status{
			Balance: total,
		},
	}, nil
}
