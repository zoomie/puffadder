package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const dataPath = "/Users/andrew/go/src/github.com/zoomie/puffadder/data.puff"

func init() {
	if _, err := os.Stat(dataPath); err == nil {
		fmt.Println("Data path exists at: ", dataPath)
	} else {
		_, err := os.Create(dataPath)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Data path created at:", dataPath)
		}
	}
}

func getValueHandler(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query()
	key := data["key"][0]
	value := getValue(key)
	fmt.Fprintf(w, value)
}
func setValueHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "TODO")
}
func main() {
	http.HandleFunc("/get", getValueHandler)
	http.HandleFunc("/set", setValueHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
