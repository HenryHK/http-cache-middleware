package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// ContactRequest defines the structure of a contact in request
type ContactRequest struct {
	Contact Contact `json:"contact,omitempty"`
}

// Contact defines a structure used to add/update a contact via autopilot api
type Contact struct {
	ContactID         string      `json:"contact_id,omitempty"`
	Email             string      `json:"Email,omitempty"`
	Twitter           string      `json:"Twitter,omitempty"`
	FirstName         string      `json:"FirstName,omitempty"`
	LastName          string      `json:"LastName,omitempty"`
	Salutation        string      `json:"Salutation,omitempty"`
	Company           string      `json:"Company,omitempty"`
	NumberOfEmployees string      `json:"NumberOfEmployees,omitempty"`
	Title             string      `json:"Title,omitempty"`
	Industry          string      `json:"Industry,omitempty"`
	Phone             string      `json:"Phone,omitempty"`
	MobilePhone       string      `json:"MobilePhone,omitempty"`
	Fax               string      `json:"Fax,omitempty"`
	Website           string      `json:"Website,omitempty"`
	MailingStreet     string      `json:"MailingStreet,omitempty"`
	MailingCity       string      `json:"MailingCity,omitempty"`
	MailingState      string      `json:"MailingState,omitempty"`
	MailingPostalCode string      `json:"MailingPostalCode,omitempty"`
	MailingCountry    string      `json:"MailingCountry,omitempty"`
	Owner             string      `json:"owner_name,omitempty"`
	LeadSource        string      `json:"LeadSource,omitempty"`
	Status            string      `json:"Status,omitempty"`
	LinkedIn          string      `json:"LinkedIn,omitempty"`
	Unsubscribed      string      `json:"unsubscribed,omitempty"`
	Custom            interface{} `json:"custom,omitempty"`
	SessionID         string      `json:"_autopilot_session_id,omitempty"`
	List              string      `json:"_autopilot_list,omitempty"`
	Notify            string      `json:"notify,omitempty"`
}

// AddOrUpdateResponse defines the structure of returned value of add/update operation via autopilot api
type AddOrUpdateResponse struct {
	ContactID string `json:"contact_id,omitempty"`
}

// ReadAccessKey read access key from text file
func ReadAccessKey() string {
	file, err := os.Open("./access")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

// PostHandler queries autopilot api via post
func PostHandler(w http.ResponseWriter, r *http.Request) {
	accessKey := ReadAccessKey()
	// init client
	client := &http.Client{}
	// get request body
	req, _ := http.NewRequest(r.Method, "https://api2.autopilothq.com/v1/contact", r.Body)
	// This part of code is for ease of test
	if len(r.Header.Get("autopilotapikey")) == 0 {
		req.Header.Add("autopilotapikey", accessKey)
	} else {
		req.Header.Add("autopilotapikey", r.Header.Get("autopilotapikey"))
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(res.StatusCode)
	b, err := ioutil.ReadAll(res.Body)
	w.Write(b)
}

// GetHandler queries autopilot via get
func GetHandler(w http.ResponseWriter, r *http.Request) {
	accessKey := ReadAccessKey()
	// init client
	client := &http.Client{}
	// get params
	id := strings.TrimPrefix(r.URL.Path, "/contact/")
	fmt.Printf("Queried id: %v\n", id)
	req, _ := http.NewRequest(r.Method, fmt.Sprintf("https://api2.autopilothq.com/v1/contact/%s", id), nil)
	// This part of code is for ease of test
	if len(r.Header.Get("autopilotapikey")) == 0 {
		req.Header.Add("autopilotapikey", accessKey)
	} else {
		req.Header.Add("autopilotapikey", r.Header.Get("autopilotapikey"))
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(res.StatusCode)
	b, err := ioutil.ReadAll(res.Body)
	w.Write(b)
}
