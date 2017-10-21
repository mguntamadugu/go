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

func getCustomersHandler(w http.ResponseWriter, r *http.Request) {

	err := json.NewEncoder(w).Encode(customers) // encode customer array into Json and pass to writer
	if err != nil {
		fmt.Fprintf(w, err.Error())
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

	if bFound == false {
		w.Write([]byte("Record not found for customer id: " + id))
	} else {
		err := json.NewEncoder(w).Encode(result) // encode found customer  into Json and pass to writer
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
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
			break;
		}
	}

	if bRemoved == false {
		w.Write([]byte("Record not found for customer id: " + id))
	} else {
		err := json.NewEncoder(w).Encode(customers) // encode new customers list  into Json and pass to writer
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
	}

}

func editCustomerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world, edit")
}

func createCustomerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world, edit")
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