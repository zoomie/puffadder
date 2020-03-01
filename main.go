package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

// The acciont projection is stored in a keyValueStore interface.
// The options are: hashTable, orderedArray, binaryTree.
const algorithmType = "binaryTree" // default
const lineOffSet = 31

type keyValueStore interface {
	get(key string) (int, bool)
	add(key string, value int)
}

type projectionStore struct {
	projection keyValueStore
}

var dataPath string

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

func chooseProjection(algType string) keyValueStore {
	if algType == "hashTable" {
		return &hashTable{}
	} else if algType == "orderedArray" {
		return &orderedArray{}
	} else {
		return &btree{}
	}
}

func createAccountProjection(projection keyValueStore) {
	file, _ := os.Open(dataPath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var currentOffset int64 = 0
	for scanner.Scan() {
		line := scanner.Text()
		name, eventType, changeAmount, err := decodeLine(line)
		if eventType == createEvent {
			projection.add(name, 0)
		} else if eventType == addEvent {
			currentAmount, _ := projection.get(name)
			updatedAmount := currentAmount + changeAmount
			projection.add(name, updatedAmount)

		} else if eventType == withdrawEvent {
			currentAmount, _ := projection.get(name)
			updatedAmount := currentAmount - changeAmount
			projection.add(name, updatedAmount)
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
}

func (p projectionStore) createAccount(w http.ResponseWriter, r *http.Request) {
	// need to cove the case when the user already exists in the system
	accountName := r.FormValue("accountName")
	err := createAccountEvent(p.projection, accountName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, "account created with name:", accountName)
}

func (p projectionStore) viewCurrentAccount(w http.ResponseWriter, r *http.Request) {
	accountName := r.FormValue("accountName")
	amount, ok := p.projection.get(accountName)
	if !ok {
		http.Error(w, "Account does not exist", http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, amount)
}

func (p projectionStore) addMoney(w http.ResponseWriter, r *http.Request) {
	accountName := r.FormValue("accountName")
	addAmount, err := strconv.Atoi(r.FormValue("addAmount"))
	if err != nil {
		http.Error(w, "addAmount is invalid", http.StatusBadRequest)
		return
	}
	err = addMoneyEvent(p.projection, accountName, addAmount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	fmt.Fprintln(w, "Added money, amount:", addAmount)
}

func (p projectionStore) withdrawMoney(w http.ResponseWriter, r *http.Request) {
	accountName := r.FormValue("accountName")
	subtractAmount, err := strconv.Atoi(r.FormValue("subtractAmount"))
	if err != nil {
		http.Error(w, "invalid subtract value", http.StatusBadRequest)
	}
	err = subtractMoneyEvent(p.projection, accountName, subtractAmount)
	if err != nil {
		http.Error(w, "not enough money", http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, "withrdrew:", subtractAmount)
}

func (p projectionStore) transfer(w http.ResponseWriter, r *http.Request) {
	fromAccount := r.FormValue("fromAccount")
	toAccount := r.FormValue("toAccount")
	transferAmount, err := strconv.Atoi(r.FormValue("transferAmount"))
	if err != nil {
		http.Error(w, "transferAmount is not a number", http.StatusBadRequest)
		return
	}
	fromAccountAmount, okFrom := p.projection.get(fromAccount)
	_, okTo := p.projection.get(toAccount)
	if !okFrom || !okTo {
		http.Error(w, "fromAccount and/or toAccount does not exist", http.StatusBadRequest)
		return
	}
	if transferAmount > fromAccountAmount {
		http.Error(w, "fromAccount does not have enouth money", http.StatusBadRequest)
		return
	}
	// add error handling to transactions between accounts
	err = subtractMoneyEvent(p.projection, fromAccount, transferAmount)
	if err != nil {
		// need to think about how to propgate erros here
		http.Error(w, "unable to subtract amount from account", http.StatusBadRequest)
		return
	}
	_ = addMoneyEvent(p.projection, toAccount, transferAmount)
	fmt.Fprintln(w, "transaction successful, amount:", transferAmount)
}

func main() {
	projection := chooseProjection(algorithmType)
	createAccountProjection(projection)
	store := projectionStore{projection: projection}
	http.HandleFunc("/create-account", store.createAccount)
	http.HandleFunc("/view-current-account", store.viewCurrentAccount)
	http.HandleFunc("/add-money", store.addMoney)
	http.HandleFunc("/withdraw-money", store.withdrawMoney)
	http.HandleFunc("/transfer", store.transfer)

	log.Fatal(http.ListenAndServe(":8090", nil))
}
