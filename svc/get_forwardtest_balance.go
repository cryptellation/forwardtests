package svc

import (
	"fmt"
	"time"

	candlesticksapi "github.com/cryptellation/candlesticks/api"
	"github.com/cryptellation/candlesticks/pkg/period"
	"github.com/cryptellation/forwardtests/api"
	"go.temporal.io/sdk/workflow"
)

var (
	// ErrNoActualPrice is the error when there is no actual price when requesting balance.
	ErrNoActualPrice = fmt.Errorf("no actual price")
)

const (
	// DefaultBalanceSymbol is the default symbol used to have the total balance.
	DefaultBalanceSymbol = "USDT"
)

// GetForwardtestBalanceWorkflow is the workflow to get the forwardtest balance.
func (wf *workflows) GetForwardtestBalanceWorkflow(
	ctx workflow.Context,
	params api.GetForwardtestBalanceWorkflowParams,
) (api.GetForwardtestBalanceWorkflowResults, error) {
	// Read forwardtest from database
	ft, err := wf.readForwardtestFromDB(ctx, params.ForwardtestID)
	if err != nil {
		return api.GetForwardtestBalanceWorkflowResults{},
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
				return api.GetForwardtestBalanceWorkflowResults{},
					fmt.Errorf("could not get candlesticks from service: %w", err)
			}

			c := csRes.List[len(csRes.List)-1]

			// Calculate value
			total += balance * c.Close
		}
	}

	return api.GetForwardtestBalanceWorkflowResults{
		Balance: total,
	}, nil
}
