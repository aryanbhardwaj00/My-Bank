package db

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/Bank/pkg/models"
	"github.com/Bank/pkg/utils"
)

type customer struct {
	DB string
}

type Customer interface {
	InsertIntoDB(models.Customer) error
	UpdateCustomer(input string, cust models.Customer) (models.Customer, error)
	DeleteCustomer(models.Customer) error
	SearchCustomer(string) (models.Customer, error)
}

func NewCustomer() Customer {
	return &customer{}
}

func (c *customer) InsertIntoDB(input models.Customer) error {
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

	result, err := res.RowsAffected()
	if result == 0 || err != nil {
		log.Println("No such record found", err)
		return errors.New("No such record found")
	}

	return nil
}

func (c *customer) SearchCustomer(input string) (models.Customer, error) {

	var cst models.Customer

	err := utils.Connection.NewSelect().Model(&cst).Where("name=?", input).Scan(context.Background())
	if err != nil {
		log.Println("Error in searching field", err)
		return cst, err
	}

	return cst, nil
}

func (c *customer) UpdateCustomer(name string, cust models.Customer) (models.Customer, error) {
	result, err := utils.Connection.NewUpdate().Model(&cust).SetColumn("age", strconv.Itoa(cust.Age)).Where("name=?", name).Exec(context.Background())

	if err != nil {
		log.Println("Error in updating data", err)
		return models.Customer{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 || err != nil {
		log.Println("No such record found", err)
		return cust, errors.New("No record")
	}

	return cust, nil
}
