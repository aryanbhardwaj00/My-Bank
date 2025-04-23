package db

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Bank/pkg/customerrors"
	"github.com/Bank/pkg/models"
	"github.com/Bank/pkg/utils"
)

type customer struct {
	DB string
}

type Customer interface {
	InsertIntoDB(models.Customer) error
	UpdateCustomerInDB(string, models.Customer) (models.Customer, error)
	DeleteCustomerInDB(string) error
	GetCustomerInDB(string) (models.Customer, error)
	ListCustomers([]string, string) ([]models.Customer, error)
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

func (c *customer) DeleteCustomerInDB(input string) error {
	var cst models.Customer
	sqlRes, err := utils.Connection.NewDelete().Model((&cst)).Where("uid=?", input).Exec(context.Background())

	if err != nil {
		log.Println("Error in deleting field", err)
		return err
	}

	rowsAffected, err := sqlRes.RowsAffected()
	if rowsAffected == 0 || err != nil {
		log.Println("No such record found", err)
		return customerrors.ErrNotFound
	}

	return nil
}

func (c *customer) GetCustomerInDB(input string) (models.Customer, error) {
	log.Println("Inside Search customer in DB")
	var cst models.Customer
	log.Println("Inside Search customer in DB, Value of Search Criteria", input)

	err := utils.Connection.NewSelect().Model(&cst).Where("uid=?", input).Scan(context.Background())

	log.Println("After DBcall:", cst)
	if err != nil {
		log.Println("Error in searching field", err)
		return cst, err
	}
	return cst, nil
}

func (c *customer) UpdateCustomerInDB(input string, updtCust models.Customer) (models.Customer, error) {

	// Update the DB and check for errors (if any)
	log.Println("Inside the update db function:")

	responseFromDb, err := utils.Connection.NewUpdate().Model(&updtCust).Where("uid=?", input).Exec(context.Background())
	log.Println("Executed the Update query")

	if err != nil {
		log.Println("Error in updating data", err)
		return models.Customer{}, err
	}

	log.Println("Verifying rows affected")
	rowsAffected, err := responseFromDb.RowsAffected()
	if rowsAffected == 0 || err != nil {
		log.Println("No such record found", err)
		return models.Customer{}, customerrors.ErrNotFound
	}
	log.Println("returning from update db function:")

	return updtCust, nil
}

func (c *customer) ListCustomers(arr []string, inp string) ([]models.Customer, error) {
	log.Println("Inside List customer in DB")
	rawQuery := fmt.Sprintf("select %v from customer %v ", strings.Join(arr, ","), inp)
	log.Println("RawQuery:", rawQuery)
	var cst []models.Customer
	// log.Println("Inside List customer in DB, Value of Search Criteria", input)

	_, err := utils.Connection.NewRaw(rawQuery).Exec(context.Background(), &cst)
	if err != nil {
		log.Println("Error in searching field", err)
		return cst, err
	}
	return cst, nil
}
