package svc

import (
	"fmt"
	"time"

	forwardtestsapi "github.com/cryptellation/forwardtests/api"
	"github.com/cryptellation/forwardtests/pkg/forwardtest"
	"github.com/cryptellation/forwardtests/svc/db"
	"github.com/cryptellation/runtime"
	"go.temporal.io/sdk/workflow"
)

// RunForwardtestWorkflow runs a forwardtest by executing the init callback.
func (wf *workflows) RunForwardtestWorkflow(
	ctx workflow.Context,
	params forwardtestsapi.RunForwardtestWorkflowParams,
) (forwardtestsapi.RunForwardtestWorkflowResults, error) {
	// Load forwardtest from database to get callbacks
	ft, err := wf.readForwardtestFromDB(ctx, params.ForwardtestID)
	if err != nil {
		return forwardtestsapi.RunForwardtestWorkflowResults{}, fmt.Errorf("loading forwardtest from database: %w", err)
	}

	// Update forwardtest status to running
	ft.Status = forwardtest.StatusRunning
	err = workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, db.DefaultActivityOptions()),
		wf.db.UpdateForwardtestActivity, db.UpdateForwardtestActivityParams{
			Forwardtest: ft,
		}).Get(ctx, nil)
	if err != nil {
		return forwardtestsapi.RunForwardtestWorkflowResults{}, fmt.Errorf("updating forwardtest status to running: %w", err)
	}

	// Execute the init callback workflow
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		// Unique identifier for this child workflow execution
		WorkflowID: fmt.Sprintf("forwardtest-%s-init", params.ForwardtestID.String()),
		// Task queue where the child workflow will be executed
		TaskQueue: ft.Callbacks.OnInitCallback.TaskQueueName,
		// Maximum time allowed for the child workflow to complete
		WorkflowExecutionTimeout: time.Second * 30,
	}

	// Check if the timeout is set
	if ft.Callbacks.OnInitCallback.ExecutionTimeout > 0 {
		childWorkflowOptions.WorkflowExecutionTimeout = ft.Callbacks.OnInitCallback.ExecutionTimeout
	}

	var result forwardtestsapi.RunForwardtestWorkflowResults
	err = workflow.ExecuteChildWorkflow(
		workflow.WithChildOptions(ctx, childWorkflowOptions),
		ft.Callbacks.OnInitCallback.Name,
		runtime.OnInitCallbackWorkflowParams{
			Context: runtime.Context{
				ID:              params.ForwardtestID,
				Mode:            runtime.ModeForwardtest,
				Now:             workflow.Now(ctx),
				ParentTaskQueue: workflow.GetInfo(ctx).TaskQueueName,
			},
		},
	).Get(ctx, &result)

	if err != nil {
		return forwardtestsapi.RunForwardtestWorkflowResults{}, fmt.Errorf("could not execute init callback: %w", err)
	}

	return result, nil
}
