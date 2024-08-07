package customerv1handler

import (
	"encoding/json"
	"log"
	"net/http"

	// "github.com/Bank/pkg/db"
	"github.com/Bank/pkg/db"
	"github.com/Bank/pkg/models"
	"github.com/gorilla/mux"
)

func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var cust models.Customer

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cust)
	log.Println(cust)
	// json.Decoder , decodes/ reads from request body(r.body) and
	// Decode(&cust) will store whatever was read from request body and store it in cust variable

	if err != nil {
		log.Println("Error in reading from request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c := db.NewCustomer()
	err = c.InsertIntoDB(cust)
	if err != nil {
		log.Println("Database insert failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	a, err := json.Marshal(cust)
	if err != nil {
		log.Println("Unable to read", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := "Created new field"
	w.Write([]byte(response))
	w.Write(a)
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("INSIDE DELETE CUSTOMER IN HANDLER")

	var customer1 models.Customer
	err := json.NewDecoder(r.Body).Decode(&customer1)
	if err != nil {
		log.Println("Error in reading request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("READ THE REQUEST")

	c := db.NewCustomer()
	err = c.DeleteCustomer(customer1)
	if err != nil {
		log.Println("Data deletion failed", err)
		w.WriteHeader(http.StatusBadRequest)
		Response := "Respective field not found"
		w.Write([]byte(Response))
		return
	} else {
		log.Println("Successfully deleted data")
	}
	log.Println("DELETED DATA")
	Response := "Deleted the respective field"
	w.Write([]byte(Response))
}

func SearchCustomer(w http.ResponseWriter, r *http.Request) {
	// Read the request , reeturn error if fails
	// Call the SearchCustomer function
	// Return if the requested field not found
	// If found , send it as response
	m := mux.Vars(r)
	log.Println("Map of path parameter", m)

	result, err := db.NewCustomer().SearchCustomer(m["id"])
	if err != nil {
		log.Println("Could not find the requested field", err)
		w.WriteHeader(http.StatusBadRequest)
		response := "Cannot find the requested field"
		w.Write([]byte(response))
		return
	}

	a, err := json.Marshal(result)
	if err != nil {
		log.Println("Unable to send", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := "Account Number:"
	w.Write([]byte(response))
	w.Write(a)
}

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	var cust models.Customer

	err := json.NewDecoder(r.Body).Decode(&cust)
	if err != nil {
		log.Println("Error in reading request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m := mux.Vars(r)
	c := db.NewCustomer()
	err = c.UpdateCustomer(m["id"], cust)
	if err != nil {
		log.Println("Error in Updating data", err)
		response := "Failed to updated the requested field"
		w.Write([]byte(response))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response := "Updated the requested field"
	w.Write([]byte(response))
}
