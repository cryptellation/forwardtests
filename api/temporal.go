package api

import (
	"github.com/cryptellation/forwardtests/pkg/forwardtest"
	"github.com/cryptellation/runtime"
	"github.com/cryptellation/runtime/account"
	"github.com/cryptellation/runtime/order"
	"github.com/google/uuid"
)

const (
	// WorkerTaskQueueName is the name of the task queue for the cryptellation worker.
	WorkerTaskQueueName = "CryptellationForwardtestsTaskQueue"
)

// CreateForwardtestWorkflowName is the name of the CreateForwardtestWorkflow.
const CreateForwardtestWorkflowName = "CreateForwardtestWorkflow"

type (
	// CreateForwardtestWorkflowParams is the input for the CreateForwardtestWorkflow.
	CreateForwardtestWorkflowParams struct {
		Accounts  map[string]account.Account
		Callbacks runtime.Callbacks
	}

	// CreateForwardtestWorkflowResults is the output for the CreateForwardtestWorkflow.
	CreateForwardtestWorkflowResults struct {
		ID uuid.UUID
	}
)

// ListForwardtestsWorkflowName is the name of the ListForwardtestsWorkflow.
const ListForwardtestsWorkflowName = "ListForwardtestsWorkflow"

type (
	// ListForwardtestsWorkflowParams is the input for the ListForwardtestsWorkflow.
	ListForwardtestsWorkflowParams struct{}

	// ListForwardtestsWorkflowResults is the output for the ListForwardtestsWorkflow.
	ListForwardtestsWorkflowResults struct {
		Forwardtests []forwardtest.Forwardtest
	}
)

// CreateForwardtestOrderWorkflowName is the name of the CreateForwardtestOrderWorkflow.
const CreateForwardtestOrderWorkflowName = "CreateForwardtestOrderWorkflow"

type (
	// CreateForwardtestOrderWorkflowParams is the input for the CreateForwardtestOrderWorkflow.
	CreateForwardtestOrderWorkflowParams struct {
		ForwardtestID uuid.UUID
		Order         order.Order
	}

	// CreateForwardtestOrderWorkflowResults is the output for the CreateForwardtestOrderWorkflow.
	CreateForwardtestOrderWorkflowResults struct{}
)

// ListForwardtestAccountsWorkflowName is the name of the ListForwardtestAccountsWorkflow.
const ListForwardtestAccountsWorkflowName = "ListForwardtestAccountsWorkflow"

type (
	// ListForwardtestAccountsWorkflowParams is the input for the ListForwardtestAccountsWorkflow.
	ListForwardtestAccountsWorkflowParams struct {
		ForwardtestID uuid.UUID
	}

	// ListForwardtestAccountsWorkflowResults is the output for the ListForwardtestAccountsWorkflow.
	ListForwardtestAccountsWorkflowResults struct {
		Accounts map[string]account.Account
	}
)

// GetForwardtestStatusWorkflowName is the name of the GetForwardtestStatusWorkflow.
const GetForwardtestStatusWorkflowName = "GetForwardtestStatusWorkflow"

type (
	// GetForwardtestStatusWorkflowParams is the input for the GetForwardtestStatusWorkflow.
	GetForwardtestStatusWorkflowParams struct {
		ForwardtestID uuid.UUID
	}

	// GetForwardtestStatusWorkflowResults is the output for the GetForwardtestStatusWorkflow.
	GetForwardtestStatusWorkflowResults struct {
		Status forwardtest.Status
	}
)

// StartForwardtestWorkflowName is the name of the StartForwardtestWorkflow.
const StartForwardtestWorkflowName = "StartForwardtestWorkflow"

type (
	// StartForwardtestWorkflowParams is the input for the StartForwardtestWorkflow.
	StartForwardtestWorkflowParams struct {
		ForwardtestID uuid.UUID
	}

	// StartForwardtestWorkflowResults is the output for the StartForwardtestWorkflow.
	StartForwardtestWorkflowResults struct{}
)

const (
	// ServiceInfoWorkflowName is the name of the workflow to get the service info.
	ServiceInfoWorkflowName = "ServiceInfoWorkflow"
)

type (
	// ServiceInfoParams contains the parameters of the service info workflow.
	ServiceInfoParams struct{}

	// ServiceInfoResults contains the result of the service info workflow.
	ServiceInfoResults struct {
		Version string
	}
)

// StopForwardtestWorkflowName is the name of the StopForwardtestWorkflow.
const StopForwardtestWorkflowName = "StopForwardtestWorkflow"

type (
	// StopForwardtestWorkflowParams is the input for the StopForwardtestWorkflow.
	StopForwardtestWorkflowParams struct {
		ForwardtestID uuid.UUID
	}

	// StopForwardtestWorkflowResults is the output for the StopForwardtestWorkflow.
	StopForwardtestWorkflowResults struct{}
)

// SubscribeToPriceWorkflowName is the name of the SubscribeToPriceWorkflow.
const SubscribeToPriceWorkflowName = "SubscribeToPriceWorkflow"

type (
	// SubscribeToPriceWorkflowParams is the input for the SubscribeToPriceWorkflow.
	SubscribeToPriceWorkflowParams struct {
		ForwardtestID   uuid.UUID
		OnPriceCallback runtime.CallbackWorkflow
		Exchange        string
		Pair            string
	}

	// SubscribeToPriceWorkflowResults is the output for the SubscribeToPriceWorkflow.
	SubscribeToPriceWorkflowResults struct{}
)
