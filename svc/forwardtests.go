package svc

import (
	"github.com/cryptellation/candlesticks/pkg/clients"
	"github.com/cryptellation/forwardtests/api"
	"github.com/cryptellation/forwardtests/svc/db"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// Forwardtests is the interface for the forwardtests domain.
type Forwardtests interface {
	Register(w worker.Worker)

	CreateForwardtestWorkflow(
		ctx workflow.Context,
		params api.CreateForwardtestWorkflowParams,
	) (api.CreateForwardtestWorkflowResults, error)

	ListForwardtestsWorkflow(
		ctx workflow.Context,
		params api.ListForwardtestsWorkflowParams,
	) (api.ListForwardtestsWorkflowResults, error)

	CreateForwardtestOrderWorkflow(
		ctx workflow.Context,
		params api.CreateForwardtestOrderWorkflowParams,
	) (api.CreateForwardtestOrderWorkflowResults, error)

	ListForwardtestAccountsWorkflow(
		ctx workflow.Context,
		params api.ListForwardtestAccountsWorkflowParams,
	) (api.ListForwardtestAccountsWorkflowResults, error)

	GetForwardtestStatusWorkflow(
		ctx workflow.Context,
		params api.GetForwardtestStatusWorkflowParams,
	) (api.GetForwardtestStatusWorkflowResults, error)
}

var _ Forwardtests = &workflows{}

type workflows struct {
	db           db.DB
	candlesticks clients.WfClient
}

// New creates a new Forwardtests instance.
func New(db db.DB) Forwardtests {
	return &workflows{
		candlesticks: clients.NewWfClient(),
		db:           db,
	}
}

// Register registers the workflows to the worker.
func (wf *workflows) Register(worker worker.Worker) {
	worker.RegisterWorkflowWithOptions(wf.CreateForwardtestWorkflow, workflow.RegisterOptions{
		Name: api.CreateForwardtestWorkflowName,
	})
	worker.RegisterWorkflowWithOptions(wf.ListForwardtestsWorkflow, workflow.RegisterOptions{
		Name: api.ListForwardtestsWorkflowName,
	})
	worker.RegisterWorkflowWithOptions(wf.CreateForwardtestOrderWorkflow, workflow.RegisterOptions{
		Name: api.CreateForwardtestOrderWorkflowName,
	})
	worker.RegisterWorkflowWithOptions(wf.ListForwardtestAccountsWorkflow, workflow.RegisterOptions{
		Name: api.ListForwardtestAccountsWorkflowName,
	})
	worker.RegisterWorkflowWithOptions(wf.GetForwardtestStatusWorkflow, workflow.RegisterOptions{
		Name: api.GetForwardtestStatusWorkflowName,
	})

	worker.RegisterWorkflowWithOptions(ServiceInfoWorkflow, workflow.RegisterOptions{
		Name: api.ServiceInfoWorkflowName,
	})
}
