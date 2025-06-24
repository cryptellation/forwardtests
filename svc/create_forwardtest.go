package svc

import (
	"fmt"

	"github.com/cryptellation/forwardtests/api"
	"github.com/cryptellation/forwardtests/pkg/forwardtest"
	"github.com/cryptellation/forwardtests/svc/db"
	"go.temporal.io/sdk/workflow"
)

// CreateForwardtestWorkflow creates a new forwardtest and saves it to the database.
func (wf *workflows) CreateForwardtestWorkflow(
	ctx workflow.Context,
	params api.CreateForwardtestWorkflowParams,
) (api.CreateForwardtestWorkflowResults, error) {
	// Validate callbacks
	if err := params.Callbacks.Validate(); err != nil {
		return api.CreateForwardtestWorkflowResults{}, fmt.Errorf("validating callbacks: %w", err)
	}

	payload := forwardtest.NewForwardtestParams{
		Accounts:  params.Accounts,
		Callbacks: params.Callbacks,
	}

	// Create new forwardtest and save it to database
	ft, err := forwardtest.New(payload)
	if err != nil {
		return api.CreateForwardtestWorkflowResults{}, fmt.Errorf("creating a new forwardtest from request: %w", err)
	}

	err = workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, db.DefaultActivityOptions()),
		wf.db.CreateForwardtestActivity, db.CreateForwardtestActivityParams{
			Forwardtest: ft,
		}).Get(ctx, nil)
	if err != nil {
		return api.CreateForwardtestWorkflowResults{}, fmt.Errorf("adding forwardtest to db: %w", err)
	}

	return api.CreateForwardtestWorkflowResults{
		ID: ft.ID,
	}, nil
}
