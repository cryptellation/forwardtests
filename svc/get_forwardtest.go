package svc

import (
	"fmt"

	forwardtestsapi "github.com/cryptellation/forwardtests/api"
	"go.temporal.io/sdk/workflow"
)

// GetForwardtestWorkflow retrieves a forwardtest from the database by its ID.
func (wf *workflows) GetForwardtestWorkflow(
	ctx workflow.Context,
	params forwardtestsapi.GetForwardtestWorkflowParams,
) (forwardtestsapi.GetForwardtestWorkflowResults, error) {
	// Read forwardtest from database
	ft, err := wf.readForwardtestFromDB(ctx, params.ForwardtestID)
	if err != nil {
		return forwardtestsapi.GetForwardtestWorkflowResults{}, fmt.Errorf("could not read forwardtest from db: %w", err)
	}

	return forwardtestsapi.GetForwardtestWorkflowResults{
		Forwardtest: ft,
	}, nil
}
