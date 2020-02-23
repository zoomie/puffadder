package main

import (
	"bufio"
	"fmt"
	"log"

	// "log"
	"net/http"
	"os"
	"path"
	"strconv"
)

const algorithmType = "binaryTree" // default
const lineOffSet = 22

type indexOffset interface {
	get(key string) (int64, bool)
	add(key string, value int64)
}

var dataPath string
var currentIndex indexOffset

func setupDataFile() {
	// Check if data file exists
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dataPath = path.Join(workingDir, "data.puff")
	if _, err := os.Stat(dataPath); err == nil {
		fmt.Println("Data path exists at: ", dataPath)
	} else {
		_, err := os.Create(dataPath)
		if err != nil {
			panic(fmt.Errorf("Could not create data file: %w", err))
		} else {
			fmt.Println("Data path created at:", dataPath)
		}
	}
}

func chooseIndex() {
	if algorithmType == "hashTable" {
		currentIndex = &hashTable{}
	} else if algorithmType == "orderedArray" {
		currentIndex = &orderedArray{}
	} else {
		currentIndex = &btree{}
	}
}

func loadInMemoryMapping() {
	file, _ := os.Open(dataPath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var currentOffset int64 = 0
	for scanner.Scan() {
		line := scanner.Text()
		name, _, err := decodeLine(line)
		if err != nil {
			panic(fmt.Errorf("data file corrupt: %w", err))
		}
		currentIndex.add(name, currentOffset)
		currentOffset += lineOffSet
	}
}

func init() {
	setupDataFile()
	chooseIndex()
	loadInMemoryMapping()
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	// need to cove the case when the user already exists in the system
	data := r.URL.Query()
	accountName := data["accountName"][0]
	startingAmount := data["startingAmount"][0]
	inputAmount, _ := strconv.Atoi(startingAmount)
	err := setAmount(currentIndex, accountName, inputAmount)
	if err != nil {
		fmt.Fprintf(w, "Set amount failed")
		return
	}
	fmt.Fprintf(w, "account created")
}

func viewCurrentAccount(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query()
	accountName := data["accountName"][0]
	amount, err := getAmount(currentIndex, accountName)
	if err != nil {
		fmt.Fprintf(w, "error")
		return
	}
	fmt.Println(amount)
}

func addMoney(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query()
	accountName := data["accountName"][0]
	addAmount, _ := strconv.Atoi(data["addAmount"][0])
	currentAmount, _ := getAmount(currentIndex, accountName)
	newAmount := currentAmount + addAmount
	_ = setAmount(currentIndex, accountName, newAmount)
}

func withdrawMoney(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query()
	accountName := data["accountName"][0]
	subtractAmount, _ := strconv.Atoi(data["subtractAmount"][0])
	currentAmount, _ := getAmount(currentIndex, accountName)
	if subtractAmount > currentAmount {
		fmt.Println("can't do transaction")
		return
	}
	newAmount := currentAmount - subtractAmount
	_ = setAmount(currentIndex, accountName, newAmount)
}

func transfer(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query()
	fromAccount := data["fromAccount"][0]
	toAccount := data["toAccount"][0]
	transferAmount, _ := strconv.Atoi(data["transferAmount"][0])
	fromAccountAmount, _ := getAmount(currentIndex, fromAccount)
	if transferAmount > fromAccountAmount {
		fmt.Println("from account does not have enouth money")
		return
	}
	_ = setAmount(currentIndex, fromAccount, fromAccountAmount-transferAmount)
	toCurrentAmount, _ := getAmount(currentIndex, toAccount)
	_ = setAmount(currentIndex, toAccount, toCurrentAmount+transferAmount)
}

func main() {
	http.HandleFunc("/create-account", createAccount)
	http.HandleFunc("/view-current-account", viewCurrentAccount)
	http.HandleFunc("/add-money", addMoney)
	http.HandleFunc("/withdraw-money", withdrawMoney)
	http.HandleFunc("/transfer", transfer)

	log.Fatal(http.ListenAndServe(":8090", nil))
}
