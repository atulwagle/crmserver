package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Customer struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	Email     string    `json:"email"`
	Phone     int       `json:"phone"`
	Contacted bool      `json:"contacted"`
}

type ErrorResponse struct {
	Errormessage string `json:"errormessage"`
}

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

var uuid_1, error1 = uuid.Parse("f3f1cc7d-1f32-4652-b016-52dc3ac5bf50")
var uuid_2, error2 = uuid.Parse("2f7a9959-084a-41d7-a85f-5430cbf6d90a")
var uuid_3, error3 = uuid.Parse("2532aacc-a5ae-4540-a22c-d04d889f142e")

var customer1 = Customer{uuid_1, "John", "Supervisor", "john@gmail.com", 4083457834, false}
var customer2 = Customer{uuid_2, "Bill", "Tester", "bill@gmail.com", 4083557844, false}
var customer3 = Customer{uuid_3, "Harry", "Manager", "harry@gmail.com", 4083657854, false}

var customerDb = []Customer{customer1, customer2, customer3}

func getCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	InfoLogger.Println("GetCustomers: Records successfully retrieved.")
	json.NewEncoder(w).Encode(customerDb)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tempCustomer = new(Customer)

	params := mux.Vars(r)
	var customerid = params["id"]
	if len(strings.TrimSpace(customerid)) == 0 {
		senderrormessage(w, "Invalid customer id.", http.StatusBadRequest)
		ErrorLogger.Println("GetCustomer: Invalid customer id.")
	} else {
		var found bool
		for i := 0; i < len(customerDb); i++ {
			if customerDb[i].Id.String() == customerid {
				tempCustomer.Id = customerDb[i].Id
				tempCustomer.Contacted = customerDb[i].Contacted
				tempCustomer.Email = customerDb[i].Email
				tempCustomer.Name = customerDb[i].Name
				tempCustomer.Phone = customerDb[i].Phone
				tempCustomer.Role = customerDb[i].Role
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(tempCustomer)
				found = true
				InfoLogger.Println("GetCustomer: Record found. Customerid:", customerid)
				break
			}
		}
		if !found {
			senderrormessage(w, "Record not found.", http.StatusNotFound)
			ErrorLogger.Println("GetCustomer: Record not found. Customerid:", customerid)
		}
	}
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	var customerid = params["id"]
	if len(strings.TrimSpace(customerid)) == 0 {
		senderrormessage(w, "Invalid customer id.", http.StatusBadRequest)
		ErrorLogger.Println("DeleteCustomer: Invalid customer id.")
	} else {
		var found bool
		for i := 0; i < len(customerDb); i++ {
			if customerDb[i].Id.String() == customerid {
				customerDb = append(customerDb[:i], customerDb[i+1:]...)
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(customerDb)
				found = true
				InfoLogger.Println("DeleteCustomer: Record successfully deleted. Customerid:", customerid)
				break
			}
		}
		if !found {
			senderrormessage(w, "Record not found.", http.StatusNotFound)
			ErrorLogger.Println("DeleteCustomer: Record not found. Customerid:", customerid)
		}
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "<html><body>")
	fmt.Fprintf(w, "<h1>List of supported APIs:</h1>\n")
	fmt.Fprintf(w, "<ul>")
	fmt.Fprintf(w, "<li>Getting a single customer through a /customers/{id} path</li>")
	fmt.Fprintf(w, "<li>Getting all customers through a the /customers path</li>")
	fmt.Fprintf(w, "<li>Creating a customer through a /customers path</li>")
	fmt.Fprintf(w, "<li>Updating a customer through a /customers/{id} path</li>")
	fmt.Fprintf(w, "<li>Deleting a customer through a /customers/{id} path</li>")
	fmt.Fprintf(w, "</ul>")
	fmt.Fprintf(w, "</body></html>")
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
	// Set the appropriate Content-Type in the response header
	w.Header().Set("Content-Type", "application/json")

	var newCustomerEntry = Customer{Id: uuid.New(), Name: "", Role: "", Email: "", Phone: 0, Contacted: false}

	// Read the HTTP request body
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorLogger.Printf("AddCustomer: Error during reading the request body: %v\n", err)
	}

	// Encode the request body into a Golang value so that we can work with the data
	err = json.Unmarshal(reqBody, &newCustomerEntry)
	if err != nil {
		ErrorLogger.Printf("AddCustomer: Error during unmarshalling: %v\n", err)
	}

	var found bool
	for i := 0; i < len(customerDb); i++ {
		if customerDb[i].Contacted == newCustomerEntry.Contacted &&
			customerDb[i].Email == newCustomerEntry.Email &&
			customerDb[i].Name == newCustomerEntry.Name &&
			customerDb[i].Role == newCustomerEntry.Role &&
			customerDb[i].Phone == newCustomerEntry.Phone {
			found = true

			senderrormessage(w, "Record already present.", http.StatusNotFound)
			ErrorLogger.Println("AddCustomer: Record already present. Customer Name:", newCustomerEntry.Name)

			break
		}
	}

	if !found {
		customerDb = append(customerDb, newCustomerEntry)
		w.WriteHeader(http.StatusCreated)
		InfoLogger.Println("AddCustomer: Record successfully added. Customer id:", newCustomerEntry.Id)

		// Regardless of successful resource creation or not, return the newly added customer record
		json.NewEncoder(w).Encode(newCustomerEntry)
	}
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	// Set the appropriate Content-Type in the response header
	w.Header().Set("Content-Type", "application/json")

	var updCustomerEntry = new(Customer)

	params := mux.Vars(r)
	var customerid = params["id"]
	if len(strings.TrimSpace(customerid)) == 0 {
		senderrormessage(w, "Invalid customer id.", http.StatusBadRequest)
		ErrorLogger.Println("UpdateCustomer: Invalid customer id.")
	} else {
		// Read the HTTP request body
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			ErrorLogger.Printf("UpdateCustomer: Error during reading the request body: %v\n", err)
		}

		// Encode the request body into a Golang value so that we can work with the data
		err = json.Unmarshal(reqBody, &updCustomerEntry)
		if err != nil {
			ErrorLogger.Printf("UpdateCustomer: Error during unmarshalling: %v\n", err)
		}

		var found bool
		for i := 0; i < len(customerDb); i++ {
			if customerDb[i].Id.String() == customerid {
				customerDb[i].Contacted = updCustomerEntry.Contacted
				customerDb[i].Email = updCustomerEntry.Email
				customerDb[i].Name = updCustomerEntry.Name
				customerDb[i].Role = updCustomerEntry.Role
				customerDb[i].Phone = updCustomerEntry.Phone
				found = true
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(customerDb[i])
				InfoLogger.Println("UpdateCustomer: Record updated successfully. CustomerId: ", customerid)
				break
			}
		}

		if !found {
			senderrormessage(w, "Record not found.", http.StatusNotFound)
			ErrorLogger.Println("UpdateCustomer: Record could not updated. Record not found. CustomerId: ", customerid)
		}
	}
}

func updateCustomers(w http.ResponseWriter, r *http.Request) {
	// Set the appropriate Content-Type in the response header
	w.Header().Set("Content-Type", "application/json")

	type UpdCustomersResponse struct {
		Customerentry Customer `json:"customerentry"`
		Statuscode    int      `json:"status"`
	}

	var updCustomerEntries []Customer

	// Read the HTTP request body
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorLogger.Printf("UpdateCustomers: Error during reading the request body: %v\n", err)
	}

	// Encode the request body into a Golang value so that we can work with the data
	err = json.Unmarshal(reqBody, &updCustomerEntries)
	if err != nil {
		ErrorLogger.Printf("UpdateCustomers: Error during unmarshalling: %v\n", err)
	}

	var updCustomerEntriesResponse []UpdCustomersResponse = make([]UpdCustomersResponse, len(updCustomerEntries))

	for j := 0; j < len(updCustomerEntries); j++ {
		var found bool
		for i := 0; i < len(customerDb); i++ {
			if customerDb[i].Id.String() == updCustomerEntries[j].Id.String() {
				customerDb[i].Contacted = updCustomerEntries[j].Contacted
				customerDb[i].Email = updCustomerEntries[j].Email
				customerDb[i].Name = updCustomerEntries[j].Name
				customerDb[i].Role = updCustomerEntries[j].Role
				customerDb[i].Phone = updCustomerEntries[j].Phone
				found = true

				updCustomerEntriesResponse[j].Customerentry = updCustomerEntries[j]
				updCustomerEntriesResponse[j].Statuscode = http.StatusOK
				InfoLogger.Println("UpdateCustomers: Record successfully updated. CustomerId:", updCustomerEntries[j].Id)
				break
			}
		}
		if !found {
			updCustomerEntriesResponse[j].Customerentry = updCustomerEntries[j]
			updCustomerEntriesResponse[j].Statuscode = http.StatusNotFound
			ErrorLogger.Println("UpdateCustomers: Record could not be updated. CustomerId:", updCustomerEntries[j].Id)
		}
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(updCustomerEntriesResponse)
	if err != nil {
		ErrorLogger.Printf("UpdateCustomers: Error during unmarshalling: %v\n", err)
	}
}

func senderrormessage(w http.ResponseWriter, message string, errorcode int) {
	var errormsg = new(ErrorResponse)
	errormsg.Errormessage = message
	w.WriteHeader(errorcode)
	json.NewEncoder(w).Encode(errormsg)
}

func init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	// Instantiate a new route
	router := mux.NewRouter()

	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	router.HandleFunc("/customers", getCustomers).Methods("GET")
	router.HandleFunc("/customers", addCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", updateCustomer).Methods("PATCH", "PUT")
	router.HandleFunc("/customers", updateCustomers).Methods("PATCH", "PUT")
	router.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")

	fmt.Println("Server is starting on port 3000...")
	InfoLogger.Println("Server is starting on port 3000...")

	http.ListenAndServe(":3000", router)
}
