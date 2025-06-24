package entities

import (
	"encoding/json"
	"time"

	"github.com/cryptellation/forwardtests/pkg/forwardtest"
	"github.com/google/uuid"
)

// ForwardtestData is the data for a forwardtest.
type ForwardtestData struct {
	Accounts  map[string]Account `json:"accounts"`
	Orders    []Order            `json:"orders"`
	Callbacks Callbacks          `json:"callbacks"`
}

// Forwardtest is the entity for a forwardtest.
type Forwardtest struct {
	ID        string    `db:"id"`
	UpdatedAt time.Time `db:"updated_at"`
	Data      []byte    `db:"data"`
}

// ToModel converts a Forwardtest entity to a Forwardtest model.
func (ft Forwardtest) ToModel() (forwardtest.Forwardtest, error) {
	id, err := uuid.Parse(ft.ID)
	if err != nil {
		return forwardtest.Forwardtest{}, err
	}

	// Get data
	var data ForwardtestData
	if err := json.Unmarshal(ft.Data, &data); err != nil {
		return forwardtest.Forwardtest{}, err
	}

	orders, err := ToOrderModels(data.Orders)
	if err != nil {
		return forwardtest.Forwardtest{}, err
	}

	return forwardtest.Forwardtest{
		ID:        id,
		UpdatedAt: ft.UpdatedAt,
		Accounts:  ToAccountModels(data.Accounts),
		Orders:    orders,
		Callbacks: data.Callbacks.ToCallbacksModel(),
	}, nil
}

// FromForwardtestModel converts a Forwardtest model to a Forwardtest entity.
func FromForwardtestModel(ft forwardtest.Forwardtest) (Forwardtest, error) {
	data := ForwardtestData{
		Accounts:  FromAccountModels(ft.Accounts),
		Orders:    FromOrderModels(ft.Orders),
		Callbacks: FromCallbacksModel(ft.Callbacks),
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return Forwardtest{}, err
	}

	return Forwardtest{
		ID:        ft.ID.String(),
		UpdatedAt: time.Now(),
		Data:      dataBytes,
	}, nil
}
