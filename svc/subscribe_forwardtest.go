package svc

import (
	"fmt"
	"time"

	"github.com/cryptellation/forwardtests/api"
	"github.com/cryptellation/forwardtests/pkg/forwardtest"
	"github.com/cryptellation/runtime"
	ticksapi "github.com/cryptellation/ticks/api"
	"github.com/cryptellation/ticks/pkg/tick"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

// SubscribeToPriceWorkflow subscribes to price updates for a forwardtest.
func (wf *workflows) SubscribeToPriceWorkflow(
	ctx workflow.Context,
	params api.SubscribeToPriceWorkflowParams,
) (api.SubscribeToPriceWorkflowResults, error) {
	// Create callback workflow
	callback := runtime.CallbackWorkflow{
		Name:          forwardNewPriceToForwardTestWorkflowName,
		TaskQueueName: workflow.GetInfo(ctx).TaskQueueName,
	}

	_, err := wf.ticks.ListenToTicks(ctx, ticksapi.RegisterForTicksListeningWorkflowParams{
		RequesterID: params.ForwardtestID,
		Exchange:    params.Exchange,
		Pair:        params.Pair,
		Callback:    callback,
	})

	return api.SubscribeToPriceWorkflowResults{}, err
}

const (
	// forwardNewPriceToForwardTestWorkflowName is the name of the ForwardNewPriceToForwardTestWorkflow.
	forwardNewPriceToForwardTestWorkflowName = "ForwardNewPriceToForwardTestWorkflow"
)

// forwardNewPriceToForwardTestWorkflow is a private proxy workflow that forwards price updates
// to the forwardtest callback.
func (wf *workflows) forwardNewPriceToForwardTestWorkflow(
	ctx workflow.Context,
	params ticksapi.ListenToTicksCallbackWorkflowParams,
) error {
	logger := workflow.GetLogger(ctx)
	logger.Debug("Forwarding price update to forwardtest callback",
		"forwardtest_id", params.RequesterID,
		"tick", params.Tick)

	// Read forwardtest from database to get callbacks
	ft, err := wf.readForwardtestFromDB(ctx, params.RequesterID)
	if err != nil {
		return fmt.Errorf("could not read forwardtest from db: %w", err)
	}

	// Check if forwardtest is finished - if so, unsubscribe from ticks and return
	if ft.Status == forwardtest.StatusFinished {
		return wf.handleFinishedForwardtest(ctx, params)
	}

	// Execute the OnNewPricesCallback workflow
	return wf.executeOnNewPricesCallback(ctx, params, ft)
}

// handleFinishedForwardtest handles the case when a forwardtest is finished.
func (wf *workflows) handleFinishedForwardtest(
	ctx workflow.Context,
	params ticksapi.ListenToTicksCallbackWorkflowParams,
) error {
	logger := workflow.GetLogger(ctx)
	logger.Debug("Forwardtest is finished, unsubscribing from ticks",
		"forwardtest_id", params.RequesterID)

	// Unsubscribe from ticks using exchange and pair from the tick
	_, err := wf.ticks.StopListeningToTicks(ctx, ticksapi.UnregisterFromTicksListeningWorkflowParams{
		RequesterID: params.RequesterID,
		Exchange:    params.Tick.Exchange,
		Pair:        params.Tick.Pair,
	})
	if err != nil {
		logger.Error("Failed to unsubscribe from ticks", "error", err)
		// Don't return error here as we want to exit gracefully
	}

	return nil
}

// executeOnNewPricesCallback executes the OnNewPricesCallback workflow.
func (wf *workflows) executeOnNewPricesCallback(
	ctx workflow.Context,
	params ticksapi.ListenToTicksCallbackWorkflowParams,
	ft forwardtest.Forwardtest,
) error {
	// Create child workflow options
	opts := workflow.ChildWorkflowOptions{
		// Unique identifier for this child workflow execution
		WorkflowID: fmt.Sprintf("forwardtest-%s-on-new-prices-%s",
			params.RequesterID.String(), params.Tick.Time.Format(time.RFC3339Nano)),
		// Task queue where the child workflow will be executed
		TaskQueue: ft.Callbacks.OnNewPricesCallback.TaskQueueName,
		// Maximum time allowed for the child workflow to complete
		WorkflowExecutionTimeout: time.Second * 30,
		// Policy for what happens to the child workflow when parent closes
		// ABANDON means the child continues running independently
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
		// Policy for handling duplicate workflow IDs
		// REJECT_DUPLICATE prevents concurrent execution of workflows with the same ID
		WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
	}

	// Check if the timeout is set
	if ft.Callbacks.OnNewPricesCallback.ExecutionTimeout > 0 {
		opts.WorkflowExecutionTimeout = ft.Callbacks.OnNewPricesCallback.ExecutionTimeout
	}

	// Execute the OnNewPricesCallback workflow
	err := workflow.ExecuteChildWorkflow(
		workflow.WithChildOptions(ctx, opts),
		ft.Callbacks.OnNewPricesCallback.Name,
		runtime.OnNewPricesCallbackWorkflowParams{
			Context: runtime.Context{
				ID:              params.RequesterID,
				Mode:            runtime.ModeForwardtest,
				Now:             params.Tick.Time,
				ParentTaskQueue: workflow.GetInfo(ctx).TaskQueueName,
			},
			Ticks: []tick.Tick{params.Tick},
		}).Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not execute OnNewPricesCallback workflow: %w", err)
	}

	logger := workflow.GetLogger(ctx)
	logger.Debug("Successfully forwarded price update to forwardtest callback",
		"forwardtest_id", params.RequesterID)
	return nil
}
