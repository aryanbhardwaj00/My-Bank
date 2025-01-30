package customerv1handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/Bank/pkg/customerrors"
	"github.com/Bank/pkg/db"
	"github.com/Bank/pkg/models"
	"github.com/Bank/pkg/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var cust models.Customer

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cust)
	// json.Decoder , decodes/ reads from request body(r.body) and
	// Decode(&cust) will store whatever was read from request body and store it in "cust" variable

	if err != nil {
		log.Println("Error in reading from request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Set the UID and status
	cust.UID = uuid.New()
	cust.Status = "Active"

	// Check if the essential fields are empty or not
	if cust.Name == "" {
		log.Println("Empty Name")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Empty name field."))
		return
	}

	if cust.Age <= 0 {
		log.Println("Invalid value received for age")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Empty or invalid value for age."))
		return
	}

	if cust.Contact == 0 || len(strconv.Itoa(cust.Contact)) != 10 {
		log.Println("Invalid contact number")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Empty or invalid contact."))
		return
	}

	isValidEmail := func(email string) bool {
		_, err := mail.ParseAddress(email)
		return err == nil
	}

	if isValidEmail(cust.PrimaryEmail) == false {
		log.Println("Invalid Primary Email")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Empty or invalid email address."))
		return
	}

	if cust.SecondaryEmail != "" {
		if isValidEmail(cust.SecondaryEmail) == false {
			log.Println("Invalid Secondary Email")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Empty or invalid email address."))
			return
		}
	}

	if cust.AccountID <= 0 {
		log.Println("Invalid value received for Acc_Id")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid value for ACCOUNT ID."))
		return
	}

	if cust.Address.City == "" || cust.Address.State == "" {
		log.Println("Empty Address")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Empty ADDRESS field."))
		return
	}
	// Creating an Instance of Customer Interface, so that we can use its underlying methods
	c := db.NewCustomer()

	// Insert into Database or return error(if any)
	err = c.InsertIntoDB(cust)
	if err != nil {
		log.Println("Database insert failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("Successfully inserted data")

	// Send back response
	w.Write([]byte("Created new field."))
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside Delete function in handler")
	// Extract the path parameter[UID]
	mapOfPathParameters := mux.Vars(r)
	log.Println("Map of path parameter:", mapOfPathParameters)

	c := db.NewCustomer()

	err := c.DeleteCustomerInDB(mapOfPathParameters["uid"])

	if err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
			log.Println("No such record found.")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No such record found."))
			return
		} else {
			log.Println("Error in deleting record", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error in deleting record."))
			return
		}
	}

	w.Write([]byte("Deleted the respective field."))
}

func GetCustomer(w http.ResponseWriter, r *http.Request) {
	// Extract value from path parameters
	// Call the SearchCustomer function
	// Return if the requested field not found
	// If found , send it as response
	log.Println("Inside Search handler")
	m := mux.Vars(r)
	log.Println("Path parameters:", m)
	result, err := db.NewCustomer().GetCustomerInDB(m["uid"])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No such record found")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No record found."))
			return
		} else {
			log.Println("Cannot find the requested field", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Cannot find the requested field."))
			return
		}
	}
	log.Println("marshalling data")
	response, err := json.Marshal(result)
	if err != nil {
		log.Println("Unable to send", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	// Create a Variable to store the incoming data
	// Decode and store the incoming data and check for errors
	// Store the searching criteria(received through path parameter) in a variable
	// First fetch the record which needs to be updated
	// Make updates in fetched record by verifying non empty fields
	// Pass the updated variable in update customer function and update the DB
	log.Println("Inside Update Handler")
	var updtCust models.Customer
	err := json.NewDecoder(r.Body).Decode(&updtCust)
	if err != nil {
		log.Println("error in reading request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	searchingCriteria := mux.Vars(r)
	log.Println("Map of path parameters:", searchingCriteria)
	dbInterface := db.NewCustomer()

	oldCust, err := dbInterface.GetCustomerInDB(searchingCriteria["uid"])
	log.Println("Fetched data from DB:", oldCust)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("Error in searching for the requested field:", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No such record found."))
			return
		} else {
			log.Println("Error in searching: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server error."))
			return
		}

	}

	// Check fields which are not empty and make changes in record fetched earlier

	if updtCust.Name != oldCust.Name {
		if updtCust.Name == "" {
			log.Println("Empty name field.")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Empty Name field."))
			return
		}
		oldCust.Name = updtCust.Name
	}

	if updtCust.Age != oldCust.Age {
		if updtCust.Age <= 0 {
			log.Println("Invalid value received for age")
			w.WriteHeader(400)
			w.Write([]byte("Invalid value for age."))
			return
		}
		oldCust.Age = updtCust.Age
		log.Println("updated age", oldCust)
	}

	if updtCust.Contact != oldCust.Contact {
		if len(strconv.Itoa(updtCust.Contact)) == 10 && updtCust.Contact > 0 {
			oldCust.Contact = updtCust.Contact
			log.Println("updated contact", oldCust)
		} else {
			log.Println("Invalid contact")
			w.WriteHeader(400)
			w.Write([]byte("Invalid value for Contact."))
			return
		}
	}

	// Validate Email
	isValidEmail := func(email string) bool {
		_, err := mail.ParseAddress(email)
		return err == nil
	}

	if updtCust.PrimaryEmail != oldCust.PrimaryEmail {
		if isValidEmail(updtCust.PrimaryEmail) == true {
			oldCust.PrimaryEmail = updtCust.PrimaryEmail
			log.Println("updated primary email", oldCust)
		} else {
			log.Println("Invalid Primary Email")
			w.WriteHeader(400)
			w.Write([]byte("Invalid email."))
			return
		}
	}

	if updtCust.SecondaryEmail != oldCust.SecondaryEmail {
		if isValidEmail(updtCust.SecondaryEmail) == true {
			oldCust.SecondaryEmail = updtCust.SecondaryEmail
			log.Println("updated secondary email", oldCust)
		} else {
			log.Println("Invalid Email")
			w.WriteHeader(400)
			w.Write([]byte("Invalid Email."))
			return
		}
	}

	if updtCust.Status != oldCust.Status {
		if updtCust.Status == "" {
			log.Println("Invalid status")
			w.WriteHeader(400)
			w.Write([]byte("Invalid Status."))
			return
		}
		oldCust.Status = updtCust.Status
		log.Println("updated status", oldCust)
	}

	if updtCust.Address.City != oldCust.Address.City {
		if updtCust.Address.City == "" {
			log.Println("Invalid city")
			w.WriteHeader(400)
			w.Write([]byte("Invalid City."))
			return
		}
		oldCust.Address.City = updtCust.Address.City
	}

	if updtCust.Address.State != oldCust.Address.State {
		if updtCust.Address.State == "" {
			log.Println("Invalid State")
			w.WriteHeader(400)
			w.Write([]byte("Invalid Sity."))
			return
		}
		oldCust.Address.State = updtCust.Address.State
	}

	// After making all the necessary changes call the UpdateCustomer function
	updatedRecord, err := dbInterface.UpdateCustomerInDB(searchingCriteria["uid"], oldCust)
	log.Println("Called the update db function:", updatedRecord)

	if err != nil {
		log.Println("Error in updating requested field", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to update the requested field."))
		return
	}
	log.Println("Marshalling data")

	finalResponse, err := json.Marshal(updatedRecord)
	if err != nil {
		log.Println("Error in marshalling", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(finalResponse))
}

func SearchCustomers(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside List Api")
	var customers []models.Customer
	// A searching criteria[name,age,status,city,state]
	name := r.URL.Query().Get("name")
	age := r.URL.Query().Get("age")
	status := r.URL.Query().Get("status")
	city := r.URL.Query().Get("city")
	state := r.URL.Query().Get("state")

	if name == "" && age == "" && status == "" && city == "" && state == "" {
		log.Println("No query parameters passed")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Empty query parameters."))
		return
	}

	log.Println("name:", name, "age:", age, "state:", state, "city:", city, "status:", status)
	log.Println("Extracted query parameters")

	// Check on which criteria do we have to search
	if name != "" {
		log.Println("Verified name", name)
		err := utils.Connection.NewSelect().Model(&customers).Where("name=?", name).Scan(context.Background())
		log.Println("Called DB")
		if err != nil {
			log.Println("Error in searching requested fields", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(customers) == 0 {
			log.Println("No such records found")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No such record found."))
			return
		}
	}
	ageInInt, err := strconv.Atoi(age)
	if err != nil {
		log.Println("Error in converting age to int:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if ageInInt != 0 {
		if ageInInt < 0 {
			log.Println("Invalid Age")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid Age"))
			return
		}
		log.Println("Verified age")
		err := utils.Connection.NewSelect().Model(&customers).Where("age=?", ageInInt).Scan(context.Background())
		log.Println(err)
		if err != nil {
			log.Println("Error in searching requested fields", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(customers) == 0 {
			log.Println("No such records found")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No such record found."))
			return
		}
	}

	if status != "" {
		log.Println("Verified status")
		err := utils.Connection.NewSelect().Model(&customers).Where("status=?", status).Scan(context.Background())
		if err != nil {
			log.Println("Error in searching requested fields", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(customers) == 0 {
			log.Println("No such records found")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No such record found."))
			return
		}
	}

	if city != "" {
		log.Println("Verified city")
		err := utils.Connection.NewSelect().Model(&customers).Where("address->>'city'=?", city).Scan(context.Background())
		if err != nil {
			log.Println("Error in searching requested fields", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(customers) == 0 {
			log.Println("No such records found")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No such record found."))
			return
		}
	}

	if state != "" {
		log.Println("Verified state")
		err := utils.Connection.NewSelect().Model(&customers).Where("address->>'state'=?", state).Scan(context.Background())
		if err != nil {
			log.Println("Error in searching requested fields", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(customers) == 0 {
			log.Println("No such records found")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No such record found."))
			return
		}
	}
	log.Println("Before marshalling", customers)
	response, err := json.Marshal(customers)
	if err != nil {
		log.Println("Error in marshaling", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(response)
	log.Println("End of list api")
}
