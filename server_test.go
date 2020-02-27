package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

// Find way to use different data file
// I could make the datafile and projection local vars?
func clearDataDir() {
	_, err := os.Stat(dataPath)
	if err == nil {
		_ = os.Remove(dataPath)
	}
	os.Create(dataPath)
}

func TestCreateAccount(t *testing.T) {
	clearDataDir()
	chooseIndex()
	url := "localhost:8090/create-account"
	params := "?accountName=joe&startingAmount=10"
	fullURL := path.Join(url, params)
	request, _ := http.NewRequest("GET", fullURL, nil)
	recorder := httptest.NewRecorder()
	createAccount(recorder, request)
	response := recorder.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Creating task failed")
	}
}

func TestViewAndAddMoney(t *testing.T) {
	clearDataDir()
	chooseIndex()

	urlCreate := "localhost:8090/create-account"
	paramsCreate := "?accountName=joe&startingAmount=10"
	fullURLCreate := path.Join(urlCreate, paramsCreate)
	requestCreate, _ := http.NewRequest("GET", fullURLCreate, nil)
	recorderCreate := httptest.NewRecorder()
	createAccount(recorderCreate, requestCreate)

	urlView := "localhost:8090/view-current-account"
	paramsView := "?accountName=joe"
	fullURLView := path.Join(urlView, paramsView)
	requestView, _ := http.NewRequest("GET", fullURLView, nil)
	recorderView := httptest.NewRecorder()
	viewCurrentAccount(recorderView, requestView)
	responseView := recorderView.Result()
	if responseView.StatusCode != http.StatusOK {
		t.Errorf("Creating task failed")
	}

	urlAdd := "localhost:8090/add-money"
	paramsAdd := "?accountName=joe&addAmount=10"
	fullURLAdd := path.Join(urlAdd, paramsAdd)
	requestAdd, _ := http.NewRequest("GET", fullURLAdd, nil)
	recorderAdd := httptest.NewRecorder()
	addMoney(recorderAdd, requestAdd)
	responseAdd := recorderAdd.Result()
	if responseAdd.StatusCode != http.StatusOK {
		t.Errorf("Creating task failed")
	}
}
