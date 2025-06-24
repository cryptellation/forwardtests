package svc

import (
	"fmt"
	"time"

	forwardtestsapi "github.com/cryptellation/forwardtests/api"
	"github.com/cryptellation/runtime"
	"go.temporal.io/sdk/workflow"
)

// StartForwardtestWorkflow starts a forwardtest by executing the init callback.
func (wf *workflows) StartForwardtestWorkflow(
	ctx workflow.Context,
	params forwardtestsapi.StartForwardtestWorkflowParams,
) (forwardtestsapi.StartForwardtestWorkflowResults, error) {
	// Load forwardtest from database to get callbacks
	ft, err := wf.readForwardtestFromDB(ctx, params.ForwardtestID)
	if err != nil {
		return forwardtestsapi.StartForwardtestWorkflowResults{}, fmt.Errorf("loading forwardtest from database: %w", err)
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

	var result forwardtestsapi.StartForwardtestWorkflowResults
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
		return forwardtestsapi.StartForwardtestWorkflowResults{}, fmt.Errorf("could not execute init callback: %w", err)
	}

	return result, nil
}
