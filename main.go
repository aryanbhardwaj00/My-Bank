package main

import (
	"log"
	"net/http"

	customerv1handler "github.com/Bank/pkg/api/v1/customers"
	"github.com/Bank/pkg/utils"
	"github.com/gorilla/mux"
)

func main() {
	err := utils.ConnectToDB()
	if err != nil {
		log.Fatalln("Exiting: Unable to connect to DB", err)
	}
	defer utils.Connection.Close()
	// Create a new Router
	// A router handles all the coming requests and direct them to the function to which they are bound

	// Created a new Router
	log.Println("Creating new Router")
	newRouter := mux.NewRouter()

	// Here we are binding the path of incoming request to the function
	// HandleFunc takes two arguments , Path and the Function

	log.Println("Binding path with Function")
	newRouter.HandleFunc("/api/v1/customers", customerv1handler.CreateCustomer).Methods("POST")
	newRouter.HandleFunc("/api/v1/customers/{id}", customerv1handler.SearchCustomer).Methods("GET")
	newRouter.HandleFunc("/api/v1/customers/{abc}", customerv1handler.DeleteCustomer).Methods("DELETE")
	newRouter.HandleFunc("/api/v1/customers/{id}", customerv1handler.UpdateCustomer).Methods("PATCH")

	log.Println("Starting the Server")
	http.ListenAndServe(":8080", newRouter)
}
