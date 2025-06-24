package clients

import (
	"context"

	"github.com/cryptellation/forwardtests/api"
	temporalclient "go.temporal.io/sdk/client"
)

// RawClient is a client for the cryptellation backtests service with just the
// calls to the temporal workflows.
type RawClient interface {
	CreateForwardtest(
		ctx context.Context,
		params api.CreateForwardtestWorkflowParams,
	) (api.CreateForwardtestWorkflowResults, error)
	GetForwardtest(
		ctx context.Context,
		params api.GetForwardtestWorkflowParams,
	) (api.GetForwardtestWorkflowResults, error)
	GetForwardtestBalance(
		ctx context.Context,
		params api.GetForwardtestBalanceWorkflowParams,
	) (api.GetForwardtestBalanceWorkflowResults, error)
	ListForwardtests(
		ctx context.Context,
		params api.ListForwardtestsWorkflowParams,
	) (api.ListForwardtestsWorkflowResults, error)
	CreateForwardtestOrder(
		ctx context.Context,
		params api.CreateForwardtestOrderWorkflowParams,
	) (api.CreateForwardtestOrderWorkflowResults, error)
	ListForwardtestAccounts(
		ctx context.Context,
		params api.ListForwardtestAccountsWorkflowParams,
	) (api.ListForwardtestAccountsWorkflowResults, error)
	RunForwardtest(
		ctx context.Context,
		params api.RunForwardtestWorkflowParams,
	) (api.RunForwardtestWorkflowResults, error)
	StopForwardtest(
		ctx context.Context,
		params api.StopForwardtestWorkflowParams,
	) (api.StopForwardtestWorkflowResults, error)
}

var _ RawClient = raw{}

type raw struct {
	temporal temporalclient.Client
}

// NewRaw creates a new raw client to execute temporal workflows.
func NewRaw(cl temporalclient.Client) RawClient {
	return &raw{temporal: cl}
}

func (c raw) CreateForwardtest(
	ctx context.Context,
	params api.CreateForwardtestWorkflowParams,
) (api.CreateForwardtestWorkflowResults, error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.CreateForwardtestWorkflowName, params)
	if err != nil {
		return api.CreateForwardtestWorkflowResults{}, err
	}

	// Get result and return
	var res api.CreateForwardtestWorkflowResults
	err = exec.Get(ctx, &res)

	return res, err
}

func (c raw) GetForwardtest(
	ctx context.Context,
	params api.GetForwardtestWorkflowParams,
) (api.GetForwardtestWorkflowResults, error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.GetForwardtestWorkflowName, params)
	if err != nil {
		return api.GetForwardtestWorkflowResults{}, err
	}

	// Get result and return
	var res api.GetForwardtestWorkflowResults
	err = exec.Get(ctx, &res)

	return res, err
}

func (c raw) GetForwardtestBalance(
	ctx context.Context,
	params api.GetForwardtestBalanceWorkflowParams,
) (api.GetForwardtestBalanceWorkflowResults, error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.GetForwardtestBalanceWorkflowName, params)
	if err != nil {
		return api.GetForwardtestBalanceWorkflowResults{}, err
	}

	// Get result and return
	var res api.GetForwardtestBalanceWorkflowResults
	err = exec.Get(ctx, &res)

	return res, err
}

func (c raw) ListForwardtests(
	ctx context.Context,
	params api.ListForwardtestsWorkflowParams,
) (api.ListForwardtestsWorkflowResults, error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.ListForwardtestsWorkflowName, params)
	if err != nil {
		return api.ListForwardtestsWorkflowResults{}, err
	}

	// Get result and return
	var res api.ListForwardtestsWorkflowResults
	err = exec.Get(ctx, &res)

	return res, err
}

func (c raw) CreateForwardtestOrder(
	ctx context.Context,
	params api.CreateForwardtestOrderWorkflowParams,
) (api.CreateForwardtestOrderWorkflowResults, error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.CreateForwardtestOrderWorkflowName, params)
	if err != nil {
		return api.CreateForwardtestOrderWorkflowResults{}, err
	}

	// Get result and return
	var res api.CreateForwardtestOrderWorkflowResults
	err = exec.Get(ctx, &res)

	return res, err
}

func (c raw) ListForwardtestAccounts(
	ctx context.Context,
	params api.ListForwardtestAccountsWorkflowParams,
) (api.ListForwardtestAccountsWorkflowResults, error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.ListForwardtestAccountsWorkflowName, params)
	if err != nil {
		return api.ListForwardtestAccountsWorkflowResults{}, err
	}

	// Get result and return
	var res api.ListForwardtestAccountsWorkflowResults
	err = exec.Get(ctx, &res)

	return res, err
}

func (c raw) RunForwardtest(
	ctx context.Context,
	params api.RunForwardtestWorkflowParams,
) (api.RunForwardtestWorkflowResults, error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.RunForwardtestWorkflowName, params)
	if err != nil {
		return api.RunForwardtestWorkflowResults{}, err
	}

	// Get result and return
	var res api.RunForwardtestWorkflowResults
	err = exec.Get(ctx, &res)

	return res, err
}

func (c raw) StopForwardtest(
	ctx context.Context,
	params api.StopForwardtestWorkflowParams,
) (api.StopForwardtestWorkflowResults, error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.StopForwardtestWorkflowName, params)
	if err != nil {
		return api.StopForwardtestWorkflowResults{}, err
	}

	// Get result and return
	var res api.StopForwardtestWorkflowResults
	err = exec.Get(ctx, &res)

	return res, err
}
