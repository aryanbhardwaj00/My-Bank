package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Customer struct {
	bun.BaseModel   `bun:"table:customer,alias:c"`
	UID             uuid.UUID
	Name            string
	Age             int
	PrimaryEmail    string
	SecondaryEmail  string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Status          string
	AccountID       int
	Address         int
	ResourceVersion time.Time
	Tags            []string
}

type Account struct {
	bun.BaseModel `bun:"table:account,alias:acc"`
	UID           uuid.UUID
	Type          string
	Balance       float32
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	CustomerID    int
	Tags          []string
}
