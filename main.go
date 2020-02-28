package main

import (
	"bufio"
	"fmt"
	"log"
	"strconv"

	// "log"
	"net/http"
	"os"
	"path"
)

const algorithmType = "binaryTree" // default
const lineOffSet = 31

type keyValueStore interface {
	get(key string) (int, bool)
	add(key string, value int)
}

var dataPath string
var accountProjection keyValueStore

func setupDataFile() {
	// Check if data file exists
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dataPath = path.Join(workingDir, "data.puff")
	_, err = os.Stat(dataPath)
	if err == nil {
		fmt.Println("Data path exists at: ", dataPath)
		return
	}
	_, err = os.Create(dataPath)
	if err != nil {
		panic(fmt.Errorf("Could not create data file: %w", err))
	}
	fmt.Println("Data path created at:", dataPath)

}

func chooseIndex() {
	if algorithmType == "hashTable" {
		accountProjection = &hashTable{}
	} else if algorithmType == "orderedArray" {
		accountProjection = &orderedArray{}
	} else {
		accountProjection = &btree{}
	}
}

func createAccountProjection() {
	file, _ := os.Open(dataPath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var currentOffset int64 = 0
	for scanner.Scan() {
		line := scanner.Text()
		name, eventType, changeAmount, err := decodeLine(line)
		if eventType == createEvent {
			accountProjection.add(name, 0)
		} else if eventType == addEvent {
			currentAmount, _ := accountProjection.get(name)
			updatedAmount := currentAmount + changeAmount
			accountProjection.add(name, updatedAmount)

		} else if eventType == withdrawEvent {
			currentAmount, _ := accountProjection.get(name)
			updatedAmount := currentAmount - changeAmount
			accountProjection.add(name, updatedAmount)
		} else {
			panic(fmt.Errorf("incorrect event"))
		}
		if err != nil {
			panic(fmt.Errorf("data file corrupt: %w", err))
		}
		currentOffset += lineOffSet
	}
}

func init() {
	setupDataFile()
	chooseIndex()
	createAccountProjection()
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	// need to cove the case when the user already exists in the system
	accountName := r.FormValue("accountName")
	err := createAccountEvent(accountProjection, accountName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, "account created with name:", accountName)
}

func viewCurrentAccount(w http.ResponseWriter, r *http.Request) {
	accountName := r.FormValue("accountName")
	amount, ok := accountProjection.get(accountName)
	if !ok {
		http.Error(w, "Account does not exist", http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, amount)
}

func addMoney(w http.ResponseWriter, r *http.Request) {
	accountName := r.FormValue("accountName")
	addAmount, err := strconv.Atoi(r.FormValue("addAmount"))
	if err != nil {
		http.Error(w, "addAmount is invalid", http.StatusBadRequest)
		return
	}
	err = addMoneyEvent(accountProjection, accountName, addAmount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func withdrawMoney(w http.ResponseWriter, r *http.Request) {
	accountName := r.FormValue("accountName")
	subtractAmount, err := strconv.Atoi(r.FormValue("subtractAmount"))
	if err != nil {
		http.Error(w, "invalid subtract value", http.StatusBadRequest)
	}
	err = subtractMoneyEvent(accountProjection, accountName, subtractAmount)
	if err != nil {
		http.Error(w, "not enough money", http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, "withrdrew:", subtractAmount)
}

func transfer(w http.ResponseWriter, r *http.Request) {
	fromAccount := r.FormValue("fromAccount")
	toAccount := r.FormValue("toAccount")
	transferAmount, err := strconv.Atoi(r.FormValue("transferAmount"))
	if err != nil {
		http.Error(w, "transferAmount is not a number", http.StatusBadRequest)
		return
	}
	fromAccountAmount, okFrom := accountProjection.get(fromAccount)
	_, okTo := accountProjection.get(toAccount)
	if !okFrom || !okTo {
		http.Error(w, "fromAccount and/or toAccount does not exist", http.StatusBadRequest)
		return
	}
	if transferAmount > fromAccountAmount {
		http.Error(w, "fromAccount does not have enouth money", http.StatusBadRequest)
		return
	}
	// add error handling to transactions between accounts
	err = subtractMoneyEvent(accountProjection, fromAccount, transferAmount)
	if err != nil {
		// need to think about how to propgate erros here
		http.Error(w, "unable to subtract amount from account", http.StatusBadRequest)
		return
	}
	_ = addMoneyEvent(accountProjection, toAccount, transferAmount)
	fmt.Fprintln(w, "transaction successful, amount:", transferAmount)
}

func main() {
	http.HandleFunc("/create-account", createAccount)
	http.HandleFunc("/view-current-account", viewCurrentAccount)
	http.HandleFunc("/add-money", addMoney)
	http.HandleFunc("/withdraw-money", withdrawMoney)
	http.HandleFunc("/transfer", transfer)

	log.Fatal(http.ListenAndServe(":8090", nil))
}
