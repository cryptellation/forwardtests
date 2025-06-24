package svc

import (
	candlesticksclients "github.com/cryptellation/candlesticks/pkg/clients"
	"github.com/cryptellation/forwardtests/api"
	"github.com/cryptellation/forwardtests/svc/db"
	tickclients "github.com/cryptellation/ticks/pkg/clients"
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

	GetForwardtestBalanceWorkflow(
		ctx workflow.Context,
		params api.GetForwardtestBalanceWorkflowParams,
	) (api.GetForwardtestBalanceWorkflowResults, error)

	GetForwardtestWorkflow(
		ctx workflow.Context,
		params api.GetForwardtestWorkflowParams,
	) (api.GetForwardtestWorkflowResults, error)

	RunForwardtestWorkflow(
		ctx workflow.Context,
		params api.RunForwardtestWorkflowParams,
	) (api.RunForwardtestWorkflowResults, error)

	StopForwardtestWorkflow(
		ctx workflow.Context,
		params api.StopForwardtestWorkflowParams,
	) (api.StopForwardtestWorkflowResults, error)

	SubscribeToPriceWorkflow(
		ctx workflow.Context,
		params api.SubscribeToPriceWorkflowParams,
	) (api.SubscribeToPriceWorkflowResults, error)
}

var _ Forwardtests = &workflows{}

type workflows struct {
	db           db.DB
	candlesticks candlesticksclients.WfClient
	ticks        tickclients.WfClient
}

// New creates a new Forwardtests instance.
func New(db db.DB) Forwardtests {
	return &workflows{
		candlesticks: candlesticksclients.NewWfClient(),
		ticks:        tickclients.NewWfClient(),
		db:           db,
	}
}

// Register registers the workflows to the worker.
func (wf *workflows) Register(worker worker.Worker) {
	// Private workflows
	worker.RegisterWorkflowWithOptions(wf.forwardNewPriceToForwardTestWorkflow, workflow.RegisterOptions{
		Name: forwardNewPriceToForwardTestWorkflowName,
	})

	// Public workflows
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
	worker.RegisterWorkflowWithOptions(wf.GetForwardtestBalanceWorkflow, workflow.RegisterOptions{
		Name: api.GetForwardtestBalanceWorkflowName,
	})
	worker.RegisterWorkflowWithOptions(wf.GetForwardtestWorkflow, workflow.RegisterOptions{
		Name: api.GetForwardtestWorkflowName,
	})
	worker.RegisterWorkflowWithOptions(wf.RunForwardtestWorkflow, workflow.RegisterOptions{
		Name: api.RunForwardtestWorkflowName,
	})
	worker.RegisterWorkflowWithOptions(wf.StopForwardtestWorkflow, workflow.RegisterOptions{
		Name: api.StopForwardtestWorkflowName,
	})
	worker.RegisterWorkflowWithOptions(wf.SubscribeToPriceWorkflow, workflow.RegisterOptions{
		Name: api.SubscribeToPriceWorkflowName,
	})

	worker.RegisterWorkflowWithOptions(ServiceInfoWorkflow, workflow.RegisterOptions{
		Name: api.ServiceInfoWorkflowName,
	})
}
