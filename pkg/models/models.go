package models

import (
	"log"
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

func GetListColumns() []string {
	return []string{"name", "age", "uid"}
}

type Listparameters struct {
	PageNumber int
	PageSize   int
	Input      string
	OrderBy    string
	SortIn     string
}

func ValidateListParam(lst Listparameters) bool {
	log.Println("Inside Validate List Para fn")
	log.Println("Input Received:", lst)
	switch {
	case lst.PageNumber < 0, lst.PageSize <= 0:
		log.Println("Page number or page size less than 0:", lst.PageNumber)
		return false
	case lst.SortIn != "asc" && lst.SortIn != "desc":
		log.Println("Empty or invalid sort by.")
		return false
	case lst.OrderBy != "name" && lst.OrderBy != "updated_at" && lst.OrderBy != "created_at":
		log.Println("Empty or invalid order by.")
		return false
	case lst.Input == "":
		log.Println("Empty Input.")
		return false
	default:
		return true
	}
}
