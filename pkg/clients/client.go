package clients

import (
	"context"

	"github.com/cryptellation/forwardtests/api"
	temporalclient "go.temporal.io/sdk/client"
)

// Client is a client for the cryptellation forwardtests service.
type Client interface {
	// NewForwardtest creates a new forwardtest.
	NewForwardtest(
		ctx context.Context,
		params api.CreateForwardtestWorkflowParams,
	) (Forwardtest, error)
	// ListForwardtests lists the forwardtests.
	ListForwardtests(
		ctx context.Context,
		params api.ListForwardtestsWorkflowParams,
	) ([]Forwardtest, error)
	// Info calls the service info.
	Info(ctx context.Context) (api.ServiceInfoResults, error)
}

type client struct {
	temporal  temporalclient.Client
	rawClient RawClient
}

// New creates a new client to execute temporal workflows.
func New(cl temporalclient.Client) Client {
	return &client{
		temporal:  cl,
		rawClient: NewRaw(cl),
	}
}

// NewForwardtest creates a new forwardtest.
func (c client) NewForwardtest(
	ctx context.Context,
	params api.CreateForwardtestWorkflowParams,
) (Forwardtest, error) {
	res, err := c.rawClient.CreateForwardtest(ctx, params)
	return Forwardtest{
		ID:        res.ID,
		rawClient: c.rawClient,
	}, err
}

// ListForwardtests lists the forwardtests.
func (c client) ListForwardtests(
	ctx context.Context,
	params api.ListForwardtestsWorkflowParams,
) ([]Forwardtest, error) {
	res, err := c.rawClient.ListForwardtests(ctx, params)
	if err != nil {
		return nil, err
	}

	forwardtests := make([]Forwardtest, len(res.Forwardtests))
	for i, ft := range res.Forwardtests {
		forwardtests[i] = Forwardtest{
			ID:        ft.ID,
			rawClient: c.rawClient,
		}
	}

	return forwardtests, nil
}

// Info calls the service info.
func (c client) Info(ctx context.Context) (res api.ServiceInfoResults, err error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.ServiceInfoWorkflowName)
	if err != nil {
		return api.ServiceInfoResults{}, err
	}

	// Get result and return
	err = exec.Get(ctx, &res)
	return res, err
}
