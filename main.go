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

const algorithmType = "binaryTree"
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
	if algorithmType == "binaryTree" {
		currentIndex = &btree{}
	} else {
		currentIndex = &orderedArray{}
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
	}
	fmt.Fprintf(w, "account created")
}

func viewCurrentAccount(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query()
	accountName := data["accountName"][0]
	amount, err := getAmount(currentIndex, accountName)
	if err != nil {
		fmt.Fprintf(w, "error")
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
	} else {
		newAmount := currentAmount - subtractAmount
		_ = setAmount(currentIndex, accountName, newAmount)
	}
}

func exchange(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "still need to implement")
}

func main() {
	http.HandleFunc("/create-account", createAccount)
	http.HandleFunc("/view-current-account", viewCurrentAccount)
	http.HandleFunc("/add-money", addMoney)
	http.HandleFunc("/withdraw-money", withdrawMoney)
	http.HandleFunc("/exchange", exchange)

	log.Fatal(http.ListenAndServe(":8090", nil))
}
