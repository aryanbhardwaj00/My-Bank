package main

import (
	"log"
	"net/http"

	accountv1handler "github.com/Bank/pkg/api/v1/accounts"
	customerv1handler "github.com/Bank/pkg/api/v1/customers"
	"github.com/Bank/pkg/db"
	"github.com/Bank/pkg/services/moneyhandler"
	"github.com/Bank/pkg/utils"
	"github.com/gorilla/mux"
)

func main() {
	// Connect to Database
	err := utils.ConnectToDB()
	if err != nil {
		log.Fatalln("Exiting: Unable to connect to DB", err)
	}
	defer utils.Connection.Close()

	// Create a new Router
	// A router handles all the coming requests and direct them to the function to which they are bound

	// Created a new Router
	newRouter := mux.NewRouter()

	// Here we are binding the path of incoming request to the function
	// HandleFunc takes two arguments , Path and the Function

	newRouter.HandleFunc("/api/v1/customer", customerv1handler.CreateCustomer).Methods("POST")
	newRouter.HandleFunc("/api/v1/customers/{uid}", customerv1handler.GetCustomer).Methods("GET")
	newRouter.HandleFunc("/api/v1/customers/{uid}", customerv1handler.DeleteCustomer).Methods("DELETE")
	newRouter.HandleFunc("/api/v1/customers/{uid}", customerv1handler.UpdateCustomer).Methods("PATCH")
	newRouter.HandleFunc("/api/v1/customers", customerv1handler.ListCustomer).Methods("POST")

	// Path for Account related requests
	db := db.NewAccount()
	a := accountv1handler.NewAccount(db)
	newRouter.HandleFunc("/api/v1/accounts", a.CreateAccount).Methods("POST")
	newRouter.HandleFunc("/api/v1/accounts/{uid}", a.GetAccount).Methods("GET")
	newRouter.HandleFunc("/api/v1/accounts/{uid}", a.DeleteAccount).Methods("DELETE")
	newRouter.HandleFunc("/api/v1/accounts/{uid}", a.UpdateAccount).Methods("PATCH")

	//Path for Transaction related requests
	newRouter.HandleFunc("/services/transaction1/{uid}", moneyhandler.AddDeposit).Methods("POST")
	newRouter.HandleFunc("/services/transaction2/{uid}", moneyhandler.WithdrawMoney).Methods("POST")
	newRouter.HandleFunc("/services/transaction3/{uid}", moneyhandler.TransferMoney).Methods("POST")
	http.ListenAndServe(":8080", newRouter)
}
