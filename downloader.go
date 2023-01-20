package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type download struct {
	Title    string `json:"title,omitempty"`
	Location string `json:"location,omitempty"`
}

func status(respose http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(respose, "Hello!")
}

func handleDownloadRequest(response http.ResponseWriter, request *http.Request) {
	var downloadRequest download
	r, err := ioutil.ReadAll(request.Body)
	err = json.Unmarshal(r, &downloadRequest)
	if err != nil {
		http.Error(response, "bad request: "+err.Error(), 400)
		log.Println(err)
		return
	}
	log.Printf("%#v", downloadRequest)

	fmt.Fprintf(response, "Download!")
}

func main() {
	fmt.Println("Downloader")

	http.HandleFunc("/", status)
	http.HandleFunc("/download", handleDownloadRequest)
	http.ListenAndServe(":3000", nil)
}
