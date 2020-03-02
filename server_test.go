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
func setUp() server {
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
	projection := chooseProjection(algorithmType)
	channelStream := setUpChannelStream(projection)
	return server{cmds: channelStream}
}

func createAccountSetUp(accountName string) (*httptest.ResponseRecorder, *http.Request) {
	params := "?accountName=" + accountName
	url := path.Join(baseURL+"create-account", params)
	request, _ := http.NewRequest("POST", url, nil)
	recorder := httptest.NewRecorder()
	return recorder, request
}

func TestCreateAccount(t *testing.T) {
	s := setUp()
	recorder, request := createAccountSetUp("joe")
	s.createAccount(recorder, request)
	response := recorder.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Creating account failed")
	}
}

func TestViewAccount(t *testing.T) {
	s := setUp()
	s.createAccount(createAccountSetUp("joe"))

	params := "?accountName=joe"
	fullURL := path.Join(baseURL+"view-current-account", params)
	request, _ := http.NewRequest("GET", fullURL, nil)
	recorder := httptest.NewRecorder()
	s.viewCurrentAccount(recorder, request)
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
	s := setUp()
	s.createAccount(createAccountSetUp("joe"))

	recorder, request := AddSetUp("joe", "10")
	s.addMoney(recorder, request)
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
	s := setUp()
	accountName := "joe"
	s.createAccount(createAccountSetUp(accountName))
	s.addMoney(AddSetUp(accountName, "10"))

	recorder, request := withdrawSetUp(accountName, "5")
	s.withdrawMoney(recorder, request)
	response := recorder.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Test withdraw failed")
	}
}

func TestTransfer(t *testing.T) {
	s := setUp()
	fromAccountName := "joe"
	s.createAccount(createAccountSetUp(fromAccountName))
	s.addMoney(AddSetUp(fromAccountName, "10"))

	toAccount := "paul"
	s.createAccount(createAccountSetUp(toAccount))

	params := "?fromAccount=" + fromAccountName + "&toAccount=" + toAccount + "&transferAmount=5"
	fullURL := path.Join(baseURL+"transfer", params)
	request, _ := http.NewRequest("GET", fullURL, nil)
	recorder := httptest.NewRecorder()
	s.transfer(recorder, request)
	response := recorder.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Creating task failed")
	}
}
