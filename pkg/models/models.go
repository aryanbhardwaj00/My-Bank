package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Customer struct {
	bun.BaseModel `bun:"table:customer,alias:c"`
	UID           uuid.UUID
	Name          string
	Age           int
	Account       Account
}

type Account struct {
	AccountNumber  int
	CurrentBalance float32
}
