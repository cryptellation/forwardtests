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

// StopForwardtestWorkflow stops a forwardtest by executing the exit callback.
func (wf *workflows) StopForwardtestWorkflow(
	ctx workflow.Context,
	params forwardtestsapi.StopForwardtestWorkflowParams,
) (forwardtestsapi.StopForwardtestWorkflowResults, error) {
	// Load forwardtest from database to get callbacks
	ft, err := wf.readForwardtestFromDB(ctx, params.ForwardtestID)
	if err != nil {
		return forwardtestsapi.StopForwardtestWorkflowResults{}, fmt.Errorf("loading forwardtest from database: %w", err)
	}

	// Update forwardtest status to finished
	ft.Status = forwardtest.StatusFinished
	err = workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, db.DefaultActivityOptions()),
		wf.db.UpdateForwardtestActivity, db.UpdateForwardtestActivityParams{
			Forwardtest: ft,
		}).Get(ctx, nil)
	if err != nil {
		return forwardtestsapi.StopForwardtestWorkflowResults{},
			fmt.Errorf("updating forwardtest status to finished: %w", err)
	}

	// Execute the exit callback workflow
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		// Unique identifier for this child workflow execution
		WorkflowID: fmt.Sprintf("forwardtest-%s-exit", params.ForwardtestID.String()),
		// Task queue where the child workflow will be executed
		TaskQueue: ft.Callbacks.OnExitCallback.TaskQueueName,
		// Maximum time allowed for the child workflow to complete
		WorkflowExecutionTimeout: time.Second * 30,
	}

	// Check if the timeout is set
	if ft.Callbacks.OnExitCallback.ExecutionTimeout > 0 {
		childWorkflowOptions.WorkflowExecutionTimeout = ft.Callbacks.OnExitCallback.ExecutionTimeout
	}

	var result forwardtestsapi.StopForwardtestWorkflowResults
	err = workflow.ExecuteChildWorkflow(
		workflow.WithChildOptions(ctx, childWorkflowOptions),
		ft.Callbacks.OnExitCallback.Name,
		runtime.OnExitCallbackWorkflowParams{
			Context: runtime.Context{
				ID:              params.ForwardtestID,
				Mode:            runtime.ModeForwardtest,
				Now:             workflow.Now(ctx),
				ParentTaskQueue: workflow.GetInfo(ctx).TaskQueueName,
			},
		},
	).Get(ctx, &result)

	if err != nil {
		return forwardtestsapi.StopForwardtestWorkflowResults{}, fmt.Errorf("could not execute exit callback: %w", err)
	}

	return result, nil
}
