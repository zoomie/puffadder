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

type reply struct {
	value string
	err   error
}

type command struct {
	typ           string
	accountName   string // during transaction this acts as the fromAccount
	toAccountName string
	amount        int
	replyChan     chan reply
}

type server struct {
	cmds chan<- command
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

func (s server) createAccount(w http.ResponseWriter, r *http.Request) {
	// need to cove the case when the user already exists in the system
	accountName := r.FormValue("accountName")
	replyChan := make(chan reply)
	s.cmds <- command{typ: createEvent, accountName: accountName, replyChan: replyChan}
	result := <-replyChan
	if result.err != nil {
		http.Error(w, "failed:"+result.err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, "account created with name:", accountName)
}

func (s server) viewCurrentAccount(w http.ResponseWriter, r *http.Request) {
	accountName := r.FormValue("accountName")
	replyChan := make(chan reply)
	s.cmds <- command{typ: viewEvent, accountName: accountName, replyChan: replyChan}
	result := <-replyChan
	if result.err != nil {
		http.Error(w, "failed:"+result.err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, result.value)
}

func (s server) addMoney(w http.ResponseWriter, r *http.Request) {
	accountName := r.FormValue("accountName")
	addAmount, err := strconv.Atoi(r.FormValue("addAmount"))
	if err != nil {
		http.Error(w, "addAmount is invalid", http.StatusBadRequest)
		return
	}
	replyChan := make(chan reply)
	s.cmds <- command{typ: addEvent, accountName: accountName, amount: addAmount, replyChan: replyChan}
	result := <-replyChan
	if result.err != nil {
		http.Error(w, "failed:"+result.err.Error(), http.StatusBadRequest)
	}
	fmt.Fprintln(w, "Added money, amount:", addAmount)
}

func (s server) withdrawMoney(w http.ResponseWriter, r *http.Request) {
	accountName := r.FormValue("accountName")
	subtractAmount, err := strconv.Atoi(r.FormValue("subtractAmount"))
	if err != nil {
		http.Error(w, "invalid subtract value", http.StatusBadRequest)
	}
	replyChan := make(chan reply)
	s.cmds <- command{typ: withdrawEvent, accountName: accountName, amount: subtractAmount, replyChan: replyChan}
	result := <-replyChan
	if result.err != nil {
		http.Error(w, "failed:"+result.err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, "withrdrew:", subtractAmount)
}

func (s server) transfer(w http.ResponseWriter, r *http.Request) {
	fromAccount := r.FormValue("fromAccount")
	toAccount := r.FormValue("toAccount")
	transferAmount, err := strconv.Atoi(r.FormValue("transferAmount"))
	if err != nil {
		http.Error(w, "transferAmount is not a number", http.StatusBadRequest)
		return
	}
	replyChan := make(chan reply)
	s.cmds <- command{typ: transactionEvent, accountName: fromAccount, toAccountName: toAccount, replyChan: replyChan}
	result := <-replyChan
	if result.err != nil {
		http.Error(w, "failed:"+result.err.Error(), http.StatusBadRequest)
	}
	fmt.Fprintln(w, "transaction successful, amount:", transferAmount)
}

func main() {
	projection := chooseProjection(algorithmType)
	createAccountProjection(projection)
	channelStream := setUpChannelStream(projection)
	s := server{cmds: channelStream}
	http.HandleFunc("/create-account", s.createAccount)
	http.HandleFunc("/view-current-account", s.viewCurrentAccount)
	http.HandleFunc("/add-money", s.addMoney)
	http.HandleFunc("/withdraw-money", s.withdrawMoney)
	http.HandleFunc("/transfer", s.transfer)

	log.Fatal(http.ListenAndServe(":8090", nil))
}
