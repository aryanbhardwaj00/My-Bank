package db

import (
	"testing"

	"github.com/Bank/pkg/models"
)

func TestCustomer(t *testing.T) {
	var a models.Customer
	c := NewCustomer()
	c.DeleteCustomer(a)
}
