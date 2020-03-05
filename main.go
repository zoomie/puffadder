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

const (
	// The options are: hashTable, orderedArray, binaryTree.
	algorithmType = "hashTable" // default
	lineOffSet    = 31
)

var (
	dataPath string
)

type keyValueStore interface {
	get(key string) (int, bool)
	add(key string, value int)
}

type reply struct {
	value string
	err   error
}

type createCommand struct {
	accountName string
	replyChan   chan reply
}

type viewCommand struct {
	accountName string
	replyChan   chan reply
}

type addCommand struct {
	accountName string
	amount      int
	replyChan   chan reply
}
type withdrawCommand struct {
	accountName string
	amount      int
	replyChan   chan reply
}

type transactionCommand struct {
	toAccountName   string
	fromAccountName string
	amount          int
	replyChan       chan reply
}

type channelServer struct {
	createChan      chan<- createCommand
	viewChan        chan<- viewCommand
	addChan         chan<- addCommand
	withdrawChan    chan<- withdrawCommand
	transactionChan chan<- transactionCommand
}

func init() {
	// Setup the datafile
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

func chooseAlgorithm(algType string) keyValueStore {
	if algType == "hashTable" {
		return &hashTable{}
	} else if algType == "orderedArray" {
		return &orderedArray{}
	} else {
		return &btree{}
	}
}

func createAccountStore(store keyValueStore) {
	file, _ := os.Open(dataPath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var currentOffset int64 = 0
	for scanner.Scan() {
		line := scanner.Text()
		name, eventType, changeAmount, err := decodeLine(line)
		if eventType == createEvent {
			store.add(name, 0)
		} else if eventType == addEvent {
			currentAmount, _ := store.get(name)
			updatedAmount := currentAmount + changeAmount
			store.add(name, updatedAmount)

		} else if eventType == withdrawEvent {
			currentAmount, _ := store.get(name)
			updatedAmount := currentAmount - changeAmount
			store.add(name, updatedAmount)
		} else {
			panic(fmt.Errorf("incorrect event"))
		}
		if err != nil {
			panic(fmt.Errorf("data file corrupt: %w", err))
		}
		currentOffset += lineOffSet
	}
}

func setUpChannelServer(store keyValueStore) channelServer {
	createChan := make(chan createCommand)
	viewChan := make(chan viewCommand)
	addChan := make(chan addCommand)
	withdrawChan := make(chan withdrawCommand)
	transactionChan := make(chan transactionCommand)
	go func() {
		for {
			select {
			case createCmd := <-createChan:
				err := performCreate(store, createCmd.accountName)
				createCmd.replyChan <- reply{err: err}
			case veiwCmd := <-viewChan:
				amount, err := performView(store, veiwCmd.accountName)
				veiwCmd.replyChan <- reply{value: strconv.Itoa(amount), err: err}
			case addCmd := <-addChan:
				err := performAdd(store, addCmd.accountName, addCmd.amount)
				addCmd.replyChan <- reply{err: err}
			case withdrawCmd := <-withdrawChan:
				err := performSubtract(store, withdrawCmd.accountName, withdrawCmd.amount)
				withdrawCmd.replyChan <- reply{err: err}
			case transactionCmd := <-transactionChan:
				err := performTransaction(store, transactionCmd.fromAccountName,
					transactionCmd.toAccountName, transactionCmd.amount)
				transactionCmd.replyChan <- reply{err: err}
			}
		}
	}()
	return channelServer{createChan, viewChan, addChan, withdrawChan, transactionChan}
}

func (c channelServer) createAccount(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	accountName := r.Form.Get("accountName")
	replyChan := make(chan reply)
	c.createChan <- createCommand{accountName, replyChan}
	result := <-replyChan
	if result.err != nil {
		http.Error(w, "failed: "+result.err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, "account created with name:", accountName)
}

func (c channelServer) viewCurrentAccount(w http.ResponseWriter, r *http.Request) {
	accountName := r.FormValue("accountName")
	replyChan := make(chan reply)
	c.viewChan <- viewCommand{accountName, replyChan}
	result := <-replyChan
	if result.err != nil {
		http.Error(w, "failed: "+result.err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, result.value)
}

func (c channelServer) addMoney(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	accountName := r.Form.Get("accountName")
	addAmount, err := strconv.Atoi(r.Form.Get("addAmount"))
	if err != nil {
		http.Error(w, "addAmount is invalid", http.StatusBadRequest)
		return
	}
	replyChan := make(chan reply)
	c.addChan <- addCommand{accountName, addAmount, replyChan}
	result := <-replyChan
	if result.err != nil {
		http.Error(w, "failed:"+result.err.Error(), http.StatusBadRequest)
	}
	fmt.Fprintln(w, "added money, amount:", addAmount)
}

func (c channelServer) withdrawMoney(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	accountName := r.Form.Get("accountName")
	subtractAmount, err := strconv.Atoi(r.Form.Get("subtractAmount"))
	if err != nil {
		http.Error(w, "invalid subtract value", http.StatusBadRequest)
	}
	replyChan := make(chan reply)
	c.withdrawChan <- withdrawCommand{accountName, subtractAmount, replyChan}
	result := <-replyChan
	if result.err != nil {
		http.Error(w, "failed: "+result.err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, "withrdrew: ", subtractAmount)
}

func (c channelServer) transfer(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fromAccount := r.Form.Get("fromAccount")
	toAccount := r.Form.Get("toAccount")
	transferAmount, err := strconv.Atoi(r.Form.Get("transferAmount"))
	if err != nil {
		http.Error(w, "transferAmount is not a number", http.StatusBadRequest)
		return
	}
	replyChan := make(chan reply)
	c.transactionChan <- transactionCommand{toAccount, fromAccount, transferAmount, replyChan}
	result := <-replyChan
	if result.err != nil {
		http.Error(w, "failed:"+result.err.Error(), http.StatusBadRequest)
	}
	fmt.Fprintln(w, "transaction successful, amount: ", transferAmount)
}

func main() {
	store := chooseAlgorithm(algorithmType)
	createAccountStore(store)
	channelSrv := setUpChannelServer(store)
	http.HandleFunc("/create-account", channelSrv.createAccount)
	http.HandleFunc("/view-current-account", channelSrv.viewCurrentAccount)
	http.HandleFunc("/add-money", channelSrv.addMoney)
	http.HandleFunc("/withdraw-money", channelSrv.withdrawMoney)
	http.HandleFunc("/transfer", channelSrv.transfer)

	log.Fatal(http.ListenAndServe(":8090", nil))
}
