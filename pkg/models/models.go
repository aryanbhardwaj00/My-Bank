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
