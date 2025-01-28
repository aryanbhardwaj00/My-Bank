package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Customer struct {
	bun.BaseModel   `bun:"table:customer,alias:c"`
	UID             uuid.UUID         `json:"uid"`
	Name            string            `json:"name"`
	Age             int               `json:"age"`
	Contact         int               `json:"contact"`
	PrimaryEmail    string            `json:"primary_email"`
	SecondaryEmail  string            `json:"secondary_email"`
	CreatedAt       *time.Time        `json:"created_at"`
	UpdatedAt       *time.Time        `json:"updated_at"`
	Status          string            `json:"status"`
	AccountID       int               `json:"accountid"`
	Address         *Address           `json:"address"`
	ResourceVersion int64             `json:"resourceversion"`
	Tags            map[string]string `json:"tags"`
	Labels          map[string]string `json:"label"`
}

type Address struct {
	City  string
	State string
}

type Account struct {
	bun.BaseModel `bun:"table:accounts,alias:acc"`
	UID           uuid.UUID `json:"uid"`
	Type          string    `json:"type"`
	Balance       float32   `json:"balance"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CustomerID    int       `json:"customer_id"`
	Tags          []string  `json:"tags"`
}
