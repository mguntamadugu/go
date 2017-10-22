/*
 * Implement a Simple REST service that tracks/maintains customer data
 *
 */

package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"

	"encoding/json"
)

type myHandler struct {
	name string
}

type customer struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Address string `json:"address"`
	Email string `json:"email"`
}

func (m myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello, world. I'm", m.name)
}

var customers []customer

func main() {

	// pre-initialize list of customers
	customers = append(customers, customer{"1", "Alfred", "A address", "alfred@gmail.com"})
	customers = append(customers, customer{"2", "Bob", "B address", "bob@gmail.com"})
	customers = append(customers, customer{"3", "Charles", "C address", "charles@gmail.com"})

	m := mux.NewRouter()
	m.HandleFunc("/customer", getCustomersHandler).Methods("GET") // Get All customers
	m.HandleFunc("/customer/{id}", getCustomerHandler).Methods("GET") // Get specific customer record
	m.HandleFunc("/customer/{id}", editCustomerHandler).Methods("PUT") // edit details of specific customer
	m.HandleFunc("/customer/{id}", createCustomerHandler).Methods("POST") // create new customer
	m.HandleFunc("/customer/{id}", deleteCustomerHandler).Methods("DELETE") // delete specific customer	 record
	http.ListenAndServe(":8080", m)
}

/*
 * Return all customer records in JSON format
 */
func getCustomersHandler(w http.ResponseWriter, r *http.Request) {

	// encode customer array into Json and pass to writer
	err := json.NewEncoder(w).Encode(customers)
	if err != nil {
		http.Error(w, "json encode error: " + err.Error(), http.StatusInternalServerError)
		return
	}
}

func getCustomerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var bFound bool
	var result customer

	for _, c := range customers {
		if c.Id == id {
			result = c
			bFound = true
			break
		}
	}

	// Customer record not found
	if bFound == false {
		// Should we return 204 )no content) or 404 (resource not found)
		http.Error(w, "Record not found for customer id: " + id, http.StatusNoContent)
		return
	}

	// encode found customer  into Json and pass to writer
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, "json encode error: " + err.Error(), http.StatusInternalServerError)
		return
	}


}

func deleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var bRemoved bool

	for arrIndex, c := range customers {
		if c.Id == id {

			customers = append(customers[:arrIndex], customers[arrIndex+1:]...)
			bRemoved = true
			break
		}
	}

	// Removed existing customer record or not
	if bRemoved == false {
		http.Error(w, "Record not found for customer id: " + id, http.StatusNotFound)
		return
	}

	// encode remaining customers list  into Json and pass to writer
	err := json.NewEncoder(w).Encode(customers)
	if err != nil {
		http.Error(w, "json encode error: " + err.Error(), http.StatusInternalServerError)
		return
	}


}

// update existing customer record
func editCustomerHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	var arrIndex int
	var bExisting bool
	//var cRecord customer

	for i, c := range customers {
		if c.Id == id {

			// locate index of specified customer record
			bExisting = true
			arrIndex = i
			//cRecord = c
			break
		}
	}

	// Found existing customer record to edit/replace
	if bExisting == false {
		http.Error(w, "Record does not already exist for customer id: " + id, http.StatusUnprocessableEntity)
		return
	}

	// Read JSON record into new customer object
	newCustomer := customer{}
	err := json.NewDecoder(r.Body).Decode(&newCustomer)
	if err != nil {
		http.Error(w, "json decode error: " + err.Error(), http.StatusInternalServerError)
		return
	}

	// Replace existing record with new record for specified customer
	customers[arrIndex] = newCustomer

	// encode new customers list  into Json and pass to writer
	err = json.NewEncoder(w).Encode(customers)
	if err != nil {
		http.Error(w, "json encode error: " + err.Error(), http.StatusInternalServerError)
		return
	}

}

// add a single new customer record
func createCustomerHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	var bExisting bool
	//var cRecord customer

	for _, c := range customers {
		if c.Id == id {

			bExisting = true
			//cRecord = c
			break;
		}
	}

	// Validate customer record does not exist
	if bExisting == true {
		http.Error(w, "Record already exists for customer id: " + id, http.StatusConflict)
		return
	}

	// Read JSON record into new customer object
	newCustomers := []customer{}
	err := json.NewDecoder(r.Body).Decode(&newCustomers)
	if err != nil {
		http.Error(w, "json decode error: " + err.Error(), http.StatusInternalServerError)
		return
	}

	// Validate customer id in URL matches customer id IN JSON object
	if id != newCustomers[0].Id {
		http.Error(w, "customer id in post request of url: " + id + " does not match record: " + newCustomers[0].Id,
			http.StatusBadRequest)
		return
	}

	// Append new customer(s) to list
	customers = append(customers, newCustomers...)

	// Send list
	err = json.NewEncoder(w).Encode(customers) // encode new customers list  into Json and pass to writer
	if err != nil {
		http.Error(w, "json encode error: " + err.Error(), http.StatusInternalServerError)
		return
	}


}


/*
func marshalExample() {
	d := dog{Color: "brown", Breed: "German Shepherd", Age: 5}
	b, _ := json.Marshal(&d)
	fmt.Println(string(b))
}

func unmarshalExample() {
	jsonStr := `{"color":"brown","breed":"German Shepherd","age":5}`
	d := dog{}
	json.Unmarshal([]byte(jsonStr), &d)
	fmt.Println(d)
}
*/