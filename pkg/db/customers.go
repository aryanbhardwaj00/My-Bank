package db

import (
	"context"
	"errors"
	"log"

	"github.com/Bank/pkg/models"
	"github.com/Bank/pkg/utils"
)

type customer struct {
	DB string
}

type Customer interface {
	InsertIntoDB(models.Customer) error
	UpdateCustomer(input string, cust models.Customer) error
	DeleteCustomer(models.Customer) error
	SearchCustomer(string) (int, error)
}

func NewCustomer() Customer {
	return &customer{}
}

func (c *customer) InsertIntoDB(input models.Customer) error {
	log.Println("inside db customer insert")
	// Insert Record into DB
	_, err := utils.Connection.NewInsert().Model(&input).Exec(context.Background())
	if err != nil {
		log.Println("Error in inserting data", err)
		return err
	}
	return nil
}

func (c *customer) DeleteCustomer(input models.Customer) error {

	res, err := utils.Connection.NewDelete().Model((&input)).Where("Name=?", input.Name).Exec(context.Background())
	if err != nil {
		log.Println("Error in deleting field", err)
		return err
	}
	log.Println("rows affected", res)
	result, err := res.RowsAffected()
	if result == 0 || err != nil {
		log.Println("No such record found", err)
		return errors.New("No such record found")
	}

	return nil
}

func (c *customer) SearchCustomer(input string) (int, error) {

	var cust models.Customer
	_, err := utils.Connection.NewSelect().Model(&cust).Where("Name=?", input).Exec(context.Background(), &cust)
	if err != nil {
		log.Println("Error in searching field", err)
		return cust.AccountNumber, err
	}

	return cust.AccountNumber, nil
}

func (c *customer) UpdateCustomer(input string, cust models.Customer) error {
	log.Println("inside db customer update", input, cust)

	result, err := utils.Connection.NewUpdate().Model(&cust).Where("Name=?", input).Exec(context.Background())
	if err != nil {
		log.Println("Error in updating data", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 || err != nil {
		log.Println("No such record found", err)
		return errors.New("No record")
	}
	log.Println(result)

	return nil
}
