package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

var baseURL = "localhost:8090/"

// Find way to use different data file
// I could make the datafile and projection local vars?
func setUp() {
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dataPath = path.Join(workingDir, "test_data.puff")
	_, err = os.Stat(dataPath)
	if err == nil {
		os.Remove(dataPath)
	}
	os.Create(dataPath)
	// empty key value store after each test
	accountProjection = &btree{}
}

func createAccountSetUp(accountName string) (*httptest.ResponseRecorder, *http.Request) {
	params := "?accountName=" + accountName
	url := path.Join(baseURL+"create-account", params)
	request, _ := http.NewRequest("POST", url, nil)
	recorder := httptest.NewRecorder()
	return recorder, request
}

func TestCreateAccount(t *testing.T) {
	setUp()
	recorder, request := createAccountSetUp("joe")
	createAccount(recorder, request)
	response := recorder.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Creating account failed")
	}
}

func TestViewAccount(t *testing.T) {
	setUp()
	createAccount(createAccountSetUp("joe"))

	params := "?accountName=joe"
	fullURL := path.Join(baseURL+"view-current-account", params)
	request, _ := http.NewRequest("GET", fullURL, nil)
	recorder := httptest.NewRecorder()
	viewCurrentAccount(recorder, request)
	response := recorder.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("View account failed")
	}
}

func AddSetUp(accountName, amount string) (*httptest.ResponseRecorder, *http.Request) {
	params := "?accountName=" + accountName + "&addAmount=" + amount
	fullURL := path.Join(baseURL+"add-money", params)
	request, _ := http.NewRequest("POST", fullURL, nil)
	recorder := httptest.NewRecorder()
	return recorder, request
}

func TestAdd(t *testing.T) {
	setUp()
	createAccount(createAccountSetUp("joe"))

	recorder, request := AddSetUp("joe", "10")
	addMoney(recorder, request)
	response := recorder.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Adding money failed")
	}
}

func withdrawSetUp(accountName, amount string) (*httptest.ResponseRecorder, *http.Request) {
	params := "?accountName=" + accountName + "&subtractAmount=" + amount
	fullURL := path.Join(baseURL+"withdraw-money", params)
	request, _ := http.NewRequest("POST", fullURL, nil)
	recorder := httptest.NewRecorder()
	return recorder, request
}

func TestWithdraw(t *testing.T) {
	setUp()
	accountName := "joe"
	createAccount(createAccountSetUp(accountName))
	addMoney(AddSetUp(accountName, "10"))

	recorder, request := withdrawSetUp(accountName, "5")
	withdrawMoney(recorder, request)
	response := recorder.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Test withdraw failed")
	}
}

func TestTransfer(t *testing.T) {
	fromAccountName := "joe"
	createAccount(createAccountSetUp(fromAccountName))
	addMoney(AddSetUp(fromAccountName, "10"))

	toAccount := "paul"
	createAccount(createAccountSetUp(toAccount))

	params := "?fromAccount=" + fromAccountName + "&toAccount=" + toAccount + "&transferAmount=5"
	fullURL := path.Join(baseURL+"transfer", params)
	request, _ := http.NewRequest("GET", fullURL, nil)
	recorder := httptest.NewRecorder()
	transfer(recorder, request)
	response := recorder.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Creating task failed")
	}
}
