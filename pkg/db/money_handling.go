package db

import (
	"context"
	"log"

	"github.com/Bank/pkg/customerrors"
	"github.com/Bank/pkg/models"
	"github.com/Bank/pkg/utils"
	"github.com/uptrace/bun"
)

type transaction struct {
	db *bun.DB
}

type Transactions interface {
	AddMoney(string, float32) (float32, error)
	WithdrawMoney(string, float32) (float32, error)
	TransferMoney(models.TransferMoney, string) (float32, error)
}

func NewTransaction() Transactions {
	return &transaction{db: utils.Connection}
}

func (tr *transaction) AddMoney(customerId string, amount float32) (float32, error) {
	// Create a variable to store the updated balance
	var updatedBalance float32

	// Update Balance in DB
	err := utils.Connection.NewUpdate().Table("accounts").Set("balance=balance+?", amount).Where("customer_id=?", customerId).Returning("balance").Scan(context.Background(), &updatedBalance)
	log.Println("executed query")
	// Check for errors
	if err != nil {
		log.Println("Error in updating balance.", err)
		return amount, err
	}

	// Return the updated balance
	return updatedBalance, nil
}

func (tr *transaction) WithdrawMoney(cstId string, amount float32) (float32, error) {
	// Create a variable to store the updated balance
	var currentBalance float32

	// Fetch record from DB
	log.Println(" Inside Withdrawn fn , executing query")
	err := utils.Connection.NewSelect().Table("accounts").Column("balance").Where("customer_id=?", cstId).Scan(context.Background(), &currentBalance)

	log.Println(" Current Balance:", currentBalance)
	log.Println(" Current Amount:", amount)

	log.Println(" Executed query")

	// Check for errors
	if err != nil {
		log.Println("Error in searching requested UID.", err)
		return amount, err
	}

	log.Println(" Validating Balance")

	// Validate Balance
	if currentBalance < amount {
		log.Println("Insufficient Balance.")
		return amount, customerrors.ErrInvalidInput
	}
	log.Println("Updating Balance")

	// Update balance after withdrawal
	_, err = utils.Connection.NewUpdate().Table("accounts").Set("balance=balance-?", amount).Where("customer_id=?", cstId).Exec(context.Background())
	if err != nil {
		log.Println("Error in updating balance.", err)
		return amount, err
	}

	// Return the remaining balance
	return currentBalance - amount, nil
}

func (tr *transaction) TransferMoney(acc models.TransferMoney, cstId string) (float32, error) {
	log.Println("Inside Transfer Money in DB.")
	log.Println("Creating a transaction")
	txn, err := tr.db.Begin()

	if err != nil {
		log.Println("Error in generating transaction", err)
		return acc.Amount, err
	}

	// First Withdraw from sender's account
	var ac models.Account
	log.Println("Searching the sender's account.")

	err = txn.NewSelect().Model(&ac).Where("customer_id=?", cstId).Scan(context.Background())
	if err != nil {
		log.Println("Error in searching requested UID.", err)
		return acc.Amount, err
	}
	log.Println("Verifying balance.")
	if ac.Balance < acc.Amount {
		log.Println("Invalid Balance")
		return ac.Balance, customerrors.ErrInvalidInput
	}
	var updatedBalance float32

	log.Println("Deducting amount")

	err = txn.NewUpdate().Model(&ac).Set("balance=balance-?", acc.Amount).Where("customer_id=?", cstId).Returning("balance").Scan(context.Background(), &updatedBalance)
	if err != nil {
		log.Println("Error in updating balance.", err)
		return acc.Amount, err
	}
	log.Println("Withdrawn from sender's account.")

	// Deposit it to receivers account
	res, err := txn.NewUpdate().Model(&ac).Set("balance=balance+?", acc.Amount).Where("customer_id=?", acc.ReceiverID).Exec(context.Background())
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("Executing Rollback")
		err = txn.Rollback()
		if err != nil {
			log.Println("Error in rollback.", err)
			return acc.Amount, err
		}
		log.Println("Error in updating balance.", err)
		return acc.Amount, err
	}

	if rowsAffected == 0 {
		log.Println("Executing Rollback due to no records to update.")
		err = txn.Rollback()
		if err != nil {
			log.Println("Error in rollback.", err)
			return acc.Amount, err
		}
		log.Println("No records were updated")
		return ac.Balance, customerrors.ErrNotFound
	}

	err = txn.Commit()

	return updatedBalance, nil
}
