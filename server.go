package main

import (
	"fmt"
	"log"
	"net/http"
)

var temp string = "test sttring"

func getValueHandler(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query()
	key := data["key"][0]
	value := GetValue(key)
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
