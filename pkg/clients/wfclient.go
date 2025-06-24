package clients

import (
	"time"

	"github.com/cryptellation/forwardtests/api"
	"go.temporal.io/sdk/workflow"
)

// WfClient is a client for the cryptellation forwardtests service from a workflow perspective.
type WfClient interface {
	// CreateForwardtestOrder creates a new order for a forwardtest.
	CreateForwardtestOrder(
		ctx workflow.Context,
		params api.CreateForwardtestOrderWorkflowParams,
	) (api.CreateForwardtestOrderWorkflowResults, error)

	// GetForwardtest retrieves a forwardtest from the database by its ID.
	GetForwardtest(
		ctx workflow.Context,
		params api.GetForwardtestWorkflowParams,
	) (api.GetForwardtestWorkflowResults, error)

	// SubscribeToPrice subscribes to specific price updates for a forwardtest.
	SubscribeToPrice(
		ctx workflow.Context,
		params api.SubscribeToPriceWorkflowParams,
	) (api.SubscribeToPriceWorkflowResults, error)
}

type wfClient struct{}

// NewWfClient creates a new workflow client.
// This client is used to call workflows from within other workflows.
// It is not used to call workflows from outside the workflow environment.
func NewWfClient() WfClient {
	return wfClient{}
}

// CreateForwardtestOrder creates a new order for a forwardtest.
func (c wfClient) CreateForwardtestOrder(
	ctx workflow.Context,
	params api.CreateForwardtestOrderWorkflowParams,
) (api.CreateForwardtestOrderWorkflowResults, error) {
	var res api.CreateForwardtestOrderWorkflowResults
	err := workflow.ExecuteActivity(
		ctx,
		api.CreateForwardtestOrderWorkflowName,
		params,
	).Get(ctx, &res)
	return res, err
}

// GetForwardtest retrieves a forwardtest from the database by its ID.
func (c wfClient) GetForwardtest(
	ctx workflow.Context,
	params api.GetForwardtestWorkflowParams,
) (api.GetForwardtestWorkflowResults, error) {
	// Set child workflow options with timeout
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		TaskQueue:                api.WorkerTaskQueueName,
		WorkflowExecutionTimeout: 10 * time.Second,
	}
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)

	// Execute the GetForwardtestWorkflow as a child workflow
	var res api.GetForwardtestWorkflowResults
	err := workflow.ExecuteChildWorkflow(ctx, api.GetForwardtestWorkflowName, params).Get(ctx, &res)
	return res, err
}

// SubscribeToPrice subscribes to price updates for a forwardtest.
func (c wfClient) SubscribeToPrice(
	ctx workflow.Context,
	params api.SubscribeToPriceWorkflowParams,
) (api.SubscribeToPriceWorkflowResults, error) {
	// Set child workflow options
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)

	// Execute the SubscribeToPriceWorkflow as a child workflow
	var res api.SubscribeToPriceWorkflowResults
	err := workflow.ExecuteChildWorkflow(ctx, api.SubscribeToPriceWorkflowName, params).Get(ctx, &res)

	if err != nil {
		return api.SubscribeToPriceWorkflowResults{}, err
	}

	return res, nil
}
