package svc

import (
	"fmt"

	candlesticksapi "github.com/cryptellation/candlesticks/api"
	"github.com/cryptellation/candlesticks/pkg/period"
	"github.com/cryptellation/forwardtests/api"
	"github.com/cryptellation/forwardtests/svc/db"
	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"
)

// CreateForwardtestOrderWorkflow creates a new forwardtest order and saves it to the database.
func (wf *workflows) CreateForwardtestOrderWorkflow(
	ctx workflow.Context,
	params api.CreateForwardtestOrderWorkflowParams,
) (api.CreateForwardtestOrderWorkflowResults, error) {
	logger := workflow.GetLogger(ctx)

	if params.Order.ID == uuid.Nil {
		params.Order.ID = uuid.New()
	}

	logger.Debug("Creating order on forwardtest",
		"order", params.Order,
		"forwardtest_id", params.ForwardtestID.String())

	// Read forwardtest from database
	ft, err := wf.readForwardtestFromDB(ctx, params.ForwardtestID)
	if err != nil {
		return api.CreateForwardtestOrderWorkflowResults{},
			fmt.Errorf("could not read forwardtest from db: %w", err)
	}

	// Get candlestick for order validation
	now := workflow.Now(ctx)
	csRes, err := wf.candlesticks.ListCandlesticks(ctx, candlesticksapi.ListCandlesticksWorkflowParams{
		Exchange: params.Order.Exchange,
		Pair:     params.Order.Pair,
		Period:   period.M1,
		Start:    &now,
		End:      &now,
		Limit:    1,
	}, &workflow.ChildWorkflowOptions{
		TaskQueue: candlesticksapi.WorkerTaskQueueName,
	})
	if err != nil {
		return api.CreateForwardtestOrderWorkflowResults{},
			fmt.Errorf("could not get candlesticks from service: %w", err)
	}

	cs := csRes.List[0]

	logger.Info("Adding order to forwardtest",
		"order", params.Order,
		"forwardtest", params.ForwardtestID.String())
	if err := ft.AddOrder(params.Order, cs); err != nil {
		return api.CreateForwardtestOrderWorkflowResults{}, err
	}

	// Save forwardtest to database
	err = workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, db.DefaultActivityOptions()),
		wf.db.UpdateForwardtestActivity, db.UpdateForwardtestActivityParams{
			Forwardtest: ft,
		}).Get(ctx, nil)
	if err != nil {
		return api.CreateForwardtestOrderWorkflowResults{}, err
	}

	return api.CreateForwardtestOrderWorkflowResults{}, nil
}
