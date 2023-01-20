package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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
	if err != nil {
		http.Error(response, "bad request", 400)
		log.Println(err)
		return
	}
	defer request.Body.Close()

	err = json.Unmarshal(r, &downloadRequest)
	if err != nil {
		http.Error(response, "bad request: "+err.Error(), 400)
		return
	}
	log.Printf("%#v", downloadRequest)

	err = getFile(downloadRequest)
	if err != nil {
		http.Error(response, "internal server error", 500)
		return
	}

	fmt.Fprintf(response, "Download!")
}

func getFile(downloadRequest download) error {
	u, err := url.Parse(downloadRequest.Location)
	if err != nil {
		log.Println(err)
		return err
	}

	response, err := http.Get(strings.Replace(downloadRequest.Location, " ", "%20", -1))
	if err != nil {
		log.Println(err)
		return err
	}
	defer response.Body.Close()

	save, err := createSaveDirectory(downloadRequest.Title)
	if err != nil {
		log.Println(err)
		return err
	}

	out, err := os.Create(filepath.Join(save, filepath.Base(u.Path)))
	defer out.Close()
	_, err = io.Copy(out, response.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func createSaveDirectory(title string) (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(dir, title)
	_, err = os.Stat(path)

	// creating directory.
	if err != nil {
		err = os.Mkdir(path, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	return path, nil
}

func main() {
	fmt.Println("Downloader")

	http.HandleFunc("/", status)
	http.HandleFunc("/download", handleDownloadRequest)
	http.ListenAndServe(":3000", nil)
}
