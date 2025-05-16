package svc

import (
	"github.com/cryptellation/forwardtests/pkg/forwardtest"
	"github.com/cryptellation/forwardtests/svc/db"
	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"
)

func (wf *workflows) readForwardtestFromDB(ctx workflow.Context, id uuid.UUID) (forwardtest.Forwardtest, error) {
	var readRes db.ReadForwardtestActivityResult
	err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, db.DefaultActivityOptions()),
		wf.db.ReadForwardtestActivity, db.ReadForwardtestActivityParams{
			ID: id,
		}).Get(ctx, &readRes)
	if err != nil {
		return forwardtest.Forwardtest{}, err
	}

	return readRes.Forwardtest, nil
}
