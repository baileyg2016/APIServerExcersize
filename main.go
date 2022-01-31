package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

type ID struct {
	id int
}

type Application struct {
	License     string        `yaml:"license"`
	Description string        `yaml:"description"`
	Title       string        `yaml:"title"`
	Version     string        `yaml:"version"`
	Maintainers []Maintainers `yaml:"maintainers"`
	Company     string        `yaml:"company"`
	Website     string        `yaml:"website"`
	Source      string        `yaml:"source"`
}

// Maintainers
type Maintainers struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

var id int = 0
var applications map[int]Application

func getMetaData(w http.ResponseWriter, r *http.Request){

	queries:= r.URL.Query()
	resp := make([]Application, 0, 10)
	
	for k, v := range queries {
		filterApplications(&resp, k, v[0])
	}

	fmt.Fprintf(w, "\n")
	jsonResp, err := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err != nil {
		log.Fatalln(err)
	}
	
	w.Write(jsonResp)
}

func postMetaData(w http.ResponseWriter, r *http.Request){
	var payload Application
	bytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}
	
	if yaml.Unmarshal(bytes, &payload) != nil {
		log.Fatalln(err)
	}

	if !areMetaDataFieldsCorrect(payload) {
		fmt.Fprintf(w, "Invalid payload attribute")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		id++
		applications[id] = payload
		fmt.Fprintf(w, "Successfully added %+v", payload.Title)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp, _ := json.Marshal(ID{id: id})
		fmt.Println(resp)
		w.Write(resp)
	}
}

func metaData(w http.ResponseWriter, r *http.Request){
	switch r.Method {
		case "GET":
			getMetaData(w, r)
		case "POST":
			postMetaData(w, r)
	}
}

func retreiveApplication(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
    id, valid := vars["id"]
    if !valid {
        w.WriteHeader(http.StatusBadRequest)
    } else {
		id, err := strconv.Atoi(id)

		if err != nil {
			log.Fatalln(err)
		}

		resp, err := json.Marshal(applications[id])

		if err != nil {
			log.Fatalln(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func handleRequests() {
	r := mux.NewRouter()
    r.HandleFunc("/application", metaData)
	r.HandleFunc("/application/{id}/", retreiveApplication)

    log.Fatal(http.ListenAndServe(":8081", r))
}

func main() {
	applications = map[int]Application{}
	fmt.Println("Listing on port 8081")
	handleRequests()
}
