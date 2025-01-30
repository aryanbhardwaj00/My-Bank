package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Customer struct {
	bun.BaseModel   `bun:"table:customer,alias:c"`
	UID             uuid.UUID         `json:"uid" bun:"uid"`
	Name            string            `json:"name" bun:"name"`
	Age             int               `json:"age" bun:"age"`
	Contact         int               `json:"contact" bun:"contact"`
	PrimaryEmail    string            `json:"primary_email" bun:"primary_email"`
	SecondaryEmail  string            `json:"secondary_email" bun:"secondary_email"`
	CreatedAt       *time.Time        `json:"created_at" bun:"created_at"`
	UpdatedAt       *time.Time        `json:"updated_at" bun:"updated_at"`
	Status          string            `json:"status" bun:"status"`
	AccountID       int               `json:"accountid" bun:"account_id"`
	Address         *Address          `json:"address" bun:"address"`
	ResourceVersion int64             `json:"resourceversion" bun:"resource_version"`
	Tags            map[string]string `json:"tags" bun:"tags"`
	Labels          map[string]string `json:"label" bun:"labels"`
}

type Address struct {
	City  string
	State string
}

type Account struct {
	bun.BaseModel `bun:"table:accounts,alias:acc"`
	CustomerID    uuid.UUID  `json:"customer_id"`
	Type          string     `json:"type"`
	Balance       float32    `json:"balance"`
	Status        string     `json:"status"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	Tags          []string   `json:"tags"`
}
